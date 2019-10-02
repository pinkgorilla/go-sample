package managed

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

// RedisStore ...
type RedisStore struct {
	key   string
	redis *redis.Client
	push  func(interface{}) (interface{}, error)
	pop   func(interface{}) (interface{}, error)
}

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

// NewRedisStore returns new RedisStore with default push and pop func
func NewRedisStore(key string, r *redis.Client) Store {
	return NewRedisStoreWithFunc(key, r,
		DataToJSON,
		JSONToData)
}

// NewRedisStoreWithFunc returns new RedisStore with specified push and pop func
func NewRedisStoreWithFunc(
	key string,
	r *redis.Client,
	push func(interface{}) (interface{}, error),
	pop func(interface{}) (interface{}, error)) Store {
	return &RedisStore{
		key:   key,
		redis: r,
		push:  push,
		pop:   pop,
	}
}

//Push ...
func (s *RedisStore) Push(data interface{}) error {
	str, err := s.push(data)
	if err != nil {
		return err
	}
	_, err = s.redis.LPush(s.key, str).Result()
	return err
}

//Pop ...
func (s *RedisStore) Pop() (interface{}, error) {
	r, err := s.redis.RPop(s.key).Result()
	if err != nil {
		return nil, err
	}
	return s.pop(r)
}

// Size ...
func (s *RedisStore) Size() (int, error) {
	i, err := s.redis.LLen(s.key).Result()
	return int(i), err
}

//Dispose ...
func (s *RedisStore) Dispose() {

}
