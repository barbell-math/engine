package db;

import (
    "testing"
    "github.com/barbell-math/block/util/test"
)

func TestQueryToString(t *testing.T){
    test.BasicTest("SELECT",SelectStmt.String(),
        "Select string is not correct.",t,
    );
    test.BasicTest("UPDATE",UpdateStmt.String(),
        "Update string is not correct.",t,
    );
    test.BasicTest("DELETE",DeleteStmt.String(),
        "Delete string is not correct.",t,
    );
    test.BasicTest("INSERT",InsertStmt.String(),
        "Insert string is not correct.",t,
    );
    test.BasicTest("unknown",UnknownStmt.String(),
        "Insert string is not correct.",t,
    );
}

func TestIsQueryType(t *testing.T){
    test.BasicTest(true,SelectStmt.isQueryType("SELECT * FROM Table;"),
        "Checking query type returned false negative.",t,
    );
    test.BasicTest(true,SelectStmt.isQueryType("  \tSELECT * FROM Table;"),
        "Checking query type returned false negative.",t,
    );
    test.BasicTest(true,SelectStmt.isQueryType("select * from Table;"),
        "Checking query type returned false negative.",t,
    );
    test.BasicTest(false,SelectStmt.isQueryType("UPDATE * FROM Table;"),
        "Checking query type returned false positive.",t,
    );
    test.BasicTest(false,SelectStmt.isQueryType("* FROM Table;"),
        "Checking query type returned false positive.",t,
    );
    test.BasicTest(false,SelectStmt.isQueryType(""),
        "Checking query type returned false positive.",t,
    );
}
