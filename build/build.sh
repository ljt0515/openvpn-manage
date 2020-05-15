#!/bin/bash

set -e

PKGFILE=openvpn-manage.tar.gz

cp -f ../$PKGFILE ./

docker build -t ljt0515/openvpn-manage .

rm -f $PKGFILE
