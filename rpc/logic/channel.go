package logic

import (
	"sync"

	"github.com/g-Halo/go-server/conf"
	h "github.com/g-Halo/go-server/util/hash"
)

var Channels *ChannelList

type ChannelBucket struct {
	// Data  map[string]Channel
	mutex *sync.Mutex
}

type ChannelList struct {
	Channels []*ChannelBucket
}

func NewChannelList() *ChannelList {
	l := &ChannelList{Channels: []*ChannelBucket{}}
	for i := 0; i < conf.Conf.ChannelBucketCount; i++ {
		item := &ChannelBucket{
			mutex: &sync.Mutex{},
		}
		l.Channels = append(l.Channels, item)
	}
	return l
}

func (l *ChannelList) New(key string) {

}

// 通过用户的 username 哈希到某个 bucket
func (l *ChannelList) Get(key string) {

}

func (l *ChannelList) Hash(key string) *ChannelBucket {
	h := h.NewMurmur3C()
	h.Write([]byte(key))
	idx := uint(h.Sum32()) & uint(conf.Conf.ChannelBucketCount-1)
	return l.Channels[idx]
}
