#!/bin/bash

export PATH=$PATH:/var/www/html/engine;
export GOPATH=/var/www/html/engine;
export PATH=$PATH:$GOPATH/bin;

source /etc/profile


path=$(pwd)
if [ ! -d "client" ]; then
	mkdir client
fi
# if [ ! -d "executable" ]; then
# 	mkdir executable
# fi
cd client
if [ ! -d "windows" ]; then
	mkdir windows
fi
if [ ! -d "linux" ]; then
	mkdir linux
fi
if [ ! -d "mac" ]; then
	mkdir mac
fi
if [ ! -d "raspberry" ]; then
	mkdir raspberry
fi
cd windows
if [ ! -d "32" ]; then
	mkdir 32
fi
if [ ! -d "64" ]; then
	mkdir 64
fi
cd ..
cd linux
if [ ! -d "32" ]; then
	mkdir 32
fi
if [ ! -d "64" ]; then
	mkdir 64
fi
cd ..
cd mac
if [ ! -d "32" ]; then
	mkdir 32
fi
if [ ! -d "64" ]; then
	mkdir 64
fi
cd $path

#windows 32
export GOOS=windows
export GOARCH=386
go build -v processengine/client/RaspberryPI_Client_Main.go

mv RaspberryPI_Client_Main.exe client/windows/32/smoothflow_w32.exe

#windows 64

export GOOS=windows
export GOARCH=amd64
go build -v processengine/client/RaspberryPI_Client_Main.go

mv RaspberryPI_Client_Main.exe client/windows/64/smoothflow_w64.exe

#linux 32
export GOOS=linux
export GOARCH=386
go build -v processengine/client/RaspberryPI_Client_Main.go

mv RaspberryPI_Client_Main client/linux/32/smoothflow_li32

#linux 64

export GOOS=linux
export GOARCH=amd64
go build -v processengine/client/RaspberryPI_Client_Main.go 

mv RaspberryPI_Client_Main client/linux/64/smoothflow_li64

#mac 64
export GOOS=darwin
export GOARCH=amd64
go build -v processengine/client/RaspberryPI_Client_Main.go 

mv RaspberryPI_Client_Main client/mac/64/smoothflow_mac64

#mac 32
export GOOS=darwin
export GOARCH=386
go build -v processengine/client/RaspberryPI_Client_Main.go

mv RaspberryPI_Client_Main client/mac/32/smoothflow_mac32

#raspberry

export GOOS=linux
export GOARCH=amd64
go build -v processengine/client/RaspberryPI_Client_Main.go 

mv RaspberryPI_Client_Main client/raspberry/smoothflow_raspberry

#remove download folder before replaceing
rm -r /var/www/html/download

cd $path
mv client /var/www/html/download
