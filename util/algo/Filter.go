package algo;

type Filter[V comparable] func(thing V) bool;

func NoFilter[V comparable](thing V) bool { return true; }
func AllFilter[V comparable](thing V) bool { return false; }

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
