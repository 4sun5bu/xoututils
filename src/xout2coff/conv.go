/*
 *  conv.go
 *
 *  Copyright (c) 2020 4sun5bu
 *  Released under the MIT license.
 *  See LICENSE.
 *
 *  A converter from XOUT to COFF
 *  Converted files can be linked with GNU ld.
 */

package main

import (
	"fmt"
	"sort"

	"binlib"
)

func assignBSS(xf *binlib.XoutFile) {
	// Look for the BSS Segment. If no BSS, create it.
	var bss int
	for bss = 0; bss < int(xf.Header.NumSegs); bss++ {
		if xf.SegTbl[bss].Type == binlib.XoutSegBSS {
			break
		}
	}
	if bss == int(xf.Header.NumSegs) {
		var bssSeg binlib.XoutSeg
		bssSeg.Number = 0xff
		bssSeg.Type = binlib.XoutSegBSS
		bssSeg.Length = uint16(0)
		xf.SegTbl = append(xf.SegTbl, bssSeg)
		xf.Header.NumSegs++
	}
	// Add memory space in the BSS segment.
	/*
		var symb *binlib.XoutSymbEntry
		for idx := 0; idx < xf.NumSymbs; idx++ {
			symb = &xf.SymbTbl[idx]
			if symb.SegIdx == 0xff && symb.Type == binlib.XoutSymbUndefEX && symb.Value != 0 {
				symb.Type = binlib.XoutSymbGlobal
				symb.SegIdx = byte(bss)
				size := symb.Value
				symb.Value = xf.SegTbl[bss].Length
				xf.SegTbl[bss].Length += size
				//fmt.Printf("%8s %d %d %04X\n", binlib.ConvertName(xf.SymbTbl[idx].Name), xf.SymbTbl[idx].Type, xf.SymbTbl[idx].SegIdx, xf.SymbTbl[idx].Value)
			}
		}
	*/
	// Add memory space in the BSS segment.
	for idx, symb := range xf.SymbTbl {
		if symb.SegIdx == 0xff && symb.Type == binlib.XoutSymbUndefEX && symb.Value != 0 {
			xf.SymbTbl[idx].Type = binlib.XoutSymbGlobal
			xf.SymbTbl[idx].SegIdx = byte(bss)
			size := symb.Value
			xf.SymbTbl[idx].Value = xf.SegTbl[bss].Length // set address
			xf.SegTbl[bss].Length += size
		}
	}
	//fmt.Printf("bss len %d\n", xf.SegTbl[bss].Length)
}

func calcAddr(xf *binlib.XoutFile, seg int, offset uint16) uint16 {
	pos := uint16(0)
	for idx := 0; ; {
		if idx == seg {
			break
		}
		pos += xf.SegTbl[idx].Length
		idx++
	}
	return pos + offset
}

func searchSymb(xf *binlib.XoutFile, name string) int {
	for idx, symb := range xf.SymbTbl {
		if string(symb.Name[:]) == name {
			return idx
		}
	}
	return -1
}

/* Add symbols for reloc items using offset from segment top */
/*
func addLocalSymb(xf *binlib.XoutFile) {
	nextIdx := len(xf.SymbTbl) - 1
	for ridx, reloc := range xf.RelocTbl {
		if reloc.Type == binlib.XoutRelocOFF {
			// get offset in the segment
			pos := calcAddr(xf, int(reloc.SegIdx), reloc.Location)
			offset := uint16(xf.CodePart[pos])*256 + uint16(xf.CodePart[pos+1])
			// generate a symbol name, and check if it has been existed already,
			name := fmt.Sprintf("SEG%d%04X", reloc.SymbIdx, offset)
			symIdx := searchSymb(xf, name)
			if symIdx == -1 {
				// make new symbol entry
				var symb binlib.XoutSymbEntry
				symb.SegIdx = byte(reloc.SymbIdx)
				symb.Type = binlib.XoutSymbLocal
				symb.Value = offset
				copy(symb.Name[:], []byte(name)[0:8])
				xf.SymbTbl = append(xf.SymbTbl, symb)
				nextIdx++
				xf.NumSymbs++
				xf.RelocTbl[ridx].SymbIdx = uint16(nextIdx)
			} else {
				xf.RelocTbl[ridx].SymbIdx = uint16(symIdx)
			}
		} else {
			continue
		}
	}
}
*/

