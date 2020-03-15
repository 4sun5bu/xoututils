/*
 *  xoutdump.go
 *
 *  Copyright (c) 2020 4sun5bu
 *  Released under the MIT license.
 *  See LICENSE.
 *
 *  Dump a XOUT file information
 */

package main

import (
	"fmt"
	"log"
	"os"

	"binlib"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("No input file")
	}
	infpath := os.Args[1]
	infile, err := os.Open(infpath)
	if err != nil {
		log.Fatalf("can not open %s\n", infpath)
	}
	defer infile.Close()

	xf := binlib.XoutFile{}
	err = xf.Read(infile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	fmt.Println("File =", os.Args[1])
	fmt.Printf("  Magic = 0x%4x\n", xf.Header.Magic)
	fmt.Printf("  nSegs = %d\n", xf.Header.NumSegs)
	fmt.Printf("  SegInfo    FilePos = 0x%04x\n", binlib.XoutHdrLen)
	fmt.Printf("  Code       FilePos = 0x%04x  Size = %d\n", xf.CodePos, xf.Header.CodePartLen)
	fmt.Printf("  RelocTable FilePos = 0x%04x  Size = %d\n", xf.RelocTblPos, xf.Header.RelocsLen)
	fmt.Printf("  SymbTable  FilePos = 0x%04x  Size = %d\n", xf.SymbTblPos, xf.Header.SymbsLen)
	fmt.Println()

	printSegInfo(xf)
	printRelocs(xf)
	printSymbs(xf)
}

func printSegInfo(xf binlib.XoutFile) {
	fmt.Println("Segment Info")
	for idx, seg := range xf.SegTbl {
		fmt.Printf(" %4d : No. = %1d, Type = %d, Size = %5d\n",
			idx, seg.Number, seg.Type, seg.Length)
	}
	fmt.Println()
}

func printRelocs(xf binlib.XoutFile) {
	fmt.Println("Relocation items")
	for idx, reloc := range xf.RelocTbl {
		fmt.Printf(" %4d : Seg = %3d, Type = %1d, Offset = 0x%04x, Symb = %d\n",
			idx, reloc.SegIdx, reloc.Type, reloc.Location, reloc.SymbIdx)
	}
	fmt.Println()
}

func printSymbs(xf binlib.XoutFile) {
	fmt.Println("Symbol table")
	for idx, symb := range xf.SymbTbl {
		fmt.Printf(" %4d : Seg = %3d, Type = %1d,  Val = 0x%04x, Name = %-8s \n",
			idx, symb.SegIdx, symb.Type, symb.Value, binlib.ConvertName(symb.Name))
	}
	fmt.Println()
}
