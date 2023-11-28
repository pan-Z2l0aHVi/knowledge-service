#
## Build
```sh
GOOS=linux GOARCH=amd64 go build -o knowledge .

scp -r -i pan.pem /Users/bopan/Code/backend/knowledge-service/knowledge ubuntu@ec2-16-163-30-187.ap-east-1.compute.amazonaws.com:/home/ubuntu/configure/bk/
```