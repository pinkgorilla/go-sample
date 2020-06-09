package factory_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// ## Description
// When using the Factory method design pattern, we gain an extra layer of encapsulation so that our program can grow in a controlled environment.
// With the Factory method, we delegate the creation of families of objects to a different package or object to abstract us from the knowledge of the pool of possible objects we could use.
// Imagine that you want to organize your holidays using a trip agency.
// You don't deal with hotels and traveling and you just tell the agency the destination you are interested in so that they provide you with everything you need.
// The trip agency represents a Factory of trips.
//
// ## Objectives
// After the previous description, the following objectives of the Factory Method design pattern must be clear to you:
// - Delegating the creation of new instances of structures to a different part of the program
// - Working at the interface level instead of with concrete implementations
// - Grouping families of objects to obtain a family object creator
//
// ## Example - a factory of payment methods for a shop
// For our example, we are going to implement a payments method Factory,
// which is going to provide us with different ways of paying at a shop.
// In the beginning, we will have two methods of paying--cash and credit card.
// We'll also have an interface with the method, Pay, which every struct that wants to be used as a payment method must implement.
//
// ### Acceptance criteria
// Using the previous description, the requirements for the acceptance criteria are the following:
// - To have a common method for every payment method called Pay
// - To be able to delegate the creation of payments methods to the Factory
// - To be able to add more payment methods to the library by just adding it to the factory method

// A Factory method has a very simple structure; we just need to identify how many implementations of our interface we are storing,
// and then provide a method, GetPaymentMethod, where you can pass a type of payment as an argument:

// The following lines define the interface of the payment method.
// They define a way of making a payment at the shop.
// The Factory method will return instances of types that implement this interface:
type PaymentMethod interface {
	Pay(amount float32) string
}

// We have to define the identified payment methods of the Factory
// as constants so that we can call and check the possible payment methods from outside of the package.
const (
	Cash      = 1
	DebitCard = 2
)

// The following code is the function that will create the objects for us (factory method).
// It returns a pointer, which must have an object that implements the PaymentMethod interface,
// and an error if asked for a method which is not registered.
func GetPaymentMethod(m int) (PaymentMethod, error) {
	switch m {
	case Cash:
		return new(CashPM), nil
	case DebitCard:
		return new(DebitCardPM), nil
	default:
		return nil, errors.New(fmt.Sprintf("Payment method %d  “not recognized\n", m))
	}
}

// To finish the declaration of the Factory, we create the two payment methods.
// As you can see, the CashPM and DebitCardPM structs implement the PaymentMethod interface by declaring a method, Pay(amount float32) string.
// The returned string will contain information about the payment.

type CashPM struct{}
type DebitCardPM struct{}

func (c *CashPM) Pay(amount float32) string {
	return fmt.Sprintf("%0.2f paid using cash\n", amount)
}

func (c *DebitCardPM) Pay(amount float32) string {
	return fmt.Sprintf("%#0.2f paid using debit card\n", amount)
}

// we will start by writing the tests for the first acceptance criteria:
// to have a common method to retrieve objects that implement the PaymentMethod interface:
func TestCreatePaymentMethodCash(t *testing.T) {
	// GetPaymentMethod is a common method to retrieve methods of payment. We use the constant Cash
	payment, err := GetPaymentMethod(Cash)
	if err != nil {
		t.Fatal("A payment method of type 'Cash' must exist")
	}

	msg := payment.Pay(10.30)
	if !strings.Contains(msg, "paid using cash") {
		t.Error("The cash payment method message wasn't correct")
	}
	t.Log("LOG:", msg)

	// We repeat the same operation with the debit card method.
	// We ask for the payment method defined with the constant DebitCard, and the returned message,
	// when paying with debit card, must contain the paid using debit card string.
	payment, err = GetPaymentMethod(DebitCard)
	if err != nil {
		t.Error("A payment method of type 'DebitCard' must exist")
	}

	msg = payment.Pay(22.30)
	if !strings.Contains(msg, "paid using debit card") {
		t.Error("The debit card payment method message wasn't correct")
	}
	t.Log("LOG:", msg)

	// Finally, we are going to test the situation when we request a payment method that doesn´t exist
	// (represented by the number 20, which doesn't match any recognized constant in the Factory)
	payment, err = GetPaymentMethod(20)
	if err == nil {
		t.Error("A payment method with ID 20 must return an error")
	}
	t.Log("LOG:", err)
}

// Excerpt From: “Go: Design Patterns for Real-World Projects.” Apple Books.
