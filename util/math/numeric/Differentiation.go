package math;

func (f Function[N])Derivative(h Vars[N]) Function[N] {
    twoH:=make(map[string]N,len(h));
    for k,v:=range(h) {
        twoH[k]=2*v;
    }
    return func(iVars Vars[N], err error) (Vars[N],error) {
        return 
    }
}