/* add symbols for each segment top address */
func addSegTopSymb(xf *binlib.XoutFile) {
	for idx, _ := range xf.SegTbl {
		name := fmt.Sprintf("SEG%d0000", idx)
		//fmt.Println(name, idx)
		var symb binlib.XoutSymbEntry
		symb.SegIdx = byte(idx)
		symb.Type = binlib.XoutSymbLocal
		symb.Value = 0x0000
		copy(symb.Name[:], []byte(name)[0:8])
		xf.SymbTbl = append(xf.SymbTbl, symb)
		xf.NumSymbs++
	}
}

func convSegName(segType byte) string {
	var name string
	switch segType {
	case binlib.XoutSegCODE:
		name = ".text"
	case binlib.XoutSegDATA:
		name = ".data"
	case binlib.XoutSegCONST:
		name = ".rdata"
	case binlib.XoutSegBSS:
		name = ".bss"
	default:
		name = ""
	}
	return name
}

/* Add segmemt Symbols int the table */
func addSegSymb(xf *binlib.XoutFile) int {
	for segIdx, seg := range xf.SegTbl {
		var symbIdx int
		var symb *binlib.XoutSymbEntry
		for symbIdx = 0; symbIdx < len(xf.SymbTbl); symbIdx++ {
			symb = &xf.SymbTbl[symbIdx]
			if symb.Type != binlib.XoutSymbSeg {
				continue
			}
			if int(symb.SegIdx) == segIdx {
				name := convSegName(seg.Type)
				copy(symb.Name[:], []byte(name)[0:8])
				break
			}
		}
		if symbIdx == len(xf.SymbTbl) {
			var segSymb binlib.XoutSymbEntry
			segSymb.Type = binlib.XoutSymbSeg
			segSymb.SegIdx = uint8(segIdx)
			segSymb.Value = 0
			name := convSegName(seg.Type)
			copy(segSymb.Name[:], []byte(name)[0:8])
			xf.SymbTbl = append(xf.SymbTbl, segSymb)
			xf.NumSymbs++
		}
	}
	return 0
}

// ConvHdr has to be called after converting sections, relocations and symbols
func convHdr(cf *binlib.CoffFile) {
	cf.Header.Magic = 0x8000
	cf.Header.NumSects = uint16(len(cf.SectTbl))
	cf.Header.Date = 0x00000000
	cf.Header.SymbTblFpos = int32(binlib.CoffHdrLen + len(cf.SectTbl)*binlib.CoffSectHdrLen +
		len(*cf.CodePart) + len(cf.RelocTbl)*binlib.CoffRelocItemLen)
	cf.Header.NumSymbs = uint32(len(cf.SymbTbl))
	cf.Header.OptHdrLen = 0
	cf.Header.Flags = 0x2205 // Z8002 non-segmented, for now
}

func convSegType(segType byte) uint32 {
	switch segType {
	case binlib.XoutSegBSS:
		return binlib.CoffSectBSS
	case binlib.XoutSegCODE:
		return binlib.CoffSectTEXT
	case binlib.XoutSegDATA, binlib.XoutSegCONST:
		return binlib.CoffSectDATA
	default:
		return 0x00
	}
}

// ConvSectHdrs converts xout segment table, some members are set by Finalize()
func convSectHdrs(xf *binlib.XoutFile, cf *binlib.CoffFile) {
	sectPos := int32(binlib.CoffHdrLen + xf.Header.NumSegs*binlib.CoffSectHdrLen)
	offset := int32(0)
	var cfSect binlib.CoffSectHdr
	for _, seg := range xf.SegTbl {
		name := convSegName(seg.Type)
		copy(cfSect.Name[:], []byte(name)[0:8])
		cfSect.Vaddr = 0x00000000
		cfSect.Paddr = 0x00000000
		cfSect.Length = uint32(seg.Length)
		if seg.Type == binlib.XoutSegBSS {
			cfSect.Fpos = 0
		} else {
			cfSect.Fpos = sectPos + offset
			offset += int32(seg.Length)
		}
		cfSect.RelocTblFpos = 0 // Set by Finalize()
		cfSect.LineNumsFpos = 0
		cfSect.NumRelocs = 0 // Set by Finalize()
		cfSect.NumLines = 0
		cfSect.Flags = convSegType(seg.Type)
		cf.SectTbl = append(cf.SectTbl, cfSect)
	}
}

