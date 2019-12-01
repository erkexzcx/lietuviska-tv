#!/usr/bin/env bash

# Make sure we are running this script from project's directory
if ! [ -d "src" ]; then
	echo "Please run this script from project's directory! Exitting..."
	exit 1
fi

# Prepare output directory:
mkdir -p dist

# Build Linux binaries:
env GOOS=linux GOARCH=386 go build -o "dist/lietuviskatv_linux_i386" src/*.go             # Linux i386
env GOOS=linux GOARCH=amd64 go build -o "dist/lietuviskatv_linux_x86_64" src/*.go         # Linux 64bit

# Build Linux ARM binaries:
env GOOS=freebsd GOARCH=amd64 go build -o "dist/lietuviskatv_freebsd_x86_64" src/*.go     # FreeBSD 64bit

# Build freebsd binary:
env GOOS=linux GOARCH=arm GOARM=5 go build -o "dist/lietuviskatv_linux_arm" src/*.go      # Linux armv5/armel/arm (it also works on armv6)
env GOOS=linux GOARCH=arm GOARM=7 go build -o "dist/lietuviskatv_linux_armhf" src/*.go    # Linux armv7/armhf
env GOOS=linux GOARCH=arm64 go build -o "dist/lietuviskatv_linux_aarch64" src/*.go        # Linux armv8/aarch64

# Build Windows binaries (oh God whyyy):
env GOOS=windows GOARCH=386 go build -o "dist/lietuviskatv_windows_i386.exe" src/*.go     # Windows 32bit
env GOOS=windows GOARCH=amd64 go build -o "dist/lietuviskatv_windows_x86_64.exe" src/*.go # Windows 64bit
