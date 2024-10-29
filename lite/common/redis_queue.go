package common

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	"time"
)

type QueueConfig struct {
	Address  string        `json:"address"`
	Password string        `json:"password"`
	DB       int           `json:"db"`
	Timeout  time.Duration `json:"timeout"`
}

type RedisQueue struct {
	Address  string
	Password string
	DB       int
	rdb      *redis.Client
	ctx      context.Context
	timeout  time.Duration
}

func CreateRedisQueue(params *QueueConfig) (*RedisQueue, error) {
	q := &RedisQueue{
		Address:  params.Address,
		Password: params.Password,
		DB:       params.DB,
		ctx:      context.TODO(),
	}
	if params.Timeout == 0 {
		q.timeout = time.Second * 60
	}

	q.rdb = redis.NewClient(&redis.Options{
		Addr:     params.Address,
		Password: params.Password,
		DB:       params.DB,
	})

	_, err := q.rdb.Ping(q.ctx).Result()
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (this *RedisQueue) PutJson(name string, value interface{}) error {
	marshal, _ := json.Marshal(value)
	_, err := this.rdb.LPush(this.ctx, name, string(marshal)).Result()
	return err
}

func (this *RedisQueue) Put(name string, value string) error {
	_, err := this.rdb.LPush(this.ctx, name, value).Result()
	return err
}

func (this *RedisQueue) BLGet(name string) (string, error) {
	result, err := this.rdb.BLPop(this.ctx, this.timeout, name).Result()
	if err != nil {
		return "", err
	}
	return result[1], err
}
