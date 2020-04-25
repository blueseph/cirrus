#!/bin/bash

set -e

NAME=cirrus
OUTPUT=dist

declare -A PLATFORMS=([linux]=linux [darwin]=osx [windows]=windows)
declare -A ARCHITECTURES=([386]=i386 [amd64]=amd64)

for platform in ${!PLATFORMS[@]}; do
    for architecture in ${!ARCHITECTURES[@]}; do
            full_name="${NAME}-${VERSION}_${PLATFORMS[$platform]}-${ARCHITECTURES[$architecture]}"
            bin_name="$NAME"

            if [ "$platform" == "windows" ]; then
                bin_name="${NAME}.exe"
            fi

            mkdir -p "$OUTPUT/$full_name"

            GOOS=$platform GOARCH=$architecture go build -o "$OUTPUT/${full_name}/${bin_name}"

            zip -9 -r "$OUTPUT/${full_name}.zip" "$OUTPUT/$full_name"

            rm -r "$OUTPUT/$full_name"
    done
done