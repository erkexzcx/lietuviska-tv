#!/usr/bin/env bash

# Make sure we are running this script from project's directory
if ! [ -d "src" ]; then
	echo "Please run this script from project's directory! Exitting..."
	exit 1
fi

package_name="livetv"

platforms=(
	"linux/386"
	"linux/amd64"
	"linux/arm"
	"linux/arm64"
	"linux/mips"
	"linux/mips64"
	"linux/mips64le"
	"linux/mipsle"
	"linux/ppc64"
	"linux/ppc64le"
	"linux/s390x"
	"windows/386"
	"windows/amd64"
)

for platform in "${platforms[@]}"; do

    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o dist/$output_name src/*.go
    if [ $? -ne 0 ]; then
        echo 'Error occured during GO build! Exitting...'
        exit 1
    fi

done

echo 'Building completed!'
