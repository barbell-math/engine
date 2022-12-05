package dataStruct;

import (
    "fmt"
)

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

func ZipSlices[K comparable, V any](keys []K, vals []V) (map[K]V,error) {
    rv:=make(map[K]V,len(keys));
    if len(keys)==len(vals) {
        for i,k:=range(keys) {
            if _,ok:=rv[k]; !ok {
                rv[k]=vals[i];
            } else {
                return rv,SliceZippingError(fmt.Sprintf(
                    "Keys have duplicate values | %v",k,
                ));
            }
        }
    } else {
        return rv,SliceZippingError(fmt.Sprintf(
            "Lengths are not equal. | K: %d V: %d",len(keys),len(vals),
        ));
    }
    return rv,nil;
}
