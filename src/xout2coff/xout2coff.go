/*
 *  xout2coff.go
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
	"log"
	"os"
	"path/filepath"

	"binlib"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("No input file")
	}
	infpath := os.Args[1]
	infile, err := os.Open(infpath)
	if err != nil {
		log.Fatalf("can not open %s\n", err)
	}
	defer infile.Close()

	infname := filepath.Base(infpath)
	outfpath := infname[:len(infname)-len(filepath.Ext(infname))] + ".o"
	outfile, err1 := os.Create(outfpath)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "can not create %s \n", outfpath)
		return
	}
	defer outfile.Close()

	xf := binlib.XoutFile{}
	err = xf.Read(infile)
	cf := binlib.CoffFile{}
	cf.Open(outfile)

	// prepare
	assignBSS(&xf)
	//addLocalSymb(&xf)
	checkSegSymb(&xf)
	// convert
	convSectHdrs(&xf, &cf)
	cf.CodePart = &xf.CodePart
	convSymbTbl(&xf, &cf)
	convRelocTbl(&xf, &cf)
	finalize(&xf, &cf)
	convHdr(&cf)

	// CoffPrintHdr()
	// CoffPrintSectTbl()
	// CoffPrintRelocTbl()
	// CoffPrintSymbTbl()

	if err = cf.WriteHdr(); err != nil {
		goto errhandle
	}
	if err = cf.WriteSectTbl(); err != nil {
		goto errhandle
	}
	if err = cf.WriteCodePart(); err != nil {
		goto errhandle
	}
	if err = cf.WriteRelocTbl(); err != nil {
		goto errhandle
	}
	if err = cf.WriteSymbTbl(); err != nil {
		goto errhandle
	}
        return
errhandle:
	log.Fatalln(err)
}
