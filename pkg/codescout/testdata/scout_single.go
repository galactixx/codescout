package somepackage

import "fmt"

// Structs
type Person struct {
	Name string
	Age  int
}

type Car struct {
	Make  string
	Model string
	Year  int
}

// Above above function
// Above function
// Function
func Greet(p Person) string {
	return fmt.Sprintf("Hello, %s!", p.Name)
}

// Below function

// Method on Person struct
func (p *Person) Birthday() {
	p.Age++
}

// Method on Car struct
func (c *Car) DisplayDetails() string {
	return fmt.Sprintf("%d %s %s", c.Year, c.Make, c.Model)
}

// Variables
var DefaultGreeting = "Welcome to Go!"

var cars = []Car{
	{"Toyota", "Camry", 2020},
	{"Honda", "Accord", 2021},
}
