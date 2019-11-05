package managed_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/go-redis/redis"
	"github.com/pinkgorilla/go-sample/pkg/event/managed"
)

func Test_RedisQueue(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rs := managed.NewRedisQueueWithFunc(r,
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
	defer rs.Dispose()

	type scenario struct {
		action func(s *managed.RedisQueue)
		test   func(s *managed.RedisQueue) error
	}
	sets := []*scenario{
		&scenario{
			action: func(s *managed.RedisQueue) {
				s.Clear()
				rs.Push("one")
				rs.Push("two")
				rs.Push("three")
			},
			test: func(s *managed.RedisQueue) error {
				if size, _ := s.Size(); size != 3 {
					return fmt.Errorf("test push a - %v", "invalid size")
				}
				values := []string{"one", "two", "three"}
				for _, val := range values {
					v, e := s.Pull()
					if e != nil {
						return fmt.Errorf("test push b - %v", "invalid size")
					}
					if v != val {
						return fmt.Errorf("test push c - %v", "invalid size")
					}
				}
				if size, _ := s.Size(); size != 0 {
					return fmt.Errorf("test push d - %v", "invalid size")
				}
				return nil
			},
		},
	}

	for _, s := range sets {
		s.action(rs)
		if err := s.test(rs); err != nil {
			t.Error(err)
		}
	}
}

func Test_RedisQueue_Struct_Push_Pop(t *testing.T) {
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
	s := managed.NewRedisQueueWithFunc(r,
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

	p, err := s.Pull()
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

func Test_RedisQueue_String_Push_Pop(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	key := "store:test"
	r.Del(key)
	data := "hello-storage"
	s := managed.NewRedisQueueWithFunc(r,
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

	p, err := s.Pull()
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

func Test_RedisQueue_Int_Push_Pop(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	key := "store:test"
	r.Del(key)
	data := 1987
	s := managed.NewRedisQueueWithFunc(r,
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

	p, err := s.Pull()
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
