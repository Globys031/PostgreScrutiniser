#!/bin/bash

mkdir api/

NAME=resourceConfig
for NAME in $(echo "auth file resource-config")
do
	echo "\n\nGenerating openapi files for $NAME:\n"
	openapi-generator-cli generate -i http://192.168.56.102:9090/api/docs/${NAME} -g typescript-axios -o ./
	sed -i "s:'./configuration':'../configuration':" api.ts
	sed -i "s:'./base':'../base':" api.ts
	sed -i "s:'./common':'../common':" api.ts
	mv api.ts api/$NAME.ts
done

# Get rid of redundant files
rm index.ts git_push.sh .gitignore .npmignore .openapi-generator-ignore
rm -r .openapi-generator