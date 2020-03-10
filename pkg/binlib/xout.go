package binlib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

const XoutHdrLen = 16       /* header length */
const XoutSegEntryLen = 4   /* segment entry length */
const XoutRelocItemLen = 6  /* reloc item length */
const XoutSymbEntryLen = 12 /* symbol entry length */
const XoutNameLen = 8       /* symbol name length */

const XoutMagicSeg = 0xee00           /* segmented, non executable */
const XoutMagicSegX = 0xee01          /* segmented, executable */
const XoutMagicNonSeg = 0xee02        /* non segmented, non executable */
const XoutMagicNonSegX = 0xee03       /* non segmented, executable not, shared*/
const XoutMagicNonSegShared = 0x0006  /* non segmented, non executable, shared */
const XoutMagicNonSegXShared = 0xee07 /* non segmented, executable, shared */
const XoutMagicNonSegSplit = 0xee0a   /* non segmented, non executable, split ID */
const XoutMagicNonSegXSplit = 0xee0b  /* non segmented, executable, split ID */

const XoutSegBSS = byte(1)     /* bss segment */
const XoutSegSTACK = byte(2)   /* stack segment */
const XoutSegCODE = byte(3)    /* code segment */
const XoutSegCONST = byte(4)   /* constant pool */
const XoutSegDATA = byte(5)    /* initialized data */
const XoutSegCDMIX = byte(6)   /* mixed code and data, not protectable */
const XoutSegCDMIX_P = byte(7) /* mixed code and data, protectable */
const XoutSegUNDEF = byte(0)   /* linker assign the address */

const XoutRelocOFF = byte(1)  /* 16bit non segmented   */
const XoutRelocSSG = byte(2)  /* 16bit short segmented */
const XoutRelocLSG = byte(3)  /* 32bit long segmented  */
const XoutRelocXOFF = byte(5) /* 16bit non segmented referenced by external */
const XoutRelocXSSG = byte(6) /* short segmented referenced by external     */
const XoutRelocXLSG = byte(7) /* long segmented referenced by external      */

const XoutSymbLocal = byte(1)   /* local symbol */
const XoutSymbUndefEX = byte(2) /* undefined external */
const XoutSymbGlobal = byte(3)  /* global definition */
const XoutSymbSeg = byte(4)     /* segment name */

type XoutFile struct {
	Filep       *os.File
	Length      int64
	CodePos     int64
	RelocTblPos int64
	SymbTblPos  int64
	NumRelocs   int
	NumSymbs    int

	Header   XoutHeader
	SegTbl   []XoutSeg
	CodePart []byte
	RelocTbl []XoutRelocItem
	SymbTbl  []XoutSymbEntry
}

type XoutHeader struct {
	Magic       uint16
	NumSegs     int16
	CodePartLen int32
	RelocsLen   int32
	SymbsLen    int32
}

type XoutSeg struct {
	Number byte
	Type   byte
	Length uint16
}

type XoutRelocItem struct {
	SegIdx   byte
	Type     byte
	Location uint16
	SymbIdx  uint16
}

type XoutSymbEntry struct {
	SegIdx byte
	Type   byte
	Value  uint16
	Name   [8]byte
}

func (xf *XoutFile) ReadHdr() error {
	xf.Filep.Seek(0, 0)
	err := binary.Read(xf.Filep, binary.BigEndian, &xf.Header)
	if err != nil {
		return err
	}
	xf.CodePos = int64(XoutHdrLen + XoutSegEntryLen*xf.Header.NumSegs)
	xf.RelocTblPos = xf.CodePos + int64(xf.Header.CodePartLen)
	xf.SymbTblPos = xf.RelocTblPos + int64(xf.Header.RelocsLen)
	xf.NumRelocs = int(xf.Header.RelocsLen / XoutRelocItemLen)
	xf.NumSymbs = int(xf.Header.SymbsLen / XoutSymbEntryLen)
	return nil
}