func convSymbIdx(xIdx uint16, xf *binlib.XoutFile, cf *binlib.CoffFile) uint32 {
	xname := xf.SymbTbl[xIdx].Name
	for idx, entry := range cf.SymbTbl {
		if symb, ok := entry.(binlib.CoffSymbEntry); ok {
			if symb.Name == xname {
				return uint32(idx)
			}
		}
	}
	return uint32(0xffffffff)
}

/*
func convSegTopSymbIdx(xIdx uint16, xf *binlib.XoutFile, cf *binlib.CoffFile) uint32 {
	segname := convSegName(xf.SegTbl[xIdx].Type)
	for idx, entry := range cf.SymbTbl {
		if symb, ok := entry.(binlib.CoffSymbEntry); ok {
			if segname == binlib.ConvertName(symb.Name) {
				return uint32(idx)
			}
		}
	}
	return uint32(0xffffffff)
}*/

func convSegTopSymbIdx(xIdx uint16, xf *binlib.XoutFile, cf *binlib.CoffFile) uint32 {
	segname := fmt.Sprintf("SEG%d0000", xIdx)
	for idx, entry := range cf.SymbTbl {
		if symb, ok := entry.(binlib.CoffSymbEntry); ok {
			if segname == binlib.ConvertName(symb.Name) {
				return uint32(idx)
			}
		}
	}
	return uint32(0xffffffff)
}

func convRelocTbl(xf *binlib.XoutFile, cf *binlib.CoffFile) {
	relocType := map[byte]uint16{0: 0xffff, 1: 0x0001, 2: 0xffff, 3: 0x0011,
		4: 0xffff, 5: 0x0001, 6: 0xffff, 7: 0x0011}
	// sort by Location
	sort.Slice(xf.RelocTbl, func(i, j int) bool { return xf.RelocTbl[i].Location < xf.RelocTbl[j].Location })
	// export to the coff reloc table
	for seg := 0; seg < int(xf.Header.NumSegs); seg++ {
		var cfReloc binlib.CoffRelocItem
		for _, xReloc := range xf.RelocTbl {
			if xReloc.SegIdx != byte(seg) {
				continue
			}
			cfReloc.Vaddr = uint32(xReloc.Location)
			cfReloc.Type = relocType[xReloc.Type]
			pos := calcAddr(xf, int(xReloc.SegIdx), xReloc.Location)
			cfReloc.Offset = uint32(xf.CodePart[pos])*256 + uint32(xf.CodePart[pos+1])
			cfReloc.Stuff = 0x5343
			if xReloc.Type == binlib.XoutRelocXOFF {
				cfReloc.SymbIdx = convSymbIdx(xReloc.SymbIdx, xf, cf)
			} else if xReloc.Type == binlib.XoutRelocOFF {
				cfReloc.SymbIdx = convSegTopSymbIdx(xReloc.SymbIdx, xf, cf)
			}
			cf.RelocTbl = append(cf.RelocTbl, cfReloc)
		}
	}
}

