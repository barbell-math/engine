package iter;


type IteratorFeedback int;
const (
    Continue IteratorFeedback=iota
    Break
    Iterate
);

type Iter[T any] func(f IteratorFeedback)(T,error,bool);

type Window[T any] struct {
    size int;
    vals []T;
    zeroIndex int;
};

func NewWindow[T any](size int) Window[T] {
    return Window[T]{
        size: size,
        zeroIndex: 0,
        vals: make([]T,size),
    };
}

func (w *Window[T])Add(v T){
    if len(w.vals)<w.size {
        w.vals=append(w.vals,v);
    }
}
