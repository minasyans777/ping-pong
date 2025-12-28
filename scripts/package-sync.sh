#!/bin/bash

set -e 

OPTIONS=$(getopt -o r:s:d:m:v:R:k:I:h --long arch:,help -- "$@")
eval set -- "$OPTIONS"

REPO=""
ARCH=""
PATH_NEW_PACKAGE=""
PATH_REPO_PACKAGE=""
PATH_METADATA_PACKAGE=""
VERSION=""
PATH_RELEASE=""
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
			echo " -R                Specify where the Release file should be generated"
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
	echo '#!/bin/sh
	set -e

	do_hash() {
	    HASH_NAME=$1
	    HASH_CMD=$2
	    echo "${HASH_NAME}:"
	    for f in $(find -type f); do
	        f=$(echo $f | cut -c3-) # remove ./ prefix
	        if [ "$f" = "Release" ]; then
	            continue
	        fi
	        echo " $(${HASH_CMD} ${f}  | cut -d" " -f1) $(wc -c $f)"
	    done
	}' >> ~/generate-release.sh

	cat << EOF >> ~/generate-release.sh
	cat << EOT
	Origin: Example Repository
	Label: Example
	Suite: stable
	Codename: stable
	Version: $VERSION
	Architectures: $ARCH
	Components: main
	Description: Game
	Date: $(date -Ru)
EOT
EOF
	echo '
	do_hash "MD5Sum" "md5sum"
	do_hash "SHA1" "sha1sum"
	do_hash "SHA256" "sha256sum" 
	'>> ~/generate-release.sh && chmod +x ~/generate-release.sh

	~/generate-release.sh | sudo tee "$PATH_RELEASE"
	
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



















