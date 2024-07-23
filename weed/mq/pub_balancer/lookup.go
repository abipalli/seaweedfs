package pub_balancer

import (
	"errors"
	"github.com/gateway-dao/seaweedfs/weed/pb/mq_pb"
)

var (
	ErrNoBroker = errors.New("no broker")
)

func (balancer *Balancer) LookupTopicPartitions(topic *mq_pb.Topic) (assignments []*mq_pb.BrokerPartitionAssignment) {
	// find existing topic partition assignments
	for brokerStatsItem := range balancer.Brokers.IterBuffered() {
		broker, brokerStats := brokerStatsItem.Key, brokerStatsItem.Val
		for topicPartitionStatsItem := range brokerStats.TopicPartitionStats.IterBuffered() {
			topicPartitionStat := topicPartitionStatsItem.Val
			if topicPartitionStat.TopicPartition.Namespace == topic.Namespace &&
				topicPartitionStat.TopicPartition.Name == topic.Name {
				assignment := &mq_pb.BrokerPartitionAssignment{
					Partition: &mq_pb.Partition{
						RingSize:   MaxPartitionCount,
						RangeStart: topicPartitionStat.RangeStart,
						RangeStop:  topicPartitionStat.RangeStop,
						UnixTimeNs: topicPartitionStat.UnixTimeNs,
					},
				}
				// TODO fix follower setting
				assignment.LeaderBroker = broker
				assignments = append(assignments, assignment)
			}
		}
	}
	return
}
