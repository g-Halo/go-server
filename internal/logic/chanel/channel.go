package chanel

import (
	"sync"

	h "github.com/g-Halo/go-server/pkg/util/hash"
)

type ChannelBucket struct {
	Data  map[string]*UCBuff
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
			Data:  map[string]*UCBuff{},
			Index: i,
			mutex: &sync.Mutex{},
		}
		l.Channels = append(l.Channels, item)
	}
	return l
}

func (l *ChannelList) Get(key string) (*UCBuff, *ChannelBucket) {
	b := l.HashInt(key)
	b.mutex.Lock()
	if c, ok := b.Data[key]; ok {
		b.mutex.Unlock()
		return c, b
	} else {
		c = NewUserChanBuff(key)
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
