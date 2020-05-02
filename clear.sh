#!/bin/bash
echo "======begin clear the kafka data..."
rm -rf /usr/local/var/lib/kafka-logs/*
rm -rf /usr/local/var/lib/kafka-logs-1/*
rm -rf /usr/local/var/lib/kafka-logs-2/*
rm -rf /usr/local/var/lib/zookeeper/*
echo "=======finish clear the kafa data..."

echo "=======begin restart zookeeper and kafka...."
pids=`jps |grep -E "QuorumPeerMain|(Kafka$)"|awk '{print $1}'|xargs`
echo "=======zookeeper and kafka pids:"$pids
kill -9 $pids
echo "=======finish stop zookeeper and kafka...."
bash /Users/wenxiaofei/learn-kafka/start_kafka_clusters.sh
echo "=======finish restart zookeeper and kafka...."
echo "=======begin restart the program <go-send-article-to-you>"
./restart.sh
