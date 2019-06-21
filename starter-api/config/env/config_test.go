package env

import (
	"os"
	"testing"
)

func Test_GetBoolean_1(t *testing.T) {
	val := "true"
	key := "CO_TEST_BOOLEAN_1"
	os.Setenv(key, val)
	b := getBoolean(key)
	defer os.Unsetenv(key)
	if b != true {
		t.Fatalf("expected %v got %v", true, b)
	}
}

func Test_GetBoolean_2(t *testing.T) {
	val := "false"
	key := "CO_TEST_BOOLEAN_2"
	os.Setenv(key, val)
	b := getBoolean(key)
	defer os.Unsetenv(key)
	if b != false {
		t.Fatalf("expected %v got %v", false, b)
	}
}

func Test_GetBoolean_3(t *testing.T) {
	val := ""
	key := "CO_TEST_BOOLEAN_3"
	os.Setenv(key, val)
	b := getBoolean(key)
	defer os.Unsetenv(key)
	if b != false {
		t.Fatalf("expected %v got %v", false, b)
	}
}
func Test_GetBoolean_4(t *testing.T) {
	val := "abc"
	key := "CO_TEST_BOOLEAN_4"
	os.Setenv(key, val)
	b := getBoolean(key)
	defer os.Unsetenv(key)
	if b != false {
		t.Fatalf("expected %v got %v", false, b)
	}
}
func Test_GetBoolean_5(t *testing.T) {
	key := "CO_TEST_BOOLEAN_5"
	err := os.Unsetenv(key)
	if err != nil {
		t.Fatal(err)
	}
	b := getBoolean(key)
	if b != false {
		t.Fatalf("expected %v got %v", false, b)
	}
}

func Test_GetInt_1(t *testing.T) {
	val := "100"
	key := "CO_TEST_INT_1"
	os.Setenv(key, val)
	i := getInt(key)
	defer os.Unsetenv(key)
	if i != 100 {
		t.Fatalf("expected %v got %v", 100, i)
	}
}

func Test_GetInt_2(t *testing.T) {
	val := "0"
	key := "CO_TEST_INT_2"
	os.Setenv(key, val)
	i := getInt(key)
	defer os.Unsetenv(key)
	if i != 0 {
		t.Fatalf("expected %v got %v", false, i)
	}
}

func Test_GetInt_3(t *testing.T) {
	val := ""
	key := "CO_TEST_INT_3"
	os.Setenv(key, val)
	i := getInt(key)
	defer os.Unsetenv(key)
	if i != 0 {
		t.Fatalf("expected %v got %v", 0, i)
	}
}

func Test_GetInt_4(t *testing.T) {
	val := "abc"
	key := "CO_TEST_INT_4"
	os.Setenv(key, val)
	i := getInt(key)
	defer os.Unsetenv(key)
	if i != 0 {
		t.Fatalf("expected %v got %v", 0, i)
	}
}
func Test_GetInt_5(t *testing.T) {
	key := "CO_TEST_INT_5"
	err := os.Unsetenv(key)
	if err != nil {
		t.Fatal(err)
	}
	i := getInt(key)
	defer os.Unsetenv(key)
	if i != 0 {
		t.Fatalf("expected %v got %v", 0, i)
	}
}
