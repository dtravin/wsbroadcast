wsbroadcast
===========

Broadcast websocket stream

INSTALL
========= 
$ go build

alternatevly
$ go get http://github.com/dtravin/wsbroadcast


RUN
=========
Examle: 
$ wsbroadcast -i=ws://localhost:8084 -l=7777

TEST
=========
Dump some output
$ wsdump http://localhost:7777 > out.1

Check for handshake message present at first line
$ head -2 o2 | grep -i jsmp
Binary file (standard input) matches

