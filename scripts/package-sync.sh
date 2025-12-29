#!/bin/bash

set -e 

OPTIONS=$(getopt -o r:s:d:m:v:R:k:I:D:h --long arch:,help -- "$@")
eval set -- "$OPTIONS"

REPO=""
ARCH=""
PATH_NEW_PACKAGE=""
PATH_REPO_PACKAGE=""
PATH_METADATA_PACKAGE=""
VERSION=""
PATH_RELEASE=""
PATH_RELEASE_DIR=""
GPG_KEY_ID=""
PATH_INRELEASE=""

while true;do
	case "$1" in 
		-r)
			REPO=$2
			shift 2
			;;
		--arch)
			ARCH=$2
			shift 2
			;;
		-s)
			PATH_NEW_PACKAGE=$2
			shift 2
			;;
		-d)
			PATH_REPO_PACKAGE=$2
			shift 2
			;;
		-m)
			PATH_METADATA_PACKAGE=$2
			shift 2
			;;
		-v)
			VERSION=$2
			shift 2
			;;
		-R) 
			PATH_RELEASE=$2
			shift 2
			;;
		-D)
			PATH_RELEASE_DIR=$2
			shift 2
			;;
		-k)
			GPG_KEY_ID=$2
			shift 2
			;;
		-I)
			PATH_INRELEASE=$2
			shift 2
			;;
		-h|--help)
			echo "Options:"
			echo " -r                Specify the full path of the repository"
                        echo " --arch            Specify the architecture"
			echo " -s                Specify the full path of the new package"
			echo " -d                Specify the full paths of the packages inside the repository"
			echo " -m                Specify where the Packages file should be generated"
			echo " -v                Specify the package version"
			echo " -D                Specify directory where the Release file should be generated"
			echo " -R                Specify where the Release file is located "
			echo " -k                Specify the GPG key ID"
			echo " -I                Specify where the InRelease file should be generated"
                        echo " -h|--help         Show this help message"
			exit 0
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

if [[ -z $REPO ]];then
	echo "Error: -r option is required"
	exit 1
fi

if [[ -z $PATH_NEW_PACKAGE ]];then
        echo "Error: -s option is required"
        exit 1
fi

if [[ -z  $GPG_KEY_ID ]];then
        echo "Error: -k option is required"
        exit 1
fi

if [[ -z $PATH_REPO_PACKAGE ]];then
        echo "Error: -d option is required"
        exit 1
fi

if [[ "$PATH_NEW_PACKAGE" == *.deb ]];then

	if [[ -z $ARCH ]];then
        echo "Error: --arch option is required"
        exit 1
	fi

	if [[ -z $PATH_REPO_PACKAGE ]];then
        echo "Error: -d option is required"
        exit 1
	fi    	
	
	if [[ -z $PATH_METADATA_PACKAGE ]];then
        echo "Error: -m option is required"
        exit 1
	fi
	
	if [[ -z $VERSION ]];then
        echo "Error: -v option is required"
        exit 1
	fi

	if [[ -z  $PATH_RELEASE ]];then
        echo "Error: -R option is required"
        exit 1
	fi

	if [[ -z  $PATH_INRELEASE ]];then
        echo "Error: -I option is required"
        exit 1
	fi

	sudo cp "$PATH_NEW_PACKAGE" "$PATH_REPO_PACKAGE"

	cd "$REPO" || exit 1

	sudo dpkg-scanpackages --arch "$ARCH" pool/ | sudo tee "$PATH_METADATA_PACKAGE"

	cat "$PATH_METADATA_PACKAGE" | gzip -9 |sudo tee "$PATH_METADATA_PACKAGE.gz"

	cat << EOF | sudo tee  "$PATH_RELEASE"
Origin: Example Repository
Label: Example
Suite: stable
Codename: stable
Version: $VERSION
Architectures: $ARCH
Components: main
Description: Game
Date: $(date -Ru)
EOF

	apt-ftparchive release "$PATH_RELEASE_DIR" | sudo tee -a "$PATH_RELEASE"

	cat "$PATH_RELEASE" | gpg --default-key "$GPG_KEY_ID" -abs | sudo tee "$PATH_RELEASE.gpg"

	cat "$PATH_RELEASE" | gpg --default-key "$GPG_KEY_ID" -abs --clearsign | sudo tee "$PATH_INRELEASE"
fi

if [[ "$PATH_NEW_PACKAGE" == *.rpm  ]];then
	
	echo "%_signature gpg
%_openpgp_sign_id $GPG_KEY_ID" | sudo tee ~/.rpmmacros

	cp "$PATH_NEW_PACKAGE" "$PATH_REPO_PACKAGE"

	rpm --addsign "$PATH_NEW_PACKAGE"

	cd "$REPO" || exit 1

	createrepo .

	gpg --detach-sign --armor repodata/repomd.xml
fi






























