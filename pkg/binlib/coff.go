package binlib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

const CoffHdrLen = 20
const CoffSectHdrLen = 40
const CoffRelocItemLen = 16
const CoffSymbEntryLen = 18

const CoffNameLen = 8
const CoffLongNameLen = 32

type CoffFile struct {
	Filep    *os.File
	Header   CoffHdr
	SectTbl  []CoffSectHdr
	RelocTbl []CoffRelocItem
	SymbTbl  []interface{}
	CodePart *[]byte
}

type CoffHdr struct {
	Magic       uint16
	NumSects    uint16
	Date        uint32
	SymbTblFpos int32
	NumSymbs    uint32
	OptHdrLen   uint16
	Flags       uint16
}

type CoffSectHdr struct {
	Name         [CoffNameLen]byte
	Paddr        uint32
	Vaddr        uint32
	Length       uint32
	Fpos         int32
	RelocTblFpos int32
	LineNumsFpos int32
	NumRelocs    uint16
	NumLines     uint16
	Flags        uint32
}

type CoffRelocItem struct {
	Vaddr   uint32
	SymbIdx uint32
	Offset  uint32
	Type    uint16
	Stuff   uint16
}

const CoffSectTEXT = uint32(0x0020)
const CoffSectDATA = uint32(0x0040)
const CoffSectBSS = uint32(0x0080)

type CoffSymbEntry struct {
	Name      [CoffNameLen]byte
	Value     uint32
	SectNo    int16
	Type      uint16
	StrgClass byte
	NumAux    byte
}

type CoffSymbAuxSect struct {
	Length    uint32
	NumRelocs uint16
	NumLines  uint16
	Dummy     [10]byte
}

type CoffSymbAuxFile struct {
	Name [18]byte
}

const CoffSymbClassLocal = byte(0x01)
const CoffSymbClassGlobal = byte(0x02)
const CoffSymbClassStatic = byte(0x03)
const CoffSymbClassExternal = byte(0x05)
const CoffSymbClassLabel = byte(0x06)
const CoffSymbClassFile = byte(0x67)

const CoffSymbSCNExt = int16(0)
const CoffSymbSCNAbs = int16(-1)

func (cf *CoffFile) Open(file *os.File) {
	cf.Filep = file
}

func (cf *CoffFile) WriteHdr() error {
	err := binary.Write(cf.Filep, binary.BigEndian, cf.Header)
	if err != nil {
		return errors.New("Coff Header write error")
	}
	return nil
}

func (cf *CoffFile) WriteCodePart() error {
	err := binary.Write(cf.Filep, binary.BigEndian, *cf.CodePart)
	if err != nil {
		return errors.New("Coff Code part write error")
	}
	return nil
}

func (cf *CoffFile) WriteSectTbl() error {
	for idx := 0; idx < int(cf.Header.NumSects); idx++ {
		err := binary.Write(cf.Filep, binary.BigEndian, cf.SectTbl[idx])
		if err != nil {
			return errors.New("Coff Section write error")
		}
	}
	return nil
}

func (cf *CoffFile) WriteRelocTbl() error {
	for _, reloc := range cf.RelocTbl {
		err := binary.Write(cf.Filep, binary.BigEndian, reloc)
		if err != nil {
			return errors.New("Coff Reloc table write error")
		}
	}
	return nil
}

func (cf *CoffFile) WriteSymbTbl() error {
	for _, symb := range cf.SymbTbl {
		err := binary.Write(cf.Filep, binary.BigEndian, symb)
		if err != nil {
			fmt.Println(err)
			return errors.New("Coff Symbol table write error")
		}
	}
	return nil
}

/*
func (cf *CoffFile) PrintHdr() {
	fmt.Printf("Coff Header\n")
	fmt.Printf("Magic %04X, Flags %04X, NSects %d, SymbTPos %04X, NSymbs %d\n",
		cf.Header.Magic, cf.Header.Flags, cf.Header.NumSects,
		cf.Header.SymbTblFpos, cf.Header.NumSymbs)
	fmt.Println()
}

func (cf *CoffFile) CoffPrintSectTbl() {
	fmt.Printf("Section Table %d items\n", len(cf.SectTbl))
	for idx, sect := range cf.SectTbl {
		fmt.Printf("%d : Vaddr %04X, Paddr %04X, Len %d, Name %s\n", idx,
			sect.Vaddr, sect.Paddr, sect.Length, sect.Name)
		fmt.Printf("    Pos %04X, RTblPos %04X, NRelocs %d, Flags %04X\n",
			sect.Fpos, sect.RelocTblFpos, sect.NumRelocs, sect.Flags)
	}
	fmt.Println()
}

func (cf *CoffFile) CoffPrintRelocTbl() {
	fmt.Printf("Relocation Table %d entries\n", len(cf.RelocTbl))
	for idx, reloc := range cf.RelocTbl {
		fmt.Printf("%3d : Vaddr %04X, RType %02x, Symb %d\n", idx,
			reloc.Vaddr, reloc.Type, reloc.SymbIdx)
	}
	fmt.Println()
}

func (cf *CoffFile) CoffPrintSymbTbl() {
	fmt.Printf("Symbol Table %d entries\n", len(cf.SymbTbl))
	for idx, entry := range cf.SymbTbl {
		if symb, ok := entry.(CoffSymbEntry); ok {
			fmt.Printf("%3d : Sect %d, Type %d, Name %s\n", idx, symb.SectNo, symb.Type, symb.Name)
		}
	}
	fmt.Println()
}
*/
