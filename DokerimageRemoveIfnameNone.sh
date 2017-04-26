#!/bin/bash
#will delete all docker images name "<none>"

docker images  | awk '{ print $1; }' >names.txt
docker images  | awk '{ print $3; }' >ids.txt
name=( $(cut -d ' ' -f 1 "./names.txt") )
ids=( $(cut -d ' ' -f 1 "./ids.txt") )
size=${#name[@]}
for (( i=1; i <= $size; ++i ))
do
    if [ "${name[$i]}" == "<none>" ]; tihen
  	echo "<none> image found "
  	docker rmi -f ${ids[$i]}
	fi
done
rm names.txt ids.txt
