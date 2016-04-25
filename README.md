Factory
=======

A object factory like java's spring bean factory

### Features
1. Dependency Inject
2. Dynamic create object instance


### Usage

#### Define your model
```go
type Wheel struct {
	ID string
}

type Car struct {
	owner string

	Wheel1 *Wheel
	Wheel2 *Wheel
	Wheel3 *Wheel
	Wheel4 *Wheel
}
```

####  Register your model
```go
factory.RegisterModel((*Wheel)(nil), "Michelin")
factory.RegisterModel((*Car)(nil), "Skoda")
```

#### Build factory
```go
carFactory := factory.NewClassicFactory(nil)

carFactory.Define("wheel", factory.Prototype, "Michelin", factory.DefOptOfNewObjectFunc(NewWheel))

	err = carFactory.Define("mycar", factory.Prototype, "Skoda",
		factory.DefOptOfNewObjectFunc(NewCar),
		factory.DefOptOfObjectRef("Wheel1", "wheel", factory.Options{"id": "1"}),
		factory.DefOptOfObjectRef("Wheel2", "wheel", factory.Options{"id": "2"}),
		factory.DefOptOfObjectRef("Wheel3", "wheel", factory.Options{"id": "3"}),
		factory.DefOptOfObjectRef("Wheel4", "wheel", factory.Options{"id": "4"}),
	)
```

The func `NewWheel` and `NewCar` is your user-define new func, it should define as 

```go
func(opts Options) (v interface{}, err error)
```

you could get your option from `Options`

The `car` depend on `wheel`, so we use `factory.DefOptOfObjectRef` for configurate model ref, the first param is filed of `Car`'s name, the Ref object initial order is params order, you also could use `factory.DefOptOfRefOrder` to define ref object initial order.


### Get object

```go
var myCar interface{}
myCar, err = carFactory.GetObject("mycar", factory.Options{"owner": "GoGap"})

if err != nil {
	fmt.Println(err)
	return
}

car := myCar.(*Car)
```

### Example

```go
type Hub struct {
	ID string
}

func NewHub(opts factory.Options) (hub interface{}, err error) {
	h := &Hub{}
	opts.Get("id", &h.ID)
	hub = h
	return
}

type Wheel struct {
	ID string

	Hub *Hub
}

func NewWheel(opts factory.Options) (wheel interface{}, err error) {
	w := &Wheel{}
	opts.Get("id", &w.ID)
	wheel = w
	return
}

func (p *Wheel) Run() {
	fmt.Printf("Wheel Running, ID: %s, HubID: %s\n", p.ID, p.Hub.ID)
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
	factory.RegisterModel((*Hub)(nil), "BBS")
	factory.RegisterModel((*Wheel)(nil), "Michelin")
	factory.RegisterModel((*Car)(nil), "Skoda")
}

func main() {

	var err error

	carFactory := factory.NewClassicFactory(nil)

	carFactory.Define("hub", factory.Prototype, "BBS", factory.DefOptOfNewObjectFunc(NewHub))

	carFactory.Define("wheel", factory.Prototype, "Michelin", factory.DefOptOfNewObjectFunc(NewWheel))

	err = carFactory.Define("mycar", factory.Prototype, "Skoda",
		factory.DefOptOfNewObjectFunc(NewCar),
		factory.DefOptOfObjectRef("Wheel1", "wheel", factory.Options{"id": "1"}),
		factory.DefOptOfObjectRef("Wheel2", "wheel", factory.Options{"id": "2"}),
		factory.DefOptOfObjectRef("Wheel3", "wheel", factory.Options{"id": "3"}),
		factory.DefOptOfObjectRef("Wheel4", "wheel", factory.Options{"id": "4"}),
		factory.DefOptOfObjectRef("Wheel1.Hub", "hub", factory.Options{"id": "HUB01"}),
		factory.DefOptOfObjectRef("Wheel2.Hub", "hub", factory.Options{"id": "HUB02"}),
		factory.DefOptOfObjectRef("Wheel3.Hub", "hub", factory.Options{"id": "HUB03"}),
		factory.DefOptOfObjectRef("Wheel4.Hub", "hub", factory.Options{"id": "HUB04"}),
	)

	if err != nil {
		return
	}

	var myCar interface{}
	myCar, err = carFactory.GetObject("mycar", factory.Options{"owner": "GoGap"})

	if err != nil {
		fmt.Println(err)
		return
	}

	car := myCar.(*Car)

	car.Run()
}

```

#### Output
```bash
> go run main.go
Wheel Running, ID: 1, HubID: HUB01
Wheel Running, ID: 2, HubID: HUB02
Wheel Running, ID: 3, HubID: HUB03
Wheel Running, ID: 4, HubID: HUB04
GoGap' Car Running
```