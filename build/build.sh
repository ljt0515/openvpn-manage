#!/bin/bash

set -e

PKGFILE=openvpn-manage.tar.gz

cp -f ../$PKGFILE ./

docker build -t awalach/openvpn-manage .

rm -f $PKGFILE
