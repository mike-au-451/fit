#!/bin/bash

# Quick and nasty make

FILES="$@"
if [[ -z "$FILES" ]]
then
	FILES="*"
fi

for FILE in cmd/$FILES
do
	if [[ "$FILE" == *"go" && -e "$FILE" ]]
	then
		BASE=$(basename "$FILE" .go)
		go build -o bin/$BASE $FILE
		if [[ $? != 0 ]]
		then
			break
		fi
	fi
done
