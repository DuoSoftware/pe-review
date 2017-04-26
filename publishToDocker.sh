#!/bin/bash
echo "Publishing to Docker process started."
echo ""
echo "Original File Location ":$1
echo "Location in docker ":$2
echo "Executable Name ":$3
echo "Port ":$4
echo "RAM ":$5
echo "CPU ":$6

#./smoothscipt.sh /home/pamidu/smooth/ /home/smooth/ ObjectStore 3000
#./smoothscipt.sh originalFileLocation whereToPutFileInDocker executableName Port
#echo "docker run -d -v $1:$2 -p $4:$4 $3:latest $2$3"

#docker run -d -v $1:$2 -p $4:$4 $3:latest $2/$3
echo "./smooth.sh $3 latest smoothflow $1$3 $4 $4 /$3 $5"

./smooth.sh $3 latest smoothflow $1$3 $4 $4 /$3 $5 $6

# delete any docker without a proper Name
chmod u+x ./DokerimageRemoveIfnameNone.sh
./DokerimageRemoveIfnameNone.sh
exit 0