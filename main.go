package main

import (
	"bigrisk/consumer"
	"context"
	"github.com/sirupsen/logrus"
	"gitlaball.nicetuan.net/wangjingnan/golib/gsr/log"
	"gitlaball.nicetuan.net/wangjingnan/golib/logrus-gsr/wrapper"
	"os"
	"os/signal"
	"riskengine/common"
	"sync"
	"syscall"
)

//var logger log.Logger
func init() {
	logger := wrapper.NewLogger()
	logger.Logrus.SetLevel(logrus.ErrorLevel)
	log.SetLogger(logger)
}

func main() {
	consumerGroup, err := common.GetConsumerGroup()
	if err != nil {
		log.Fatal("Error creating consumer group: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	consumer := consumer.NewConsumer()

	doConsume := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			err := consumerGroup.Consume(ctx, common.GetTopics(), &consumer)
			if err != nil {
				log.Fatal("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}

	waitGroup := &sync.WaitGroup{}
	consumerCount, _ := common.GetConsumerCount()
	waitGroup.Add(consumerCount)

	for i := 0; i < consumerCount; i++ {
		go doConsume(waitGroup)
		log.Info("Consumer goroutine %d is up and running", i)
	}

	<-consumer.Ready
	log.Info("Consumer group is up and running")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigterm:
		log.Info("Consuming terminated by signal")
	case <-ctx.Done():
		log.Info("Consuming terminated by context")
	}

	cancel()
	waitGroup.Wait()

	err = consumerGroup.Close()
	if err != nil {
		log.Fatal("Error closing consumer group: %v", err)
	}
}
