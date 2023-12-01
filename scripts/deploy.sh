#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o knowledge .

ssh -i pan.pem ubuntu@ec2-16-163-30-187.ap-east-1.compute.amazonaws.com "rm -f ubuntu:ubuntu /home/ubuntu/configure/bk/knowledge"

scp -i pan.pem /Users/bopan/Code/backend/knowledge-service/knowledge ubuntu@ec2-16-163-30-187.ap-east-1.compute.amazonaws.com:/home/ubuntu/configure/bk/

border="====================================="
echo "$border"
echo "               部署完成"
echo "$border"