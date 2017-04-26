#!/bin/bash
#$1=docker name
#$2=tag
#$3=foldername
#$4=executable location (/var/home/user/executablename )
#keep basic docker file in specific location in this case i asue docker file store in /home/dockerfile/Dockerfile .github also ok for this 
#$5:$6 =ports(3000:3000)
#$7= prosess name 
#$8= ram
#$9= cpu

newport=0
path=$(pwd)
cd PublishedDockers/
sudo mkdir $3
cd $3
sudo cp $4 $path/PublishedDockers/$3/
sudo cp /var/www/html/engine/DockerRelatedFiles/Dockerfile $path/PublishedDockers/$3/
sudo docker build -t $1:$2 .
echo "docker run -d --memory=$8 --cpus=$9 -p $5:$6 $1:$2 $7"
sudo docker run -d --memory=$8 --cpus=$9 -p $5:$6 $1:$2 $7

newport=$(($5 + 1))

cd /var/www/html/engine/
sudo ./runner 10.240.0.6 $1 $newport $1_proxy

echo "Publishing to Docker is complete!"
