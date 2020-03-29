# Xoututils
 File convert tools for CP/M-8000 XOUT format to Z8k-COFF format.
 The XOUT format is linkable and executable format used in CP/M-8000.
 These tools are made for cross develping CP/M-8000 on Linux and Windows.
 Converted COFF files can be linked by the Gnu ld linker.

## License
 This software is released under the MIT License, see LICENSE.

## Commands
- **xout2coff** converts XOUT to Z8k-COFF.
- **xarch** extracts XOUT files from a libray.  
- **xoutdump** shows information about file structure, relocations and symbols.  

These commands take one filename, such as `xout2coff xxx.rel`.  

## How to Build
Set the GOPATH at the down loaded or cloned directory. In the directory type `go build xout2coff`, `go build xarch` and `go build xoutdump`. 
