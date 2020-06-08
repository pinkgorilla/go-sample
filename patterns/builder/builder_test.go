package builder_test

import "testing"

// Builder Design Pattern - reusing an algorithm to create many implementations of an interface
//
// ## Description
// Instance creation can be as simple as providing the opening and closing braces {} and leaving the instance with zero values,
// or as complex as an object that needs to make some API calls, check states, and create objects for its fields.
// You could also have an object that is composed of many objects, something that's really idiomatic in Go, as it doesn't support inheritance.
// At the same time, you could be using the same technique to create many types of objects.
// For example, you'll use almost the same technique to build a car as you would build a bus,
// except that they'll be of different sizes and number of seats, so why don't we reuse the construction process?
// This is where the Builder pattern comes to the rescue.

// ## Objectives
// A Builder design pattern tries to:
// - Abstract complex creations so that object creation is separated from the object user
// - Create an object step by step by filling its fields and creating the embedded objects”
// - Reuse the object creation algorithm between many objects
//
// ## Example - vehicle manufacturing

// The Builder design pattern has been commonly described as the relationship between a director, a few Builders, and the product they build.
// Continuing with our example of the car, we'll create a vehicle Builder.
// The process (widely described as the algorithm) of creating a vehicle (the product) is more or less the same for every kind of vehicle
// --choose vehicle type, assemble the structure, place the wheels, and place the seats.
// If you think about it, you could build a car and a motorbike (two Builders) with this description,
// so we are reusing the description to create cars in manufacturing.
// The director is represented by the ManufacturingDirector type in our example.

type BuildProcess interface {
	SetWheels() BuildProcess
	SetSeats() BuildProcess
	SetStructure() BuildProcess
	GetVehicle() VehicleProduct
}

// This preceding interface defines the steps that are necessary to build a vehicle.
// Every builder must implement this interface if they are to be used by the manufacturing.
// On every Set step, we return the same build process, so we can chain various steps together in the same statement, as we'll see later.
// Finally, we'll need a GetVehicle method to retrieve the Vehicle instance from the builder:

type ManufacturingDirector struct {
	builder BuildProcess
}

func (f *ManufacturingDirector) SetBuilder(b BuildProcess) {
	f.builder = b
}

func (f *ManufacturingDirector) Construct() {
	f.builder.SetSeats().SetStructure().SetWheels()
}

// The ManufacturingDirector director variable is the one in charge of accepting the builders.
// It has a Construct method that will use the builder that is stored in Manufacturing, and will reproduce the required steps.
// The SetBuilder method will allow us to change the builder that is being used in the Manufacturing director:

type VehicleProduct struct {
	Wheels    int
	Seats     int
	Structure string
}

//  The product is the final object that we want to retrieve while using the manufacturing.
// In this case, a vehicle is composed of wheels, seats, and a structure:
type CarBuilder struct {
	v VehicleProduct
}

func (c *CarBuilder) SetWheels() BuildProcess {
	c.v.Wheels = 4
	return c
}

func (c *CarBuilder) SetSeats() BuildProcess {
	c.v.Seats = 5
	return c
}

func (c *CarBuilder) SetStructure() BuildProcess {
	c.v.Structure = "Car"
	return c
}

func (c *CarBuilder) GetVehicle() VehicleProduct {
	return c.v
}

// The first Builder is the Car builder.
// It must implement every method defined in the BuildProcess interface.
// This is where we'll set the information for this particular builder:
type BikeBuilder struct {
	v VehicleProduct
}

func (b *BikeBuilder) SetWheels() BuildProcess {
	b.v.Wheels = 2
	return b
}

func (b *BikeBuilder) SetSeats() BuildProcess {
	b.v.Seats = 2
	return b
}

func (b *BikeBuilder) SetStructure() BuildProcess {
	b.v.Structure = "Motorbike"
	return b
}

func (b *BikeBuilder) GetVehicle() VehicleProduct {
	return b.v
}

// The Motorbike structure must be the same as the Car structure, as they are all Builder implementations,
// but keep in mind that the process of building each can be very different. With this declaration of objects,
// we can create the following tests:

func TestBuilderPattern(t *testing.T) {
	manufacturingComplex := ManufacturingDirector{}

	carBuilder := &CarBuilder{}
	manufacturingComplex.SetBuilder(carBuilder)
	manufacturingComplex.Construct()

	car := carBuilder.GetVehicle()

	// We will start with the Manufacturing director and the Car Builder to fulfill the first two acceptance criteria.
	// In the preceding code, we are creating our Manufacturing director that will be in charge of the creation of every vehicle during the test.
	// After creating the Manufacturing director, we created a CarBuilder that we then passed to manufacturing by using the SetBuilder method.
	// Once the Manufacturing director knows what it has to construct now, we can call the Construct method to create the VehicleProduct using CarBuilder.
	// Finally, once we have all the pieces for our car, we call the GetVehicle method on CarBuilder to retrieve a Car instance:

	if car.Wheels != 4 {
		t.Errorf("Wheels on a car must be 4 and they were %d\n", car.Wheels)
	}

	if car.Structure != "Car" {
		t.Errorf("Structure on a car must be 'Car' and was %s\n", car.Structure)
	}

	if car.Seats != 5 {
		t.Errorf("Seats on a car must be 5 and they were %d\n", car.Seats)
	}
	// Now we will create tests for a Motorbike builder that covers the third and fourth acceptance criteria:
	bikeBuilder := &BikeBuilder{}

	manufacturingComplex.SetBuilder(bikeBuilder)
	manufacturingComplex.Construct()

	motorbike := bikeBuilder.GetVehicle()
	motorbike.Seats = 1

	if motorbike.Wheels != 2 {
		t.Errorf("Wheels on a motorbike must be 2 and they were %d\n", motorbike.Wheels)
	}

	if motorbike.Structure != "Motorbike" {
		t.Errorf("Structure on a motorbike must be 'Motorbike' and was %s\n", motorbike.Structure)
	}
	// The preceding code is a continuation of the car tests.
	// As you can see, we reuse the previously created manufacturing to create the bike now by passing the Motorbike builder to it.
	// Then we hit the construct button again to create the necessary parts,
	// and call the builder GetVehicle method to retrieve the motorbike instance.
}

// ## Wrapping up the Builder design pattern
// The Builder design pattern helps us maintain an unpredictable number of products by using a common construction algorithm that is used by the director.
// The construction process is always abstracted from the user of the product.
// At the same time, having a defined construction pattern helps when a newcomer to our source code needs to add a new product to the pipeline.
// The BuildProcess interface specifies what he must comply to be part of the possible builders.
// However, try to avoid the Builder pattern when you are not completely sure that the algorithm is going to be more or less stable
// because any small change in this interface will affect all your builders and
// it could be awkward if you add a new method that some of your builders need and others Builders do not.
