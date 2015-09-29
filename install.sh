#!/bin/sh

INSTALL_PREFIX="/usr/local"
APROXY_VER="0.1"
APROXY_BIN_GO="./bin/aproxy.go"

if [ ! -f "$APROXY_BIN_GO" ]; then
	echo "please enter aproxy root dir to execute this script."
	exit 1
fi


eval $(go env)

GIT_SHA=`git rev-parse --short HEAD || echo "GitNotFound"`

val=$(go version)
gover=$(echo $val | awk -F ' ' '{print $3}')

outdir="aproxy-v$APROXY_VER-$gover-git$GIT_SHA"

echo "build file to ./dist/$outdir/"

rm -rf ./dist/$outdir/*
mkdir -p ./dist/$outdir/bin
mkdir -p ./dist/$outdir/conf

go build -o ./dist/$outdir/bin/aproxy ./bin/aproxy.go
go build -o ./dist/$outdir/bin/adduser ./bin/adduser.go

yes|cp -rf ./web ./dist/$outdir/
yes|cp -f ./conf/aproxy.toml ./dist/$outdir/conf/aproxy.toml.example

echo "install to $INSTALL_PREFIX/aproxy"

mkdir -p $INSTALL_PREFIX/aproxy

yes|cp -rf ./dist/$outdir/* $INSTALL_PREFIX/aproxy/

echo "build and install done."
