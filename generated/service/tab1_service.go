package service

import (
    "encoding/json"
    "fmt"
    "vpat/databases"
    "vpat/model"
)

type tab1 struct {
	id float64 `json:"id"`
	name string `json:"name"`
	role string `json:"role"`
}


// Createtab1 inserts a new tab1 record into the database.
func  Createtab1(tab1s model.Tab1) error {
    var tab1sMap map[string]interface{}
    tab1, _ := json.Marshal(tab1s)
    json.Unmarshal(tab1, &tab1sMap)
    query := databases.DbQuery("INSERT", "tab1",tab1sMap)
    // Execute the query to insert tab1 into the database
    fmt.Println("Executing query:", query)
    _, err := databases.Vpat.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Gettab1 retrieves a tab1 record from the database by ID.
func Gettab1ByID(id int) (model.Tab1, error) {
    query := databases.DbQuery("SELECTBYID", "tab1",map[string]interface{}{"id":id})
    // Execute the query to retrieve tab1 from the database
    fmt.Println("Executing query:", query)
    tab1 := model.Tab1{} // Replace with actual retrieval logic
    res, err := databases.Vpat.Query(query)
    if err != nil {
        return tab1,err
    }
    // Implement query execution and scanning here
    res.Scan(tab1)
    return tab1, nil
}

// Updatetab1 updates an existing tab1 record in the database.
func Updatetab1(tab1s model.Tab1) error {
    var tab1sMap map[string]interface{}
    tab1, _ := json.Marshal(tab1s)
    json.Unmarshal(tab1, &tab1sMap)
    query := databases.DbQuery("{UPDATE","tab1s", tab1sMap)
    // Execute the query to update tab1 in the database
    fmt.Println("Executing query:", query)
    _, err := databases.Vpat.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Deletetab1 deletes a tab1 record from the database by ID.
func Deletetab1ByID(id int) error {
    query := databases.DbQuery("DELETE", "tab1", map[string]interface{}{"id":id})
    // Execute the query to delete tab1 from the database
    fmt.Println("Executing query:", query)
    _, err := databases.Vpat.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

