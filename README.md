wsbroadcast
===========

Broadcast websocket stream

INSTALL
========= 
```
$ go build
```

alternatively set env variable GOPATH to existing folder
```
$ echo "export GOPATH=/home/user/go" >> ~/.bashrc
$ echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
```
reload bash
```
$ . ~/.bashrc
```
install with go tools
```
$ go get github.com/dtravin/wsbroadcast
```


RUN
=========
Example: 
```
$ wsbroadcast -i=ws://localhost:8084 -l=7777
```

TEST
=========
Dump some output
```
$ wsdump http://localhost:7777 > out.1
```

Check for handshake message present at first line
```
$ head -2 out.1 | grep -i jsmp
Binary file (standard input) matches
```

