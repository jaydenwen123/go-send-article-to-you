#!/bin/bash
if [ $# == 0 ]; then
  echo "the commit info is empty..."
  exit
fi
echo "begin to commit to repositry..."
echo "commit script:<git commit -a -m $1>"
git commit -a -m $1
echo "commit success...."
echo "begin to push to remote repositry..."
git push origin master
echo "git push success...."
