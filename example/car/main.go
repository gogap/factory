package main

import (
	"fmt"

	"github.com/gogap/factory"
)

type Wheel struct {
	ID string
}

func NewWheel(opts factory.Options) (wheel interface{}, err error) {
	w := &Wheel{}
	opts.Get("id", &w.ID)
	wheel = w
	return
}

func (p *Wheel) Run() {
	fmt.Printf("Wheel Running, ID: %s\n", p.ID)
}

type Car struct {
	owner string

	Wheel1 *Wheel
	Wheel2 *Wheel
	Wheel3 *Wheel
	Wheel4 *Wheel
}

func (p *Car) Run() {
	p.Wheel1.Run()
	p.Wheel2.Run()
	p.Wheel3.Run()
	p.Wheel4.Run()

	fmt.Printf("%s' Car Running\n", p.owner)
}

func NewCar(opts factory.Options) (car interface{}, err error) {
	c := &Car{}
	opts.Get("owner", &c.owner)
	car = c
	return
}

func init() {
	factory.RegisterModel((*Wheel)(nil), "Michelin")
	factory.RegisterModel((*Car)(nil), "Skoda")
}

func main() {
	carFactory := factory.NewClassicFactory(nil)

	carFactory.Define("wheel", factory.Prototype, "Michelin", factory.DefOptOfNewObjectFunc(NewWheel))

	carFactory.Define("mycar", factory.Prototype, "Skoda",
		factory.DefOptOfNewObjectFunc(NewCar),
		factory.DefOptOfObjectRef("Wheel1", "wheel", factory.Options{"id": "1"}),
		factory.DefOptOfObjectRef("Wheel2", "wheel", factory.Options{"id": "2"}),
		factory.DefOptOfObjectRef("Wheel3", "wheel", factory.Options{"id": "3"}),
		factory.DefOptOfObjectRef("Wheel4", "wheel", factory.Options{"id": "4"}))

	myCar, err := carFactory.GetObject("mycar", factory.Options{"owner": "gogap"})

	if err != nil {
		fmt.Println(err)
		return
	}

	car := myCar.(*Car)

	car.Run()
}
