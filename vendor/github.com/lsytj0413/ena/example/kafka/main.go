package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
)

var urls = []string{"127.0.0.1:32776", "127.0.0.1:32777", "127.0.0.1:32778"}

var wait sync.WaitGroup
var count = 10
var topic = "kafka_test_p3"

func producer() {
	fmt.Println("producer start")
	// wait.Add(1)
	defer wait.Done()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	// config.Producer.Partitions =
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_0

	producer, err := sarama.NewAsyncProducer(urls, config)
	if err != nil {
		fmt.Println("producer: ", err)
		return
	}
	defer producer.AsyncClose()

	fmt.Println("producer client get ok")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < count; i++ {
			uid := uuid.NewV4()
			value := uid.String() + fmt.Sprintf("%d", i)
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(value),
				Value: sarama.ByteEncoder(value),
			}
			producer.Input() <- msg
			fmt.Println("producer: ", value)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < count; i++ {
			select {
			case suc := <-producer.Successes():
				fmt.Printf("producer: %+v\n", suc)
			case fail := <-producer.Errors():
				fmt.Println("producer: ", fail.Err)
			}
		}
	}()

	wg.Wait()
	fmt.Println("procuder done")
}

func consumer() {
	fmt.Println("consumer start")
	// wait.Add(1)
	defer wait.Done()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_0

	consumer, err := sarama.NewConsumer(urls, config)
	if err != nil {
		fmt.Println("consumer: ", err)
		return
	}
	defer consumer.Close()
	fmt.Println("consumer client get ok")

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("consumer: ", err)
		return
	}
	// fmt.Println("partitions: %v", partitions)
	// partitions := []int32{0}

	consumers := make([]sarama.PartitionConsumer, len(partitions))
	for i := 0; i < len(consumers); i++ {
		c, err := consumer.ConsumePartition(topic, partitions[i], sarama.OffsetNewest)
		if err != nil {
			fmt.Println("consumer partition: ", err)
			return
		}
		defer c.Close()

		consumers[i] = c
	}
	var i int
	for {
		for j := 0; j < len(partitions); j++ {
			select {
			case msg := <-consumers[j].Messages():
				i++
				fmt.Printf("consumer: %+v\n", msg)
				if i >= count {
					goto done
				}
			case <-consumers[j].Errors():
			case <-time.After(time.Second):
			}
		}
	}
done:

	fmt.Println("consumer done")
}

func main() {
	wait.Add(2)
	go producer()
	go consumer()
	wait.Wait()
	fmt.Println("main")
}
