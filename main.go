package main;

import (
    "fmt"
    "github.com/carmichaeljr/powerlifting-engine/db"
)

func main(){
    test,err:=db.NewCRUD("localhost",5432,"carmichaeljr","research");
    defer test.Close();
    if err!=nil {
        panic(err);
        fmt.Println("Err connecting to DB");
    } else {
        fmt.Println("Connected to DB!");
    }
}
