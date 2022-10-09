package structs

import "github.com/VincentFF/thinredis/config"

const (
	REDIS_STRING = iota
	REDIS_LIST
	REDIS_SET
	REDIS_ZSET
	REDIS_HASH
)

var ShardNum = config.Configures.ShardNum
var RedisTypes = [5]string{"string", "list", "set", "zset", "hash"}

type RedisObject struct {
	Type int
	Lru  int64
	Ptr  any
}

type RedisDb struct {
	Id           int
	Table        *Dict
	Expires      *Dict
	WatchedKeys  *Dict
	BlockingKeys *Dict
	ReadyKeys    *Dict
}
