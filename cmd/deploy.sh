#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o knowledge .

ssh -i greypan.pem ubuntu@ec2-18-162-229-84.ap-east-1.compute.amazonaws.com "rm -f ubuntu:ubuntu /home/ubuntu/configure/bk/knowledge"

scp -i greypan.pem /Users/bopan/Code/backend/knowledge-service/knowledge ubuntu@ec2-18-162-229-84.ap-east-1.compute.amazonaws.com:/home/ubuntu/configure/bk/

border="====================================="
echo "$border"
echo "               部署完成"
echo "$border"