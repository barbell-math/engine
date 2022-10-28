package db;

import (
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

func TestQueryToString(t *testing.T){
    testUtil.BasicTest("SELECT",SelectStmt.String(),
        "Select string is not correct.",t,
    );
    testUtil.BasicTest("UPDATE",UpdateStmt.String(),
        "Update string is not correct.",t,
    );
    testUtil.BasicTest("DELETE",DeleteStmt.String(),
        "Delete string is not correct.",t,
    );
    testUtil.BasicTest("INSERT",InsertStmt.String(),
        "Insert string is not correct.",t,
    );
    testUtil.BasicTest("unknown",UnknownStmt.String(),
        "Insert string is not correct.",t,
    );
}

func TestIsQueryType(t *testing.T){
    testUtil.BasicTest(true,SelectStmt.isQueryType("SELECT * FROM Table;"),
        "Checking query type returned false negative.",t,
    );
    testUtil.BasicTest(true,SelectStmt.isQueryType("  \tSELECT * FROM Table;"),
        "Checking query type returned false negative.",t,
    );
    testUtil.BasicTest(true,SelectStmt.isQueryType("select * from Table;"),
        "Checking query type returned false negative.",t,
    );
    testUtil.BasicTest(false,SelectStmt.isQueryType("UPDATE * FROM Table;"),
        "Checking query type returned false positive.",t,
    );
    testUtil.BasicTest(false,SelectStmt.isQueryType("* FROM Table;"),
        "Checking query type returned false positive.",t,
    );
    testUtil.BasicTest(false,SelectStmt.isQueryType(""),
        "Checking query type returned false positive.",t,
    );
}
