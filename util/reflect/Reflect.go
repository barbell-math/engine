package reflect;

import (
    "reflect"
)

func GetErrorFromReflectValue(in *reflect.Value) error {
    switch in.Interface().(type) {
        case error: return in.Interface().(error);
        default: return nil;
    }
}
