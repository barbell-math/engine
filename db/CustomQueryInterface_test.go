package db;

import (
    "testing"
)

func TestCustomReadQuery(t *testing.T){
    setup();
    settings.Modify(func(s *settings.Settings){
        s.DBInfo.DataVersion=1;
    });
    testDB.RunDataConversion();
    err:=CustomReadQuery(&testDB,"SELECT * FROM Exercise ORDER BY Id ASEC;",
        []any{},func(e *Exercise){

    });
}
