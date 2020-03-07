#!/bin/bash
psinfo=`ps -ef |grep  ./go-send-article-to-you |grep -v grep`
if [ `echo $psinfo |wc -l` == 1 ]; then
   echo "the ps monitor success..."
   echo $psinfo
   pid=`echo $psinfo | awk -F " " '{print $2}'`
   echo "the server pid:<"$pid">"
   # kill pid
  kill -9 $pid
  echo "stop server finish."
else
  echo "current there is no server is running...."
fi
