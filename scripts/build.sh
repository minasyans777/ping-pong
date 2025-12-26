#!/bin/bash

FILE=""
OUTPUT=""

while getopts "s:o:h" opt; do
	case $opt in
		s)
			FILE=$OPTARG
			;;
		o)
			OUTPUT=$OPTARG
			;;
		h)
          		echo "Options:"
            		echo "  -s    Specify the full path of the program to build"
            		echo "  -o    Specify the directory where to save the built program"
            		echo "  -h    Show this help message"
            		exit 0
            		;;
		*)
			echo "Invalid option!"
      			echo "Usage: $0 -s full_path_to_program -o output_directory [-h]"
      			exit 1
      			;;
	esac
done

if [[ -z $FILE ]]; then
	echo "Error: -s option is required"
	echo "Usage: $0 -s full_path_to_program -o output_directory [-h]"
	exit 1
fi

if [[ -z $OUTPUT ]]; then
	OUTPUT="."
fi

go build -o "$OUTPUT" "$FILE"
