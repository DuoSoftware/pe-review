#!/bin/bash
service nginx start
consul-template -consul=$CONSUL_URL -template="/templates/defaults.ctmpl:/etc/nginx/conf.d/default.conf:service nginx reload"