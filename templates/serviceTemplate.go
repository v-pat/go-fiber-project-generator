package templates

const ServiceTemplate = `package service

import (
    "encoding/json"
    "fmt"
    "{{.AppName}}/databases"
    "{{.AppName}}/model"
)

{{.StructCode}}

// Create{{.StructName}} inserts a new {{.StructName}} record into the database.
func  Create{{.StructName}}({{.StructName}}s model.{{.StructNameTitlecase}}) error {
    var {{.StructName}}sMap map[string]interface{}
    {{.StructName}}, _ := json.Marshal({{.StructName}}s)
    json.Unmarshal({{.StructName}}, &{{.StructName}}sMap)
    query := databases.DbQuery("INSERT", "{{.StructName}}",{{.StructName}}sMap)
    // Execute the query to insert {{.StructName}} into the database
    fmt.Println("Executing query:", query)
    _, err := databases.{{.DBName}}.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Get{{.StructName}} retrieves a {{.StructName}} record from the database by ID.
func Get{{.StructName}}ByID(id int) (model.{{.StructNameTitlecase}}, error) {
    query := databases.DbQuery("SELECTBYID", "{{.StructName}}",map[string]interface{}{"id":id})
    // Execute the query to retrieve {{.StructName}} from the database
    fmt.Println("Executing query:", query)
    {{.StructName}} := model.{{.StructNameTitlecase}}{} // Replace with actual retrieval logic
    res, err := databases.{{.DBName}}.Query(query)
    if err != nil {
        return {{.StructName}},err
    }
    // Implement query execution and scanning here
    res.Scan({{.StructName}})
    return {{.StructName}}, nil
}

// Update{{.StructName}} updates an existing {{.StructName}} record in the database.
func Update{{.StructName}}({{.StructName}}s model.{{.StructNameTitlecase}}) error {
    var {{.StructName}}sMap map[string]interface{}
    {{.StructName}}, _ := json.Marshal({{.StructName}}s)
    json.Unmarshal({{.StructName}}, &{{.StructName}}sMap)
    query := databases.DbQuery("{UPDATE","{{.StructName}}s", {{.StructName}}sMap)
    // Execute the query to update {{.StructName}} in the database
    fmt.Println("Executing query:", query)
    _, err := databases.{{.DBName}}.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Delete{{.StructName}} deletes a {{.StructName}} record from the database by ID.
func Delete{{.StructName}}ByID(id int) error {
    query := databases.DbQuery("DELETE", "{{.StructName}}", map[string]interface{}{"id":id})
    // Execute the query to delete {{.StructName}} from the database
    fmt.Println("Executing query:", query)
    _, err := databases.{{.DBName}}.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

`
