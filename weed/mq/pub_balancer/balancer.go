package pub_balancer

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/gateway-dao/seaweedfs/weed/mq/topic"
	"github.com/gateway-dao/seaweedfs/weed/pb/mq_pb"
)

const (
	MaxPartitionCount  = 8 * 9 * 5 * 7 //2520
	LockBrokerBalancer = "broker_balancer"
)

// Balancer collects stats from all brokers.
//
//	When publishers wants to create topics, it picks brokers to assign the topic partitions.
//	When consumers wants to subscribe topics, it tells which brokers are serving the topic partitions.
//
// When a partition needs to be split or merged, or a partition needs to be moved to another broker,
// the balancer will let the broker tell the consumer instance to stop processing the partition.
// The existing consumer instance will flush the internal state, and then stop processing.
// Then the balancer will tell the brokers to start sending new messages in the new/moved partition to the consumer instances.
//
// Failover to standby consumer instances:
//
//	A consumer group can have min and max number of consumer instances.
//	For consumer instances joined after the max number, they will be in standby mode.
//
//	When a consumer instance is down, the broker will notice this and inform the balancer.
//	The balancer will then tell the broker to send the partition to another standby consumer instance.
type Balancer struct {
	Brokers cmap.ConcurrentMap[string, *BrokerStats] // key: broker address
	// Collected from all brokers when they connect to the broker leader
	TopicToBrokers    cmap.ConcurrentMap[string, *PartitionSlotToBrokerList] // key: topic name
	OnPartitionChange func(topic *mq_pb.Topic, assignments []*mq_pb.BrokerPartitionAssignment)
	OnAddBroker       func(broker string, brokerStats *BrokerStats)
	OnRemoveBroker    func(broker string, brokerStats *BrokerStats)
}

func NewBalancer() *Balancer {
	return &Balancer{
		Brokers:        cmap.New[*BrokerStats](),
		TopicToBrokers: cmap.New[*PartitionSlotToBrokerList](),
	}
}

func (balancer *Balancer) AddBroker(broker string) (brokerStats *BrokerStats) {
	var found bool
	brokerStats, found = balancer.Brokers.Get(broker)
	if !found {
		brokerStats = NewBrokerStats()
		if !balancer.Brokers.SetIfAbsent(broker, brokerStats) {
			brokerStats, _ = balancer.Brokers.Get(broker)
		}
	}
	balancer.onPubAddBroker(broker, brokerStats)
	balancer.OnAddBroker(broker, brokerStats)
	return brokerStats
}

func (balancer *Balancer) RemoveBroker(broker string, stats *BrokerStats) {
	balancer.Brokers.Remove(broker)

	// update TopicToBrokers
	for _, topic := range stats.Topics {
		partitionSlotToBrokerList, found := balancer.TopicToBrokers.Get(topic.String())
		if !found {
			continue
		}
		pickedBroker := pickBrokers(balancer.Brokers, 1)
		if len(pickedBroker) == 0 {
			partitionSlotToBrokerList.RemoveBroker(broker)
		} else {
			partitionSlotToBrokerList.ReplaceBroker(broker, pickedBroker[0])
		}
	}
	balancer.onPubRemoveBroker(broker, stats)
	balancer.OnRemoveBroker(broker, stats)
}

func (balancer *Balancer) OnBrokerStatsUpdated(broker string, brokerStats *BrokerStats, receivedStats *mq_pb.BrokerStats) {
	brokerStats.UpdateStats(receivedStats)

	// update TopicToBrokers
	for _, topicPartitionStats := range receivedStats.Stats {
		topicKey := topic.FromPbTopic(topicPartitionStats.Topic).String()
		partition := topicPartitionStats.Partition
		partitionSlotToBrokerList, found := balancer.TopicToBrokers.Get(topicKey)
		if !found {
			partitionSlotToBrokerList = NewPartitionSlotToBrokerList(MaxPartitionCount)
			if !balancer.TopicToBrokers.SetIfAbsent(topicKey, partitionSlotToBrokerList) {
				partitionSlotToBrokerList, _ = balancer.TopicToBrokers.Get(topicKey)
			}
		}
		partitionSlotToBrokerList.AddBroker(partition, broker)
	}
}

// OnPubAddBroker is called when a broker is added for a publisher coordinator
func (balancer *Balancer) onPubAddBroker(broker string, brokerStats *BrokerStats) {
}

// OnPubRemoveBroker is called when a broker is removed for a publisher coordinator
func (balancer *Balancer) onPubRemoveBroker(broker string, brokerStats *BrokerStats) {
}
