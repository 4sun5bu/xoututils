你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
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
Down load or clone xoututils. Move src/ to a directory that GOPATH points. In the directory directory type `go build xout2coff`, `go build xarch` and `go build xoutdump`. 

## To Build CP/M-8000 with GNU Binutils 
You need to convert cpmsys.rel and libcpm.a to buid CP/M-8000. I confirmed it possible to convert these two files in the **CP/M-8000 1.1** at **The Unofficial CP/M Web site**.  http://www.cpm.z80.de/download/cpm8k11.zip

To conver libcpm.a, there is a simple script in the xarch directory. The script extracts xout files from a library, converts them to COFF files and makes a library file. This script makes a lot of *.rel and *.o files so I recomend to do it in a working directory only for this. The generated library file has the same name as original xout library file, and original file is renamed to preserve. 
