#!/bin/sh
set -x
set -e

# Set temp environment vars
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin
export BUILDPATH=${GOPATH}/src/github.com/vkorehov/go-git2consul
export PKG_CONFIG_PATH="/usr/lib/pkgconfig/:/usr/local/lib/pkgconfig/"
export PKG_CONFIG_PATH="${PKG_CONFIG_PATH}:/tmp/libgit2/install/lib/pkgconfig:/tmp/openssl/install/lib/pkgconfig:/tmp/libssh2/build/src"

# Install libraries
/build/install-openssl.sh
/build/install-libssh2.sh
/build/install-libgit2.sh

# Set up go environment
mkdir -p $(dirname ${BUILDPATH})
ln -s /app ${BUILDPATH}
