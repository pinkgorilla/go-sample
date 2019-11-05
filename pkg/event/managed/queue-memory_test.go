package managed_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/pinkgorilla/go-sample/pkg/event/managed"
)

func Test_InMemoryQueue(t *testing.T) {
	push := func(i interface{}) (interface{}, error) {
		return fmt.Sprint(i), nil
	}
	pop := func(i interface{}) (interface{}, error) {
		s, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("Pop failed:%s", "cannot assert interface{} to string")
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	key := func(i interface{}) interface{} {
		return i
	}
	im := managed.NewInMemoryQueueWithFn(push, pop, key)
	store := im
	// push value to store
	store.Push(1)
	// store should not be empty
	if store.IsEmpty() {
		t.Fatal("isEmpty")
	}
	// pop data from store
	v, err := store.Pull()
	if err != nil {
		t.Fatal(err)
	}
	// if value is not what is pushed should error
	if v != 1 {
		t.Fatal("v")
	}
	// if poped but it is not empty, should error
	if !store.IsEmpty() {
		t.Fatal("isEmpty")
	}
	store.Push(1979)
	store.Push(2088)
	store.Push(2020)
	bs, err := ioutil.ReadAll(store)
	if err != nil {
		t.Fatal(err)
	}

	buff := bytes.NewBuffer(bs)
	x := managed.NewInMemoryQueueWithFn(push, pop, key)
	ns := x
	err = managed.LoadInMemoryQueueWithReader(ns, buff)
	if err != nil {
		t.Fatal(err)
	}

	size, err := ns.Size()
	if err != nil {
		t.Fatal(err)
	}
	if size != 3 {
		t.Fatal("size")
	}
}
