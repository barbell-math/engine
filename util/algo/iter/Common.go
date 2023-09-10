package iter;


type IteratorFeedback int;
const (
    Continue IteratorFeedback=iota
    Break
    Iterate
);

type Iter[T any] func(f IteratorFeedback)(T,error,bool);
