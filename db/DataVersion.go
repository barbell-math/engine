package db;

import (
    "github.com/carmichaeljr/powerlifting-engine/util"
)

type DataVersionConversion func(crud *CRUD) error;
var DataVersionOps map[int]DataVersionConversion=map[int]DataVersionConversion{
    1: zeroToOne,
};

func zeroToOne(crud *CRUD) error {
    crud.ResetDB();
    crud.setDataVersion(1);
    var err error=nil;
    typeMap, focusMap:=make(map[string]int), make(map[string]int);
    util.CSVFileSplitter("./sql/ExerciseTypesInit.csv",',',true,
        func(columns []string) bool {
            tmp,err:=Create(crud,ExerciseType{
                T: columns[0], Description: columns[1],
            });
            typeMap[columns[0]]=tmp[0];
            return err==nil;
        },
    );
    util.CSVFileSplitter("./sql/ExerciseFocusInit.csv",',',true,
        func(columns []string) bool {
            tmp,err:=Create(crud,ExerciseFocus{Focus: columns[0]});
            focusMap[columns[0]]=tmp[0];
            return err==nil;
        },
    );
    util.CSVFileSplitter("./sql/ExercisesInit.csv",',',true,
        func(columns []string) bool {
            _,err:=Create(crud,Exercise{
                Name: columns[0],
                TypeID: typeMap[columns[1]],
                FocusID: focusMap[columns[2]],
            });
            return err==nil;
        },
    );
    return err;
}
