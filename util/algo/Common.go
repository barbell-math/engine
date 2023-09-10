package algo;

import (
	"fmt"
	"github.com/barbell-math/engine/util/algo/iter"
	customerr "github.com/barbell-math/engine/util/err"
)

func SlicesEqual[T comparable](one []T, two []T) bool {
    if len(one)!=len(two) {
        return false;
    }
    i:=-1;
    res,_:=iter.SliceElems(one).All(func(other T) (bool,error) {
        i++;
        return other==two[i],nil;
    });
    return res;
}

func ZipSlices[K comparable, V any](keys []K, vals []V) (map[K]V,error) {
    rv:=make(map[K]V,len(keys));
    if err:=customerr.ArrayDimsArgree(keys,vals,
        "Keys and values are different lengths.",
    ); err!=nil {
        return rv,err;
    }
    for i,k:=range(keys) {
        if _,ok:=rv[k]; !ok {
            rv[k]=vals[i];
        } else {
            return rv,SliceZippingError(fmt.Sprintf(
                "Keys have duplicate values | %v",k,
            ));
        }
    }
    return rv,nil;
}

func AppendWithPreallocation[T any](slices ...[]T) []T {
    var totLen,i int=0, 0;
    for _,s:=range(slices) {
        totLen+=len(s);
    }
    rv:=make([]T,totLen);
    for _,tmp:=range(slices){
        i+=copy(rv[i:],tmp);
    }
    return rv;
}
