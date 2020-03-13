# xoututils
 File convert tools for CP/M-8000 XOUT format object files.
 The XOUT format is linkable and executable format used in CP/M-8000.

## License
 This software is released under the MIT License, see LICENSE.

## xout2coff
 A converter from XOUT to Z8k-COFF written in Go.
 Converted COFF files can be linked by the Gnu ld linker in binutils.
 
  xout2coff xxx.rel
 
## xarch
 A de-archiver for library file.
 
  xarch libxxx.a
