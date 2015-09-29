#!/bin/sh

APROXY_VER="0.1"

if [ ! -f "$APROXY_BIN_GO" ]; then
	echo "please enter aproxy root dir."
	exit 1
fi

eval $(go env)

GIT_SHA=`git rev-parse --short HEAD || echo "GitNotFound"`
val=$(go version)
gover=$(echo $val | awk -F ' ' '{print $3}')

outdir="aproxy-v$APROXY_VER-$gover-git$GIT_SHA"

all=('linux 386' 'linux amd64' 'darwin 386' 'darwin amd64' 'windows 386' 'windows amd64')

for  i  in "${all[@]}" ; do 
	b=($i)
	os=${b[0]}
	bit=${b[1]}
	outdir="aproxy-v$APROXY_VER-$os-$bit-$gover-git$GIT_SHA"
	echo "start build [$os-$bit] to ./release/$outdir"
	rm -rf ./release/$outdir/*
	mkdir -p ./release/$outdir/bin
	mkdir -p ./release/$outdir/conf

	if [ "$os"x = "windows"x ]; then
		GOOS=$os GOARCH=$bit go build -o ./release/$outdir/bin/aproxy.exe ./bin/aproxy.go
		GOOS=$os GOARCH=$bit go build -o ./release/$outdir/bin/adduser.exe ./bin/adduser.go
	else
		GOOS=$os GOARCH=$bit go build -o ./release/$outdir/bin/aproxy ./bin/aproxy.go
		GOOS=$os GOARCH=$bit go build -o ./release/$outdir/bin/adduser ./bin/adduser.go
	fi

	yes|cp -rf ./web ./release/$outdir/
	yes|cp -f ./conf/aproxy.toml ./release/$outdir/conf/aproxy.toml.example
done
