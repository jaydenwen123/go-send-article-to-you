#!/bin/bash
echo  > server.log
go build && ./go-send-article-to-you > server.log &
echo "server is started...."
tail -f server.log
