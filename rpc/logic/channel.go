package logic

import (
	"sync"

	"github.com/g-Halo/go-server/model"
	h "github.com/g-Halo/go-server/util/hash"
)

var UserChannels *ChannelList
var RoomChannels *ChannelList

type Chan interface {
	PushMsg(key string, m *model.Message)
	GetMsg(key string) *model.Message
}

type ChannelBucket struct {
	Data  map[string]Chan
	Index int
	mutex *sync.Mutex
}

type ChannelList struct {
	Channels    []*ChannelBucket
	BucketCount int
}

func NewChannelList(bucketCount int) *ChannelList {
	l := &ChannelList{Channels: []*ChannelBucket{}, BucketCount: bucketCount}
	for i := 0; i < bucketCount; i++ {
		item := &ChannelBucket{
			Data:  map[string]Chan{},
			Index: i,
			mutex: &sync.Mutex{},
		}
		l.Channels = append(l.Channels, item)
	}
	return l
}

// 通过用户的 username 哈希到某个 bucket
func (l *ChannelList) Get(key string) (Chan, *ChannelBucket) {
	b := l.HashInt(key)
	b.mutex.Lock()
	if c, ok := b.Data[key]; ok {
		b.mutex.Unlock()
		return c, b
	} else {
		c = NewRoomChan(key)
		b.Data[key] = c
		b.mutex.Unlock()
		return c, b
	}
}

func (l *ChannelList) HashInt(key string) *ChannelBucket {
	h := h.NewMurmur3C()
	h.Write([]byte(key))
	idx := uint(h.Sum32()) & uint(l.BucketCount-1)
	return l.Channels[idx]
}
