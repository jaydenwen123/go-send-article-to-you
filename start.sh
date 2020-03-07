#!/bin/bash
go build && ./go-send-article-to-you > server.log &
echo "server is started...."
