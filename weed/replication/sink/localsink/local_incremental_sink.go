package localsink

import (
	"github.com/gateway-dao/seaweedfs/weed/replication/sink"
)

type LocalIncSink struct {
	LocalSink
}

func (localincsink *LocalIncSink) GetName() string {
	return "local_incremental"
}

func init() {
	sink.Sinks = append(sink.Sinks, &LocalIncSink{})
}
