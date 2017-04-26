#!/bin/bash
#echo "Generating Executable"
#echo "OS Code": $1
export GOOS=$1
#echo "System Architecture": $2
export GOARCH=$2
cd $3
#echo "Building Executable"
go build -v $4
#echo "Executable generation complete!"