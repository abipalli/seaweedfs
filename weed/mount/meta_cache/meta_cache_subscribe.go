package meta_cache

import (
	"context"
	"github.com/gateway-dao/seaweedfs/weed/filer"
	"github.com/gateway-dao/seaweedfs/weed/glog"
	"github.com/gateway-dao/seaweedfs/weed/pb"
	"github.com/gateway-dao/seaweedfs/weed/pb/filer_pb"
	"github.com/gateway-dao/seaweedfs/weed/util"
)

func SubscribeMetaEvents(mc *MetaCache, selfSignature int32, client filer_pb.FilerClient, dir string, lastTsNs int64) error {

	processEventFn := func(resp *filer_pb.SubscribeMetadataResponse) error {
		message := resp.EventNotification

		for _, sig := range message.Signatures {
			if sig == selfSignature && selfSignature != 0 {
				return nil
			}
		}

		dir := resp.Directory
		var oldPath util.FullPath
		var newEntry *filer.Entry
		if message.OldEntry != nil {
			oldPath = util.NewFullPath(dir, message.OldEntry.Name)
			glog.V(4).Infof("deleting %v", oldPath)
		}

		if message.NewEntry != nil {
			if message.NewParentPath != "" {
				dir = message.NewParentPath
			}
			key := util.NewFullPath(dir, message.NewEntry.Name)
			glog.V(4).Infof("creating %v", key)
			newEntry = filer.FromPbEntry(dir, message.NewEntry)
		}
		err := mc.AtomicUpdateEntryFromFiler(context.Background(), oldPath, newEntry)
		if err == nil {
			if message.OldEntry != nil && message.NewEntry != nil {
				oldKey := util.NewFullPath(resp.Directory, message.OldEntry.Name)
				mc.invalidateFunc(oldKey, message.OldEntry)
				if message.OldEntry.Name != message.NewEntry.Name {
					newKey := util.NewFullPath(dir, message.NewEntry.Name)
					mc.invalidateFunc(newKey, message.NewEntry)
				}
			} else if filer_pb.IsCreate(resp) {
				// no need to invalidate
			} else if filer_pb.IsDelete(resp) {
				oldKey := util.NewFullPath(resp.Directory, message.OldEntry.Name)
				mc.invalidateFunc(oldKey, message.OldEntry)
			}
		}

		return err

	}

	metadataFollowOption := &pb.MetadataFollowOption{
		ClientName:             "mount",
		ClientId:               selfSignature,
		ClientEpoch:            1,
		SelfSignature:          selfSignature,
		PathPrefix:             dir,
		AdditionalPathPrefixes: nil,
		DirectoriesToWatch:     nil,
		StartTsNs:              lastTsNs,
		StopTsNs:               0,
		EventErrorType:         pb.FatalOnError,
	}
	util.RetryUntil("followMetaUpdates", func() error {
		metadataFollowOption.ClientEpoch++
		return pb.WithFilerClientFollowMetadata(client, metadataFollowOption, processEventFn)
	}, func(err error) bool {
		glog.Errorf("follow metadata updates: %v", err)
		return true
	})

	return nil
}
