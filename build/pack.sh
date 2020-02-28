#!/bin/bash

set -e

time docker run \
    -v "$PWD/../":/go/src/openvpn-manage \
    --rm \
    -w /usr/src/myapp \
    awalach/beego:1.8.1 \
    sh -c "cd /go/src/openvpn-manage/ && bee version && bee pack -exr='^vendor|^data.db|^build|^README.md|^docs|LICENSE|README.en.md|.gitignore|go.mod|go.sum'"
