package db

import (
	"github.com/barbell-math/engine/settings"
	"github.com/barbell-math/engine/util/algo/iter"
	"github.com/barbell-math/engine/util/io/csv"
)

type DataVersionConversion func(crud *DB) error;
var DataVersionOps map[int]DataVersionConversion=map[int]DataVersionConversion{
    1: zeroToOne,
};

func zeroToOne(crud *DB) error {
    crud.ResetDB();
    crud.setDataVersion(1);
    var err error=nil;
    typeMap, focusMap:=make(map[string]int), make(map[string]int);
    if fError:=csv.CSVFileSplitter(settings.ExerciseTypeInitData(),',','#').ForEach(
    func(index int, val []string) (iter.IteratorFeedback, error) {
        tmp,err:=Create(crud,ExerciseType{
            T: val[0], Description: val[1],
        });
        typeMap[val[0]]=tmp[0];
        return iter.Continue,err;
    }); err!=nil {
        return fError;
    }
    //if fError:=csv.CSVFileSplitter(settings.ExerciseTypeInitData(),',',true,
    //    func(columns []string) bool {
    //        tmp,err:=Create(crud,ExerciseType{
    //            T: columns[0], Description: columns[1],
    //        });
    //        typeMap[columns[0]]=tmp[0];
    //        return err==nil;
    //    }); fError!=nil {
    //    return fError;
    //}
    if fError:=csv.CSVFileSplitter(settings.ExerciseTypeInitData(),',','#').ForEach(
    func(index int, val []string) (iter.IteratorFeedback, error) {
         tmp,err:=Create(crud,ExerciseFocus{Focus: val[0]});
         focusMap[val[0]]=tmp[0];
         return iter.Continue,err;
    }); err!=nil {
         return fError;
    }
    //if fError:=csv.CSVFileSplitter(settings.ExerciseFocusInitData(),',',true,
    //    func(columns []string) bool {
    //        tmp,err:=Create(crud,ExerciseFocus{Focus: columns[0]});
    //        focusMap[columns[0]]=tmp[0];
    //        return err==nil;
    //    }); fError!=nil {
    //    return fError;
    //}
    if fError:=csv.CSVFileSplitter(settings.ExerciseTypeInitData(),',','#').ForEach(
    func(index int, val []string) (iter.IteratorFeedback, error) {
        _,err:=Create(crud,Exercise{
            Name: val[0],
            TypeID: typeMap[val[1]],
            FocusID: focusMap[val[2]],
        });
        return iter.Continue,err;
    }); err!=nil {
        return fError;
    }
    //if fError:=csv.CSVFileSplitter(settings.ExerciseInitData(),',',true,
    //    func(columns []string) bool {
    //        _,err:=Create(crud,Exercise{
    //            Name: columns[0],
    //            TypeID: typeMap[columns[1]],
    //            FocusID: focusMap[columns[2]],
    //        });
    //        return err==nil;
    //    }); fError!=nil {
    //    return fError;
    //}
    return err;
}
