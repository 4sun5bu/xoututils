/*
 *  xarch.go
 *
 *  Copyright (c) 2020 4sun5bu
 *  Released under the MIT license.
 *  See Licence.txt
 *
 *  A De-archiver of XOUT Lib
 *  XOUT files are extracted from lib file.
 */

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const ArHdrLen = 26
const ArFnameLen = 14
const ArMagic = 0xff65

type ArHdr struct {
	Name [ArFnameLen]byte
	Date uint32
	UID  byte
	GID  byte
	Mode uint16
	Size uint32
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("No input file")
	}
	infpath := os.Args[1]
	infile, err := os.Open(infpath)
	if err != nil {
		log.Fatalf("can not open %s\n", infpath)
	}

	/* check the Magic of the input file */
	var magic uint16
	infile.Seek(0, 0)
	binary.Read(infile, binary.BigEndian, &magic)
	if magic != ArMagic {
		log.Fatalln("not library file", magic)
	}

	var arhdr ArHdr
	for {
		/* read an ar header */
		err3 := binary.Read(infile, binary.BigEndian, &arhdr)
		if err3 != nil {
			break
		}
		if arhdr.Name[0] == 0 || arhdr.Size == 0 {
			break
		}

		/* get a file name */
		var objpath string
		for _, c := range arhdr.Name {
			if c == 0 {
				break
			}
			objpath += string(c)
		}
		fmt.Println(objpath)
		objfile, err4 := os.Create(objpath)
		if err4 != nil {
			log.Fatalln(err4)
		}
		/* read and write an object file */
		var obj []byte
		obj = make([]byte, arhdr.Size)
		err5 := binary.Read(infile, binary.BigEndian, &obj)
		if err5 != nil {
			log.Fatalln(err4)
		}
		err6 := binary.Write(objfile, binary.BigEndian, &obj)
		if err6 != nil {
			log.Fatalln(err4)
		}

		objfile.Close()
	}
}
