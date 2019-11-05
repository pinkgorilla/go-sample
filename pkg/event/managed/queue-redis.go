package managed

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// DataToJSON Default push function, parse i to json string.
var DataToJSON = func(i interface{}) (interface{}, error) {
	bs, err := json.Marshal(i)
	return string(bs), err
}

// JSONToData is default pop function, parse json string to map[string]interface{}
var JSONToData = func(i interface{}) (interface{}, error) {
	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("Pop failed:%s", "cannot assert interface{} to string")
	}
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// RedisQueue is queue implementation on redis
type RedisQueue struct {
	key  string
	r    *redis.Client
	push func(interface{}) (interface{}, error)
	pop  func(interface{}) (interface{}, error)
}

// NewRedisQueue returns new RedisQueue instance
func NewRedisQueue(r *redis.Client) *RedisQueue {
	return NewRedisQueueWithFunc(r, DataToJSON, JSONToData)
}

// NewRedisQueueWithFunc returns new RedisStore with specified push and pop func
func NewRedisQueueWithFunc(
	r *redis.Client,
	push func(interface{}) (interface{}, error),
	pop func(interface{}) (interface{}, error)) *RedisQueue {
	return &RedisQueue{
		key:  fmt.Sprintf("queue:%v", time.Now().Nanosecond()),
		r:    r,
		push: push,
		pop:  pop,
	}
}

// Push pushes data to redis stream
func (s *RedisQueue) Push(data interface{}) error {
	str, err := s.push(data)
	if err != nil {
		return err
	}
	_, err = s.r.LPush(s.Key(), str).Result()
	return err
}

// Pull pulls data from redis stream
func (s *RedisQueue) Pull() (interface{}, error) {
	r, err := s.r.RPop(s.Key()).Result()
	if err != nil {
		return nil, err
	}
	return s.pop(r)
}

// Dispose clean resources
func (s *RedisQueue) Dispose() {
	log.Println("dispose")
	s.Clear()
}

// Key returns stream unique key
func (s *RedisQueue) Key() string {
	return s.key
}

// Clear clears the stream
func (s *RedisQueue) Clear() error {
	return s.r.Del(s.Key()).Err()
}

//Size returns stream size
func (s *RedisQueue) Size() (int, error) {
	size, err := s.r.LLen(s.Key()).Result()
	return int(size), err
}
