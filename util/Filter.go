package util;

type Filter[V comparable] func(name V) bool;

func NoFilter[V comparable](name V) bool { return true; }
func AllFilter[V comparable](name V) bool { return false; }

func GenFilter[V comparable](inverse bool, names ...V) Filter[V] {
    return func(name V) bool {
        rv:=inverse;
        for i:=0; ((inverse && rv) || (!inverse && !rv)) && i<len(names); i++ {
            if inverse {
                rv=(rv && name!=names[i]);
            } else {
                rv=(name==names[i]);
            }
        }
        return rv;
    }
}
