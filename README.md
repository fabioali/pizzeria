# Pizza Maker

Let's bake some pizza. Customers order pizzas and wait for it.

 

## Making process

The sequential order of the pizza production process comprises 4 steps.

Each production step requires some time (in milliseconds) to finish.

A customer submits an order and waits for delivery. Once a pizza baker

received an order, he/she prepares the pizza, bakes it, and eventually

checks the quality of the baked product and delivers it.

 

- Receive and process order (1ms)

- Prepare pizza (2ms)

- Bake (5ms)

- Quality check (1ms)

 

I propose the following interfaces:

 

```go

type Order struct {

}

type Pizza struct {

    isBaked bool

}

type Oven interface {

    Bake(unbakedPizza Pizza) Pizza

}

type PizzaBaker interface {

    ProcessOrder() Order

    Prepare(order) *Pizza

    QualityCheck(pizza *Pizza) (*Pizza, error)

}

```

 

## Your Task

 

Implement a prototype and evaluate the performance (pizza throughput and

latency) of your Pizza maker implementation. Please document the 

resources you are using.

You can take as much time as you need to complete the implementation.

Please share your solution using some proper code sharing platform. Once

completed, we will setup a short review meeting and discuss next steps.

If you have any additional questions do not hesitate to contact me. 