func (xf *XoutFile) ReadSegTbl() error {
	xf.Filep.Seek(XoutHdrLen, 0)
	var seg XoutSeg
	for idx := 0; idx < int(xf.Header.NumSegs); idx++ {
		err := binary.Read(xf.Filep, binary.BigEndian, &seg)
		if err != nil {
			return errors.New("Segment table read error")
		}
		xf.SegTbl = append(xf.SegTbl, seg)
	}
	return nil
}

func (xf *XoutFile) ReadCodePart() error {
	xf.Filep.Seek(xf.CodePos, 0)
	err := binary.Read(xf.Filep, binary.BigEndian, xf.CodePart)
	if err != nil {
		return errors.New("Code part read error")
	}
	bytes.NewBuffer(xf.CodePart)
	return nil
}

func (xf *XoutFile) ReadRelocTbl() error {
	xf.Filep.Seek(xf.RelocTblPos, 0)
	for idx := 0; idx < xf.NumRelocs; idx++ {
		var rbuf XoutRelocItem
		err := binary.Read(xf.Filep, binary.BigEndian, &rbuf)
		if err != nil {
			return errors.New("Relocation table read error")
		}
		if rbuf.Type == 0 {
			continue
		}
		xf.RelocTbl = append(xf.RelocTbl, rbuf)
		/*
			fmt.Printf("%d : %x %x %x %d\n", idx,
				xoutRelocTbl[idx].SegIdx,
				xoutRelocTbl[idx].Type,
				xoutRelocTbl[idx].Location,
				xoutRelocTbl[idx].SymbIdx)
		*/
	}
	xf.NumRelocs = len(xf.RelocTbl)
	return nil
}

func (xf *XoutFile) ReadSymbTbl() error {
	xf.Filep.Seek(xf.SymbTblPos, 0)
	for idx := 0; idx < xf.NumSymbs; idx++ {
		var rbuf XoutSymbEntry
		err := binary.Read(xf.Filep, binary.BigEndian, &rbuf)
		if err != nil {
			return errors.New("Symbol table read error")
		}
		xf.SymbTbl = append(xf.SymbTbl, rbuf)
	}
	return nil
}

func (xf *XoutFile) Read(file *os.File) error {
	xf.Filep = file
	xoutinfo, err := file.Stat()
	if err != nil {
		return errors.New("can not get file status")
	}
	xf.Length = xoutinfo.Size()
	res := xf.ReadHdr()
	if res != nil {
		return res
	}
	res = xf.ReadSegTbl()
	if res != nil {
		return res
	}
	xf.CodePart = make([]byte, xf.Header.CodePartLen)
	res = xf.ReadCodePart()
	if res != nil {
		return res
	}
	xf.RelocTbl = make([]XoutRelocItem, 0, 1024)
	res = xf.ReadRelocTbl()
	if res != nil {
		return res
	}
	xf.SymbTbl = make([]XoutSymbEntry, 0, 1024)
	res = xf.ReadSymbTbl()
	if res != nil {
		return res
	}
	return nil
}

/*
func (xf *XoutFile) PrintRelocTbl() {
	for idx := 0; idx < xf.NumRelocs; idx++ {
		fmt.Printf("%3d : Seg %d, Type %d, Loc %04X, Symb %d\n", idx,
			xf.RelocTbl[idx].SegIdx,
			xf.RelocTbl[idx].Type,
			xf.RelocTbl[idx].Location,
			xf.RelocTbl[idx].SymbIdx)
	}
}

func (xf *XoutFile) PrintSymbTbl() {
	for idx := 0; idx < xf.NumSymbs; idx++ {
		fmt.Printf("%3d : Seg %3d, Type %d, Val %04X, %s\n", idx,
			xf.SymbTbl[idx].SegIdx,
			xf.SymbTbl[idx].Type,
			xf.SymbTbl[idx].Value,
			xf.SymbTbl[idx].Name)
	}
}
*/
