#!/bin/bash

set -e

read -p "Please write directory name: " dir_name

mkdir -p ~/$dir_name

mkdir -p ~/$dir_name/usr/bin

cp /home/sam/pingpong/main ~/$dir_name/usr/bin/.

mkdir -p ~/$dir_name/DEBIAN

echo "Package: ping-pong
Version: 0.2.0
Maintainer: example <example@example.com>
Depends: libc6
Architecture: amd64
Homepage: http://example.com
Description: Game" \
> ~/$dir_name/DEBIAN/control

dpkg --build ~/$dir_name




