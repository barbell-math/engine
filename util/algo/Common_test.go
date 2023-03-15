package algo

import (
	"testing"
	"github.com/barbell-math/block/util/test"
	customerr "github.com/barbell-math/block/util/err"
)

func TestSlicesEqual(t *testing.T){
    test.BasicTest(true,SlicesEqual([]int{},[]int{}),
        "Slices equal returned false negative.",t,
    );
    test.BasicTest(true,SlicesEqual([]int{1},[]int{1}),
        "Slices equal returned false negative.",t,
    );
    test.BasicTest(true,SlicesEqual([]int{1,2,3,4},[]int{1,2,3,4}),
        "Slices equal returned false negative.",t,
    );
    test.BasicTest(false,SlicesEqual([]int{1},[]int{}),
        "Slices equal returned false positive.",t,
    );
    test.BasicTest(false,SlicesEqual([]int{1},[]int{2}),
        "Slices equal returned false positive.",t,
    );
    test.BasicTest(false,SlicesEqual([]int{1,1,1,1},[]int{1,1,1,2}),
        "Slices equal returned false positive.",t,
    );
}

func TestZipSlicesIncorrectDimensions(t *testing.T){
    rv,err:=ZipSlices([]int{1},[]int{});
    test.BasicTest(0, len(rv),"ZipSlices returned values it shouldn't have.",t);
    if !customerr.IsDimensionsDoNotAgree(err) {
        test.FormatError(customerr.DimensionsDoNotAgree(""),err,
            "ZipSlices returned the incorrect error.",t,
        );
    }
}

func TestZipSlices(t *testing.T){
    rv,err:=ZipSlices([]int{},[]int{});
    test.BasicTest(0, len(rv),"ZipSlices returned values it shouldn't have.",t);
    test.BasicTest(nil,err,
        "ZipSlices returned an error when it shouldn't have.",t,
    );
    res:=map[int]string{1: "one", 2: "two", 3: "three", 4: "four"};
    rv1,err:=ZipSlices([]int{1,2,3,4},[]string{"one","two","three","four"});
    test.BasicTest(nil,err,
        "ZipSlices returned an error when it shouldn't have.",t,
    );
    for k,v:=range(rv1) {
        realVal,present:=res[k];
        test.BasicTest(true,present,"ZipSlices added values to result",t);
        if present {
            test.BasicTest(realVal,v,
                "ZipSlices returned incorrectly ordered results.",t,
            );
        }
    }
}

func TestZipSlicesDuplicateKeys(t *testing.T){
    _,err:=ZipSlices([]int{1,1,3,4},[]string{"one","two","three","four"});
    if !IsSliceZippingError(err) {
        test.FormatError(SliceZippingError(""),err,
            "ZipSlices did not detect duplicate keys with correct error.",t,
        );
    }
}

func TestAppendWithPreallocation(t *testing.T){
    rv:=AppendWithPreallocation([]int{},[]int{});
    test.BasicTest(0,len(rv),
        "AppendWithPreallocation did not return all values.",t,
    );
    rv=AppendWithPreallocation([]int{1});
    test.BasicTest(1,len(rv),
        "AppendWithPreallocation did not return all values.",t,
    );
    test.BasicTest(1,rv[0],
        "AppendWithPreallocation returned incorrect value.",t,
    );
    rv=AppendWithPreallocation([]int{},[]int{1},[]int{2,3},[]int{4,5,6,7});
    test.BasicTest(7,len(rv),
        "AppendWithPreallocation did not return all values.",t,
    );
    for i,v:=range(rv) {
        test.BasicTest(i+1,v,
            "AppendWithPreallocation returned incorrect values.",t,
        );
    }
}
