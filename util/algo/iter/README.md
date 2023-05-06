# iter

A generic lazy iterator framework library that respects errors as values with reverse message passing for resource management.

```golang
func sequenceGenerator(start float64, end float64, step float64) Iter[float64] {
    cntr:=start;
    return func(f IteratorFeedback) (float64,error,bool) {
        if f!=Break && cntr<=end {
            cntr+=step;
            return cntr-1,nil,true;
        }
        return end,nil,false;
    }
}

//Find the area under the function y=5cos(2pi/5(x-5))+5 between [-100,100]
// using a Riemann sum with a given step size
val,err:=sequenceGenerator(-100.0,100.0,step).Map(func(index int, val float64) (float64,error) {
    height:=amp*math.Cos(2*math.Pi/period*(val-hShift))+vShift;
    return height*step,nil;
}).Reduce(0.0, func(accum *float64, iter float64) error {
    *accum+=iter;
    return nil;
});
fmt.Printf("Area is: %f Using step size: %f\n",total,step);
fmt.Printf("Err is: %v\n",err);

//Output:
//Area is 1000.000000 Using step size: 0.000100
//Err is: <nil>
```

## Design

This library takes advantage of the fact that functions can have methods in go-lang. This allows an iterator to be defined as a function and for its methods to call on there relative 'object' (which would be a function in this case) to get values as needed. Each method returns a new iterator, allowing the methods to be chained together, creating a recursive, lazily evaluated iterator sequence.

There are three parts to any iterator chain, with the middle part being optional:

1. Producer: The producer is responsible for creating a stream of values to pass to the rest os the iterator stream. The source can be a slice, channel, or single value.
1. Intermediary: An Intermediary is responsible for taking it's parent iterators values, mutating and/or filtering them, and passing them down to it's child iterator.
1. Consumer: The consumer is what collects all of the final values and either saves them or performs some aggregate function with them.

The intermediaries and consumers can be further sub-categorized:

1. Non-Pseudo: All iterators in the parent category (i.e. Intermediary or consumer) can be expressed using a non-pseudo iterator. For intermediaries the non-pseudo iterator is ```Next``` and for consumers the non-pseudo iterators are ```ForEach``` and ```Stop```.
1. Pseudo: Any iterator in the parent category that can be expressed using the appropriate categories non-pseudo iterator.

If you are looking to extend this library and add more iterators, it is highly recommended that any new intermediary or consumer iterators are created _using the non-pseudo iterators_. This will reduce errors and time spent needlessly banging your head against a wall.

## Reverse Message Passing

The iterators can be in three possible states:

1. Continue: Signaling to 'accept' the current value and pass it along to the child iterator.
1. Break: Signaling to ignore the current value and return the signal to stop iterating.
1. Iterate: Signaling the current iterator to continue iterating and grab the next value.

These states are managed for each individual iterator, and are passed between the child and parent iterators.

If any iterator produces an error or signals its child iterator to stop iterating, a ```Break``` command is generated. This ```Break``` command is passed down to the consumer with no action being taken by the intermediaries. The consumer then calls its parent iterator with the ```Break``` command which is propagated all the way to the producer, which performs it's resource management. Once done, the producer returns any errors and it's child iterator is allowed to perform resource management and the pattern continues all the way down to the consumer. This allows for resources to be properly destroyed in a top-down fashion.

## Benchmarking

Obviously there will be overhead when using this library instead of using plain for loops. The ```example_test.go``` not only showcases the example at the top of this readme, but contains benchmarks for three different scenarios. These scenarios are shown below for convince.


##### Scenario 1: A 'typical' functional implementation

```golang
val,err:=sequenceGenerator(-100.0,100.0,step).Map(func(index int, val float64) (float64,error) {
    height:=amp*math.Cos(2*math.Pi/period*(val-hShift))+vShift;
    return height*step,nil;
}).Reduce(0.0, func(accum *float64, iter float64) error {
    *accum+=iter;
    return nil;
});
```

##### Scenario 2: Another implementation using iterators

```golang
total:=0.0;
err:=sequenceGenerator(-100.0,100.0,step).ForEach(func(index int, val float64) (IteratorFeedback,error) {
    height:=amp*math.Cos(2*math.Pi/period*(val-hShift))+vShift;
    total+=height*step;
    return Continue,nil;
});
```

##### Scenario 3: A basic for loop

```golang
total:=0.0;
for x:=-100.0; x<=100; x+=step {
    height:=amp*math.Cos(2*math.Pi/period*(x-hShift))+vShift;
    total+=height*step;
}
```

The benchmarks (gathered from the go-lang benchmark utility) for the scenarios with various step sizes are shown below. Make of these results as you will.

| Scenario | Step Size | Time |
|----------|-----------|------|
| 1 | 1 | 7815 ns/op |
| 2 | 1 | 5930 ns/op |
| 3 | 1 | 3706 ns/op |
| 1 | 0.1 | 84634 ns/op |
| 2 | 0.1 | 59650 ns/op |
| 3 | 0.1 | 34406 ns/op |
| 1 | 0.01 | 765754 ns/op |
| 2 | 0.01 | 576928 ns/op |
| 3 | 0.01 | 352233 ns/op |
| 1 | 0.001 | 7526192 ns/op |
| 2 | 0.001 | 5810898 ns/op |
| 3 | 0.001 | 3460169 ns/op |
| 1 | 0.0001 | 73463321 ns/op |
| 2 | 0.0001 | 57991047 ns/op |
| 3 | 0.0001 | 34179406 ns/op |
