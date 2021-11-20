#!/bin/bash

set -eu

app_name="wistia-dl"
build_folder="output"

rm -rf ${build_folder}
mkdir ${build_folder}

echo "Building..."
GOOS=darwin GOARCH=amd64 go build                    -trimpath -ldflags="-w -s" -o ${build_folder}/${app_name}_darwin_amd64 .
GOOS=linux GOARCH=amd64 go build                     -trimpath -ldflags="-w -s" -o ${build_folder}/${app_name}_linux_amd64 .
GOOS=linux GOARCH=386 go build                       -trimpath -ldflags="-w -s" -o ${build_folder}/${app_name}_linux_386 .
GOOS=windows GOARCH=amd64 go build                   -trimpath -ldflags="-w -s" -o ${build_folder}/${app_name}_windows_amd64.exe .
GOOS=windows GOARCH=386 go build                     -trimpath -ldflags="-w -s" -o ${build_folder}/${app_name}_windows_386.exe .

echo "Done"
