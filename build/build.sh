#!/bin/bash

set -e

PKGFILE=openvpn-manage.tar.gz

cp -f ../$PKGFILE ./

docker build -t 913519/openvpn-manage .

rm -f $PKGFILE
