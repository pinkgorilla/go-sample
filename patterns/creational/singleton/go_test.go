package singleton_test

import "testing"

// Singleton design pattern - having a unique instance of a type in the entire program
//
// ### Description
// The Singleton pattern is easy to remember.
// As the name implies, it will provide you with a single instance of an object,
// and guarantee that there are no duplicates
// You'll use the Singleton pattern in many different situations. For example:
// - When you want to use the same connection to a database to make every query
// - When you open a Secure Shell (SSH) connection to a server to do a few tasks, and don't want to reopen the connection for each task
// - If you need to limit the access to some variable or space, you use a Singleton as the door to this variable (we'll see in the following chapters that this is more achievable in Go using channels anyway)
// - If you need to limit the number of calls to some places, you create a Singleton instance to make the calls in the accepted window”
//
// ### Objectives
// As a general guide, we consider using the Singleton pattern when the following rule applies:
// - We need a single, shared value, of some particular type.
// - We need to restrict object creation of some type to a single unit along the entire program.”

// ## Example - a unique counter
// As an example of an object of which we must ensure that there is only one instance, we will write a counter that holds the number of times it has been called during program execution. It shouldn't matter how many instances we have of the counter, all of them must count the same value and it must be consistent between the instances.
// ### Requirements and acceptance criteria
// There are some requirements and acceptance criteria to write the described single counter. They are as follows:
// When no counter has been created before, a new one is created with the value 0
// - If a counter has already been created, return this instance that holds the actual count
// - If we call the method AddOne, the count must be incremented by 1

type Singleton interface {
	AddOne() int
}
type singleton struct {
	count int
}

var instance *singleton

func GetInstance() Singleton {
	if instance == nil {
		instance = new(singleton)
	}
	return instance
}
func (s *singleton) AddOne() int {
	s.count++
	return s.count

}

func Test_UniqueCounter(t *testing.T) {
	counter1 := GetInstance()
	if counter1 == nil {
		//Test of acceptance criteria 1 failed
		t.Error("expected pointer to Singleton after calling GetInstance(), not nil")
	}
	expectedCounter := counter1

	currentCount := counter1.AddOne()
	if currentCount != 1 {
		t.Errorf("After calling for the first time to count, the count must be 1 but it is %d\n", currentCount)
	}
	counter2 := GetInstance()
	if counter2 != expectedCounter {
		//Test 2 failed
		t.Error("Expected same instance in counter2 but it got a different instance")
	}
	currentCount = counter2.AddOne()
	if currentCount != 2 {
		t.Errorf("After calling 'AddOne' using the second counter, the current count must be 2 but was %d\n", currentCount)
	}
}

// ## A few words about the Singleton design pattern
// We have seen a very simple example of the Singleton pattern, partially applied to some situation, that is, a simple counter.
// Just keep in mind that the Singleton pattern will give you the power to have a unique instance of some struct in
// your application and that no package can create any clone of this struct.
//
// With Singleton, you are also hiding the complexity of creating the object, in case it requires some computation,
// and the pitfall of creating it every time you need an instance of it if all of them are similar.
// All this code writing, checking if the variable already exists, and storage, are encapsulated in the singleton and you won't need to repeat it everywhere if you use a global variable.
//
// Here we are learning the classic singleton implementation for single threaded context. We will see a concurrent singleton implementation when we reach the chapters about concurrency because this implementation is not thread safe!
// Excerpt From: “Go: Design Patterns for Real-World Projects.”
