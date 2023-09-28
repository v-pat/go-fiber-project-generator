package templates

const ServiceTemplate = `package service

import (
    "database/sql"
    "fmt"
    "{{.AppName}}/model"
)

// {{.StructName}}Service represents the service for {{.StructName}}.
type {{.StructName}}Service struct {
    DB *sql.DB
}

{{.StructCode}}

// Create{{.StructName}} inserts a new {{.StructName}} record into the database.
func (s *{{.StructName}}Service) Create{{.StructName}}({{.StructName}} *model.{{.StructName}}) error {
    query := dbQuery("{{.Database}}", "insert", "{{.StructName}}s")
    // Execute the query to insert {{.StructName}} into the database
    fmt.Println("Executing query:", query)
    _, err := s.DB.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Get{{.StructName}} retrieves a {{.StructName}} record from the database by ID.
func (s *{{.StructName}}Service) Get{{.StructName}}(id int) (*model.{{.StructName}}, error) {
    query := dbQuery("{{.Database}}", "select", "{{.StructName}}s")
    // Execute the query to retrieve {{.StructName}} from the database
    fmt.Println("Executing query:", query)
    // Implement query execution and scanning here
    {{.StructName}} := &model.{{.StructName}}{} // Replace with actual retrieval logic
    return {{.StructName}}, nil
}

// Update{{.StructName}} updates an existing {{.StructName}} record in the database.
func (s *{{.StructName}}Service) Update{{.StructName}}({{.StructName}} *model.{{.StructName}}) error {
    query := dbQuery("{{.Database}}", "update", "{{.StructName}}s")
    // Execute the query to update {{.StructName}} in the database
    fmt.Println("Executing query:", query)
    _, err := s.DB.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Delete{{.StructName}} deletes a {{.StructName}} record from the database by ID.
func (s *{{.StructName}}Service) Delete{{.StructName}}(id int) error {
    query := dbQuery("{{.Database}}", "delete", "{{.StructName}}s")
    // Execute the query to delete {{.StructName}} from the database
    fmt.Println("Executing query:", query)
    _, err := s.DB.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

`
