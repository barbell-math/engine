package db;

import (
    "strings"
)

type QueryType int;
const (
    UnknownStmt QueryType = iota
    SelectStmt
    UpdateStmt
    DeleteStmt
    InsertStmt
)

func (q QueryType)String() string {
    switch q {
        case SelectStmt: return "SELECT";
        case UpdateStmt: return "UPDATE";
        case DeleteStmt: return "DELETE";
        case InsertStmt: return "INSERT";
        case UnknownStmt: fallthrough
        default: return "unknown";
    }
}

func (q QueryType)isQueryType(sqlStmt string) bool {
    tmp:=strings.TrimSpace(sqlStmt);
    prefix:=q.String();
    if strings.HasPrefix(tmp,prefix) ||
       strings.HasPrefix(tmp,strings.ToLower(prefix)) {
        return true;
    }
    return false;
}
