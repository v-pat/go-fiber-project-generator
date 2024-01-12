package service

import (
    "encoding/json"
    "fmt"
    "vpat/databases"
    "vpat/model"
)

type tab2 struct {
	number float64 `json:"number"`
	surname string `json:"surname"`
	job string `json:"job"`
}


// Createtab2 inserts a new tab2 record into the database.
func  Createtab2(tab2s model.Tab2) error {
    var tab2sMap map[string]interface{}
    tab2, _ := json.Marshal(tab2s)
    json.Unmarshal(tab2, &tab2sMap)
    query := databases.DbQuery("INSERT", "tab2",tab2sMap)
    // Execute the query to insert tab2 into the database
    fmt.Println("Executing query:", query)
    _, err := databases.Vpat.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Gettab2 retrieves a tab2 record from the database by ID.
func Gettab2ByID(id int) (model.Tab2, error) {
    query := databases.DbQuery("SELECTBYID", "tab2",map[string]interface{}{"id":id})
    // Execute the query to retrieve tab2 from the database
    fmt.Println("Executing query:", query)
    tab2 := model.Tab2{} // Replace with actual retrieval logic
    res, err := databases.Vpat.Query(query)
    if err != nil {
        return tab2,err
    }
    // Implement query execution and scanning here
    res.Scan(tab2)
    return tab2, nil
}

// Updatetab2 updates an existing tab2 record in the database.
func Updatetab2(tab2s model.Tab2) error {
    var tab2sMap map[string]interface{}
    tab2, _ := json.Marshal(tab2s)
    json.Unmarshal(tab2, &tab2sMap)
    query := databases.DbQuery("{UPDATE","tab2s", tab2sMap)
    // Execute the query to update tab2 in the database
    fmt.Println("Executing query:", query)
    _, err := databases.Vpat.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Deletetab2 deletes a tab2 record from the database by ID.
func Deletetab2ByID(id int) error {
    query := databases.DbQuery("DELETE", "tab2", map[string]interface{}{"id":id})
    // Execute the query to delete tab2 from the database
    fmt.Println("Executing query:", query)
    _, err := databases.Vpat.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

