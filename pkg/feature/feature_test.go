package feature_test

import (
	"errors"
	"testing"

	"github.com/pinkgorilla/go-sample/pkg/feature"
)

func Test_WhenEqualWithInvalidKey_ShouldError(t *testing.T) {
	m := feature.GetManager()
	k := "test1"
	fn := func() error {
		return nil
	}
	err := m.WhenEqual(k, 1, fn)
	if err == nil {
		t.Fatal("expecting error when key not set")
	}
}

func Test_GetStateWithInvalidKey_ShouldSuccess(t *testing.T) {
	m := feature.GetManager()
	k := "test1"
	i := m.GetState(k)
	if i != nil {
		t.Fatal("expecting nil when key not set")
	}
}

func Test_GetStateWithValidKey_ShouldSuccess(t *testing.T) {
	m := feature.GetManager()
	k := "test1"
	m.SetState(k, 1)
	i := m.GetState(k)
	if i != 1 {
		t.Fatal("expecting 1 when key not set")
	}
}

func Test_WhenEqualWithStateAsInteger_ShouldSuccess(t *testing.T) {
	m := feature.GetManager()
	k := "test1"
	m.SetState(k, 1)
	fn := func() error {
		return nil
	}

	err := m.WhenEqual(k, 1, fn)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_WhenEqualWithStateAsString_ShouldSuccess(t *testing.T) {
	m := feature.GetManager()
	k := "test1"
	m.SetState(k, "hello")

	fn := func() error {
		return nil
	}

	err := m.WhenEqual(k, "hello", fn)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_WhenEqualWithStateAsFn_ShouldSuccess(t *testing.T) {
	m := feature.GetManager()
	k := "test1"
	m.SetState(k, func() interface{} { return "hello" })

	fn := func() error {
		return nil
	}

	err := m.WhenEqual(k, "hello", fn)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_WhenEqualWhenStateNotEqual_ShouldSuccess(t *testing.T) {
	m := feature.GetManager()
	k := "test1"
	m.SetState(k, "hello")

	fn := func() error {
		return errors.New("error because this fn run when feature has invalid state")
	}

	err := m.WhenEqual(k, "world", fn)
	if err != nil {
		t.Fatal(err)
	}
}
