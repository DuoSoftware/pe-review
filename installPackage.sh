#!/bin/bash
echo "Installing Package": $1
go get $1
echo "Installation Complete!"