package abstract_factory_test

import (
	"errors"
	"fmt"
	"testing"
)

// Abstract Factory - a factory of factories
// ## Description
// The Abstract Factory design pattern is a new layer of grouping to achieve a bigger (and more complex)
// composite object, which is used through its interfaces.
// The idea behind grouping objects in families and grouping families is to have big factories that can be interchangeable and can grow more easily.
// In the early stages of development, it is also easier to work with factories and
// abstract factories than to wait until all concrete implementations are done to start your code.
// Also, you won't write an Abstract Factory from the beginning unless you know that your object's inventory
// for a particular field is going to be very large and it could be easily grouped into families.

// ## The objectives
// Grouping related families of objects is very convenient when your object number is growing so much that creating
// a unique point to get them all seems the only way to gain the flexibility of the runtime object creation.
// The following objectives of the Abstract Factory method must be clear to you:
// - Provide a new layer of encapsulation for Factory methods that return a common interface for all factories
// - Group common factories into a super Factory (also called a factory of factories)

// ## The vehicle factory example, again?
// For our example, we are going to reuse the Factory we created in the Builder design pattern.
// We want to show the similarities to solve the same problem using a different approach so that you can see the strengths and weaknesses of each “approach.
// This is going to show you the power of implicit interfaces in Go, as we won't have to touch almost anything.
// Finally, we are going to create a new Factory to create shipment orders.

// Vehicle : The interface that all objects in our factories must implement:
type Vehicle interface {
	NumWheels() int
	NumSeats() int
}

// Car : An interface for cars of types luxury (with four doors) and family (with five doors).
type Car interface {
	NumDoors() int
}

// Motorbike : An interface for motorbikes of the types sport (one seat) and cruise (two seats).
type Motorbike interface {
	GetMotorbikeType() int
}

// VehicleFactory: An interface (the Abstract Factory) to retrieve factories that implement the VehicleFactory method:
// Motorbike Factory: A factory that implements the VehicleFactory interface to return vehicle that implements the Vehicle and Motorbike interfaces.
type VehicleFactory interface {
	Build(v int) (Vehicle, error)
}

const (
	CarFactoryType       = 1
	MotorbikeFactoryType = 2
)

func BuildFactory(f int) (VehicleFactory, error) {
	switch f {
	case CarFactoryType:
		return new(CarFactory), nil
	case MotorbikeFactoryType:
		return new(MotorbikeFactory), nil
	default:
		return nil, errors.New(fmt.Sprintf("Factory with id %d not recognized\n", f))
	}
}

const (
	LuxuryCarType = 1
	FamilyCarType = 2
)

// Car Factory: A factory that implements the VehicleFactory interface to return vehicles that implement the Vehicle and Car interfaces.
type CarFactory struct{}

func (c *CarFactory) Build(v int) (Vehicle, error) {
	switch v {
	case LuxuryCarType:
		return new(LuxuryCar), nil
	case FamilyCarType:
		return new(FamilyCar), nil
	default:
		return nil, errors.New(fmt.Sprintf("Vehicle of type %d not recognized\n", v))
	}
}

type LuxuryCar struct{}

func (*LuxuryCar) NumDoors() int {
	return 4
}
func (*LuxuryCar) NumWheels() int {
	return 4
}
func (*LuxuryCar) NumSeats() int {
	return 5
}

type FamilyCar struct{}

func (*FamilyCar) NumDoors() int {
	return 5
}
func (*FamilyCar) NumWheels() int {
	return 4
}
func (*FamilyCar) NumSeats() int {
	return 5
}

const (
	SportMotorbikeType  = 1
	CruiseMotorbikeType = 2
)

// Now we need the motorbike factory, which, like the car factory, must implement the VehicleFactory interface:
// For the motorbike Factory, we have also defined two types of motorbikes using the const keywords: SportMotorbikeType and CruiseMotorbikeType.
// We will switch over the v argument in the Build method to know which type shall be returned.
type MotorbikeFactory struct{}

func (m *MotorbikeFactory) Build(v int) (Vehicle, error) {
	switch v {
	case SportMotorbikeType:
		return new(SportMotorbike), nil
	case CruiseMotorbikeType:
		return new(CruiseMotorbike), nil
	default:
		return nil, errors.New(fmt.Sprintf("Vehicle of type %d not recognized\n", v))
	}
}

type SportMotorbike struct{}

func (s *SportMotorbike) NumWheels() int {
	return 2
}
func (s *SportMotorbike) NumSeats() int {
	return 1
}
func (s *SportMotorbike) GetMotorbikeType() int {
	return SportMotorbikeType
}

type CruiseMotorbike struct{}

func (c *CruiseMotorbike) NumWheels() int {
	return 2
}
func (c *CruiseMotorbike) NumSeats() int {
	return 2
}
func (c *CruiseMotorbike) GetMotorbikeType() int {
	return CruiseMotorbikeType
}

func TestMotorbikeFactory(t *testing.T) {
	motorbikeF, err := BuildFactory(MotorbikeFactoryType)
	if err != nil {
		t.Fatal(err)
	}

	motorbikeVehicle, err := motorbikeF.Build(SportMotorbikeType)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Motorbike vehicle has %d wheels\n", motorbikeVehicle.NumWheels())

	sportBike, ok := motorbikeVehicle.(Motorbike)
	if !ok {
		t.Fatal("Struct assertion has failed")
	}
	t.Logf("Sport motorbike has type %d\n", sportBike.GetMotorbikeType())
}

// Here's a graph to help you visualize the relationships between entities
// BuildFactory() returns VehicleFactory
// 							|__ CarFactory
// 									|__ Build()
// 											|__ LuxuryCar implements Vehicle, Car
// 											|__ FamilyCar implements Vehicle, Car
// 							|__ MotorBikeFactory
// 									|__ Build()
// 											|__ SportMotorbike implements Vehicle, Motorbike
// 											|__ CruiseMotorbike implements Vehicle, Motorbike
//
// ## A few lines about the Abstract Factory method
// We have learned how to write a factory of factories that provides us with a very generic object of vehicle type.
// This pattern is commonly used in many applications and libraries, such as cross-platform GUI libraries.
// Think of a button, a generic object, and button factory that provides you with a factory for Microsoft Windows buttons while you have another factory for Mac OS X buttons.
// You don't want to deal with the implementation details of each platform, but you just want to implement the actions for some specific behavior raised by a button.
// Also, we have seen the differences when approaching the same problem with two different solutions--the Abstract factory and the Builder pattern.
// As you have seen, with the Builder pattern, we had an unstructured list of objects (cars with motorbikes in the same factory).
// Also, we encouraged reusing the building algorithm in the Builder pattern.
// In the Abstract factory, we have a very structured list of vehicles (the factory for motorbikes and a factory for cars).
// We also didn't mix the creation of cars with motorbikes, providing more flexibility in the creation process.
// The Abstract factory and Builder patterns can both resolve the same problem, but your particular needs will help you find the slight differences that should lead you to take one solution or the other.

// Excerpt From: “Go: Design Patterns for Real-World Projects.” Apple Books.
