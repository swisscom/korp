#!/bin/sh

SRC_PATTERN="(.*)image: (.*)"
REPL_PATTERN="\1image: eas-docker-virtual.artifactory.swisscom.com\/\2"
BAK_SUFFIX="_bak"
DIR="$1"

if [[ "$DIR" == "" ]]; then
	echo
	echo "/!\ No directory specified, parsing current one /!\\"
	DIR="."
fi

echo
echo "Parse $DIR to add Swisscom Artifactory docker prefix"

for filepath in $(find $DIR -type f -regex ".*.yaml"); do
	# if [[ $filepath != *"DS_Store"* ]]; then
		# filename=$(basename $filepath)
		# echo "    ... parsing $filename"
		echo "    ... parsing $filepath"
		sed -i "$BAK_SUFFIX" -E "s/$SRC_PATTERN/$REPL_PATTERN/g" "$filepath"
	# fi
done

echo
echo "Parse $DIR completed"
echo
