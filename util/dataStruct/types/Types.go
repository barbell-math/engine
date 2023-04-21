package types;

//This file should never import anything other than the std library. If anything
// else is imported the risk of import loops is very high.

type Vector[T any] interface {
    Get(idx int) (T,error);
    GetPntr(idx int) (*T,error);
    Set(idx int) (T,error);
    Append(v T) error;
    Insert(v T, idx int) error;
    Length() int;
};

type Queue[T any] interface {
    Pop() (T,error);
    Peek(idx int) (T,error);
    PeekPntr(idx int) (*T,error);
    Push(v T) (error);
    Capacity() int;
    Length() int;
};
