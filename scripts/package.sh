#!/bin/bash

set -e

OPTIONS=$(getopt -o v:s:n:h --long deb,rpm,help -- "$@")
eval set -- "$OPTIONS"

DIR=""
FILE=""
PACKAGE_TYPE=""
PACKAGE_VERSION=""
while true;do
	case "$1" in 
		-v)
			PACKAGE_VERSION=$2
			shift 2
			;;
		-n)
			DIR=$2
			shift 2
			;;
		-s)
			FILE=$2
			shift 2
			;;
		-h)
			echo "Options:"
			echo " -n         Specify the directory name for the package"
			echo " -s         Specify the full path of the binary file"
			echo " -h         Show this help message"
			echo " -v         Specify the package version"
			echo " --deb      For creating a .deb package"
			echo " --rpm      For creating a .rpm package"
			exit 0
			;;
		--deb)
			PACKAGE_TYPE="deb"
			shift
			;;
		--rpm)
			PACKAGE_TYPE="rpm"
			shift
			;;
		--)
			shift
			break
			;;

		*) 
			echo "Invalid option: $1"
			echo "Try '--help' or '-h' for more information."
			exit 1 
			;;
	esac
done

if [[ -z $PACKAGE_TYPE ]]; then
    echo "Error: You must specify either --deb or --rpm"
    echo "Hint: --deb or --rpm option is required"
    echo "Usage: $0 -n <directory_name> -s <full_path_binary> -v <package_version> [--deb|--rpm]"
    exit 1
fi

if [[ -z $DIR ]]; then
    echo "Error: -n option is required"
    exit 1
fi

if [[ -z $FILE ]]; then
    echo "Error: -s option is required"
    exit 1
fi

if [[ -z $PACKAGE_VERSION ]]; then
    echo "Error: -v option is required"
    exit 1
fi

if [[ "$PACKAGE_TYPE" == "deb" ]]; then
	mkdir -p "$DIR/usr/bin"
	cp "$FILE" "$DIR/usr/bin/"
        mkdir -p "$DIR/DEBIAN"

	cat <<EOF > "$DIR/DEBIAN/control"
Package: ping-pong
Version: $PACKAGE_VERSION
Maintainer: example <example@example.com>
Depends: libc6
Architecture: amd64
Homepage: http://example.com
Description: Game
EOF

	dpkg --build "$DIR"
fi

DATE=$(date '+%a %b %d %Y')

if [[ "$PACKAGE_TYPE" == "rpm" ]];then
	mkdir -p "$DIR"

	cat <<EOF > "$DIR/ping-pong.spec"
Summary: Game
Name: ping-pong
Version: $PACKAGE_VERSION
License: example
Requires: bash
Release: 1
%description
ping-pong game

%install
mkdir -p %{buildroot}/usr/bin/
cp "$FILE" %{buildroot}/usr/bin/ping-pong

%files
/usr/bin/ping-pong

%changelog
* $DATE example <example@example.com>
- initial example
EOF

	rpmbuild --target "x86_64" -bb "$DIR/ping-pong.spec"
	cp ~/rpmbuild/RPMS/x86_64/*.rpm .	
fi


