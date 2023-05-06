package dataStruct

import (
	"testing"

	"github.com/barbell-math/block/util/dataStruct/types"
	"github.com/barbell-math/block/util/test"
)

func TestVariantA(t *testing.T){
    var tmp int=5;
    v:=Variant[int,float64]{};
    newV:=v.SetValA(tmp);
    test.BasicTest(true,newV.HasA(),"Variant did claim to have correct value.",t);
    test.BasicTest(false,newV.HasB(),"Variant did claim to have correct value.",t);
    test.BasicTest(5,newV.ValA(),"Variant did not return correct value.",t);
    test.BasicTest(5,newV.ValAOr(1),"Variant did not return correct value.",t);
    test.BasicTest(float64(1.0),newV.ValBOr(1),
        "Variant did not return correct value.",t,
    );
}

func TestVariantB(t *testing.T){
    var tmp float64=5;
    v:=Variant[int,float64]{};
    newV:=v.SetValB(tmp);
    test.BasicTest(false,newV.HasA(),"Variant did claim to have correct value.",t);
    test.BasicTest(true,newV.HasB(),"Variant did claim to have correct value.",t);
    test.BasicTest(float64(5),newV.ValB(),"Variant did not return correct value.",t);
    test.BasicTest(1,newV.ValAOr(1),"Variant did not return correct value.",t);
    test.BasicTest(float64(5),newV.ValBOr(1),
        "Variant did not return correct value.",t,
    );
}

func interfaceTestHelper[T any, U any](v types.Variant[T,U]){}
func TestVariantInterface(t *testing.T){
    tmp:=5;
    v:=Variant[int,float64]{};
    interfaceTestHelper[int,float64](v);
    interfaceV:=v.SetValA(tmp);
    interfaceTestHelper(interfaceV);
}
