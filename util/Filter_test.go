package util;

import (
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

var testVals []string=[]string{
    "one",
    "two",
    "three",
    "four",
    "five",
};

func filterExFunc(f Filter[string], v []string) map[string]interface{} {
    rv:=make(map[string]interface{},len(v));
    for _,iter:=range(v) {
        if f(iter) {
            rv[iter]=struct{}{};
        }
    }
    return rv;
}

func TestNoFilter(t *testing.T){
    fv:=filterExFunc(NoFilter[string],testVals);
    for _,iter:=range(testVals) {
        if _,ok:=fv[iter]; !ok {
            t.Errorf("Filter removed item it was not supposed to. | Item: %v",iter);
        }
    }
}

func TestAllFilter(t *testing.T){
    fv:=filterExFunc(AllFilter[string],testVals);
    testUtil.BasicTest(0, len(fv),"AllFilter did not remove all values.",t);
}

func TestCustomFilter(t *testing.T){
    fv:=filterExFunc(GenFilter(false,"one","three"),testVals);
    _,ok:=fv["one"];
    testUtil.BasicTest(true,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["two"];
    testUtil.BasicTest(false,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["three"];
    testUtil.BasicTest(true,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["four"];
    testUtil.BasicTest(false,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["five"];
    testUtil.BasicTest(false,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
}

func TestCustomFilterInverse(t *testing.T){
    fv:=filterExFunc(GenFilter(true,"one","three"),testVals);
    _,ok:=fv["one"];
    testUtil.BasicTest(false,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["two"];
    testUtil.BasicTest(true,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["three"];
    testUtil.BasicTest(false,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["four"];
    testUtil.BasicTest(true,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
    _,ok=fv["five"];
    testUtil.BasicTest(true,ok,
        "Value was not filtered when it was supposed to be.",t,
    );
}
