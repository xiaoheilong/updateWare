@ echo off
%1 %2
ver|find "5.">nul&&goto :Admin
mshta vbscript:createobject("shell.application").shellexecute("%~s0","goto :Admin","","runas",1)(window.close)&goto :eof
:Admin
cd /d %~dp0
go env -w GO111MODULE=auto
mkdir  api
mkdir  package
mkdir  src
mkdir  pkg
mkdir  bin
echo $GOPATH
./package
cd ./api
go mod init api