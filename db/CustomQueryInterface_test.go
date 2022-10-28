package db;

import (
    "testing"
)

func TestCustomReadQuery(t *testing.T){
    setup();
    CustomReadQuery(&testDB,"SELECT * FROM Exercise ORDER BY Id DESC;",
        []any{},func(e *Exercise){

    });
}
