package db;

import (
    "github.com/carmichaeljr/powerlifting-engine/util"
)

const CURRENT_DATA_VERSION int=1;

type DataVersionConversion func(crud *CRUD) error;
var DataVersionOps map[int]DataVersionConversion=map[int]DataVersionConversion{
    1: func (crud *CRUD) error {
        crud.ResetDB();
        crud.addDataVersion(1);
        typeMap:=make(map[string]int);
        focusMap:=make(map[string]int);
        util.CSVFileSplitter("./sql/ExerciseTypesInit.csv",',',true,func(columns []string){
            typeMap[columns[0]],_=crud.CreateExerciseType(
                ExerciseType{
                    _type: columns[0],
                    description: columns[1],
                },
            );
        });
        util.CSVFileSplitter("./sql/ExerciseFocusInit.csv",',',true,func(columns []string){
            focusMap[columns[0]],_=crud.CreateExerciseFocus(
                ExerciseFocus{
                    focus: columns[0],
                },
            );
        });
        util.CSVFileSplitter("./sql/ExercisesInit.csv",',',true,func(columns []string){
            crud.CreateExercise(
                Exercise{
                    name: columns[0],
                    typeID: typeMap[columns[1]],
                    focusID: focusMap[columns[2]],
                },
            );
        });
        return nil;
    },
};
