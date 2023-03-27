#!/bin/bash

echo "This is a script for calculating how many lines of code this repo has"
echo "Script should take in at least one file extension, for example 'go'"

EXTENSIONS=$@

echo "Extensions: $EXTENSIONS"

TOTAL_LINES=0
for EXTENSION in $EXTENSIONS
do
        # Don't calculate files under node_modules or test code related files
        for FILE in $(find . -not -path "*node_modules*" -not -path "*cypress*" -not -path "*__tests__*" -iname "*\.${EXTENSION}");
        do
                TOTAL_LINES=$(( $TOTAL_LINES + $(wc -l $FILE | awk '{ print $1 }') ))
                echo "File: $FILE, Current line count: $TOTAL_LINES"
        done
done

echo "Final calculation: ${TOTAL_LINES}";
