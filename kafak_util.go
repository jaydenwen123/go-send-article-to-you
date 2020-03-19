package main

//初始化kafka的消费者和生产者

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

//消费者
var consumer *kafka.Reader

//生产者
var producter *kafka.Writer

//话题
//initKafkaConsumer 初始化kafka消费者
func initKafkaConsumer(brokers []string, groupId string, topic string) {
	consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		GroupID:   groupId,
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
}

//initKafkaProducter 初始化kafka生产者
func initKafkaProducter(brokers []string, topic string, async bool) {
	producter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		//默认轮询
		Balancer: nil,
		Async:    async,
	})
}

//createKafkaTopic 创建主题
func createKafkaTopic(network string, address string,
	topic string, numPartitions int, replicationsFactor int, ) {
	conn, err := kafka.Dial(network, address)
	if err != nil {
		logs.Error("the kafka DialContext error:%s", err.Error())
		panic(err)
	}
	defer conn.Close()
	//如果话题存在，删除话题
	err = conn.DeleteTopics(topic)
	if err != nil {
		logs.Error("DeleteTopics occurs error:%v", err)
	}

	//创建话题
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:              topic,
		NumPartitions:      numPartitions,
		ReplicationFactor:  replicationsFactor,
		ReplicaAssignments: nil,
		ConfigEntries:      nil,
	})
	if err != nil {
		logs.Error("kafka CreateTopics occurs error:%s", err.Error())
		panic(err)
	}

}

//sendMessage 发送消息
func sendMessage(data interface{}) error {
	//生产数据
	bdata, err := jsoniter.Marshal(data)
	if err != nil {
		logs.Error("the message  jsoniter.Marshal occurs error:%s", err.Error())
		return err
	}
	err = producter.WriteMessages(ctx, kafka.Message{
		Value: bdata,
	})
	if err != nil {
		logs.Error("producter.WriteMessages write the mssage error:%v", err)
		return err
	}
	return nil
}
