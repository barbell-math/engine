package main;

import (
    "fmt"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/settings"
)

func main(){
    settings.ReadSettings("./settings.json");
    test,err:=db.NewCRUD(settings.DBHost(),settings.DBPort(),settings.DBName());
    defer test.Close();
    if err!=nil {
        fmt.Println("Err connecting to DB");
        test.ResetDB();
    } else {
        fmt.Println("Connected to DB!");
    }
}
