#!/bin/bash
if [ $1 -z "" ]; then
  echo "the commit info is empty..."
    exit
fi
echo "begin to commit to repositry..."
git commit -a -m $1
echo "commit success...."
echo "begin to push to remote repositry..."
git push origin master
echo "git push success...."