func convSymbTbl(xf *binlib.XoutFile, cf *binlib.CoffFile) {
	// Add dummy
	var dmySymb binlib.CoffSymbEntry
	copy(dmySymb.Name[:], []byte(".file"))
	dmySymb.Value = 0
	dmySymb.SectNo = -2
	dmySymb.Type = 0
	dmySymb.StrgClass = binlib.CoffSymbClassFile
	dmySymb.NumAux = 1
	cf.SymbTbl = append(cf.SymbTbl, dmySymb)
	var fdmySymb binlib.CoffSymbAuxFile
	copy(fdmySymb.Name[:], []byte("fake"))
	cf.SymbTbl = append(cf.SymbTbl, fdmySymb)

	// Convert local symbols
	var cfSymb binlib.CoffSymbEntry
	for _, symb := range xf.SymbTbl {
		if symb.SegIdx == 255 || symb.Type != binlib.XoutSymbLocal {
			continue
		}
		cfSymb.Name = symb.Name
		cfSymb.Value = uint32(symb.Value)
		cfSymb.SectNo = int16(symb.SegIdx + 1)
		cfSymb.Type = 0x00
		cfSymb.StrgClass = binlib.CoffSymbClassStatic
		cfSymb.NumAux = 0
		cf.SymbTbl = append(cf.SymbTbl, cfSymb)
	}
	// Convert Section symbols
	for _, symb := range xf.SymbTbl {
		if symb.Type != binlib.XoutSymbSeg {
			continue
		}
		cfSymb.Name = symb.Name
		cfSymb.Value = 0
		cfSymb.SectNo = int16(symb.SegIdx + 1)
		cfSymb.Type = 0x00
		cfSymb.StrgClass = binlib.CoffSymbClassStatic
		cfSymb.NumAux = 1
		cf.SymbTbl = append(cf.SymbTbl, cfSymb)

		var sectAuxSymb binlib.CoffSymbAuxSect
		sectAuxSymb.Length = uint32(xf.SegTbl[symb.SegIdx].Length)
		sectAuxSymb.NumLines = 0
		sectAuxSymb.NumRelocs = 0
		for _, reloc := range xf.RelocTbl {
			if reloc.SegIdx != symb.SegIdx {
				continue
			}
			sectAuxSymb.NumRelocs++
		}
		cf.SymbTbl = append(cf.SymbTbl, sectAuxSymb)
	}
	// Convert global symbols
	for seg := 0; seg < int(xf.Header.NumSegs); seg++ {
		for _, symb := range xf.SymbTbl {
			if symb.SegIdx == byte(seg) && symb.Type == binlib.XoutSymbGlobal {
				cfSymb.Name = symb.Name
				cfSymb.Value = uint32(symb.Value)
				cfSymb.SectNo = int16(symb.SegIdx + 1)
				cfSymb.Type = 0x00
				cfSymb.StrgClass = binlib.CoffSymbClassGlobal
				cfSymb.NumAux = 0
				cf.SymbTbl = append(cf.SymbTbl, cfSymb)
			}
		}
	}
	// Convert external symbols and constats
	for _, symb := range xf.SymbTbl {
		if symb.SegIdx != 255 {
			continue
		}
		switch symb.Type {
		case binlib.XoutSymbUndefEX:
			cfSymb.Name = symb.Name
			cfSymb.Value = 0
			cfSymb.SectNo = binlib.CoffSymbSCNExt
			cfSymb.Type = 0x00
			cfSymb.StrgClass = binlib.CoffSymbClassGlobal
			cfSymb.NumAux = 0
			cf.SymbTbl = append(cf.SymbTbl, cfSymb)
		case binlib.XoutSymbLocal:
			cfSymb.Name = symb.Name
			cfSymb.Value = uint32(symb.Value)
			cfSymb.SectNo = binlib.CoffSymbSCNAbs
			cfSymb.Type = 0x00
			cfSymb.StrgClass = binlib.CoffSymbClassGlobal
			cfSymb.NumAux = 0
			cf.SymbTbl = append(cf.SymbTbl, cfSymb)
		default:
		}
	}
}

func finalize(xf *binlib.XoutFile, cf *binlib.CoffFile) {
	// Set reloc table infomation in the section table
	relocFpos := binlib.CoffHdrLen + binlib.CoffSectHdrLen*
		len(cf.SectTbl) + len(xf.CodePart)
	count := 0
	for sect := 0; sect < len(cf.SectTbl); sect++ {
		relocFpos += count * binlib.CoffRelocItemLen
		cf.SectTbl[sect].RelocTblFpos = int32(relocFpos)
		count = 0
		for _, reloc := range xf.RelocTbl {
			if sect == int(reloc.SegIdx) {
				count++
			}
		}
		cf.SectTbl[sect].NumRelocs = uint16(count)
		if count == 0 {
			cf.SectTbl[sect].RelocTblFpos = 0
		}
	}
}
