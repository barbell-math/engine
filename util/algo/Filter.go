package algo;

type Filter[V any] func(thing V) bool;

func NoFilter[V any](thing V) bool { return true; }
func AllFilter[V any](thing V) bool { return false; }

func GenFilter[V comparable](inverse bool, things ...V) Filter[V] {
    return func(thing V) bool {
        rv:=inverse;
        for i:=0; ((inverse && rv) || (!inverse && !rv)) && i<len(things); i++ {
            if inverse {
                rv=(rv && thing!=things[i]);
            } else {
                rv=(thing==things[i]);
            }
        }
        return rv;
    }
}

func NoNil[P any, V *P](thing V) bool { return thing!=nil; }
func NoNilError(err error) bool { return err!=nil; }
