package managed_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/go-redis/redis"
	"github.com/pinkgorilla/go-sample/pkg/event/managed"
)

func Test_RedisStore_Struct_Push_Pop(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	key := "store:test"
	r.Del(key)
	type T struct {
		I int    `json:"i"`
		S string `json:"s"`
	}
	data := T{99, "Baloons"}
	s := managed.NewRedisStoreWithFunc(key, r,
		func(i interface{}) (interface{}, error) {
			bs, err := json.Marshal(i)
			return string(bs), err
		},
		func(i interface{}) (interface{}, error) {
			s, ok := i.(string)
			if !ok {
				return nil, fmt.Errorf("Pop failed:%s", "cannot assert interface{} to string")
			}
			var t T
			err := json.Unmarshal([]byte(s), &t)
			if err != nil {
				return nil, err
			}
			return t, nil
		})
	err := s.Push(data)
	if err != nil {
		t.Fatal(err)
	}

	size, err := s.Size()
	if err != nil {
		t.Fatal(err)
	}
	if size != 1 {
		t.Fatal("invalid size")
	}

	p, err := s.Pop()
	if err != nil {
		t.Fatal(err)
	}
	result, ok := p.(T)
	if !ok {
		t.Fatal("assert failed")
	}
	if result.I != data.I {
		t.Fatal("mismatch!")
	}
	size, err = s.Size()
	if err != nil {
		t.Fatal(err)
	}
	if size != 0 {
		t.Fatal("invalid size")
	}
}

func Test_RedisStore_String_Push_Pop(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	key := "store:test"
	r.Del(key)
	data := "hello-storage"
	s := managed.NewRedisStoreWithFunc(key, r,
		func(i interface{}) (interface{}, error) {
			return fmt.Sprint(i), nil
		},
		func(i interface{}) (interface{}, error) {
			s, ok := i.(string)
			if !ok {
				return nil, fmt.Errorf("Pop failed:%s", "cannot assert interface{} to string")
			}
			return s, nil
		})
	err := s.Push(data)
	if err != nil {
		t.Fatal(err)
	}

	p, err := s.Pop()
	if err != nil {
		t.Fatal(err)
	}
	result, ok := p.(string)
	if !ok {
		t.Fatal("assert failed")
	}
	if result != data {
		t.Fatal("mismatch!")
	}
}

func Test_RedisStore_Int_Push_Pop(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	key := "store:test"
	r.Del(key)
	data := 1987
	s := managed.NewRedisStoreWithFunc(key, r,
		func(i interface{}) (interface{}, error) {
			return i, nil
		},
		func(i interface{}) (interface{}, error) {
			s, ok := i.(string)
			if !ok {
				return nil, fmt.Errorf("Pop failed:%s", "cannot assert interface{} to string")
			}
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
			return int(v), nil
		})
	err := s.Push(data)
	if err != nil {
		t.Fatal(err)
	}

	p, err := s.Pop()
	if err != nil {
		t.Fatal(err)
	}
	result, ok := p.(int)
	if !ok {
		t.Fatal("assert failed")
	}
	if result != data {
		t.Fatal("mismatch!")
	}
}
