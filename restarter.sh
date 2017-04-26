#!/bin/bash
echo "restart sh running "
#screen -m -d smoothflow
cd /var/www/html/serve/
./serve.sh
echo $?