package common

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

func MakeConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Minute
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	return config
}

func GetConsumerGroup() (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(GetBootstrapServers(), GetGroupId(), MakeConfig())
}

func GetBootstrapServers() []string {
	servers := beego.AppConfig.Strings("bootstrapServers")
	for _, server := range servers {
		bootstrapServers := strings.Split(server, ",")
		return bootstrapServers
	}
	return nil
}
func GetGroupId() string {
	return beego.AppConfig.String("groupId")
}

func GetTopics() []string {
	return beego.AppConfig.Strings("topics")
}
func GetConsumerCount() (int, error) {
	return beego.AppConfig.Int("consumer_count")
}
