package templates

const DatabaseConnectionTemplate = `
package databases

import (
	"database/sql"
	"fmt"
	"strings"

	_ "{{.DatabaseDriver}}" // Import the database driver package
)

var {{.DBName}} *sql.DB


// DbQuery generates SQL queries (INSERT, UPDATE, DELETE) for a given table name and values.
func DbQuery( operation string, tableName string, values map[string]interface{}) string {
    switch operation {
    case "INSERT":
        return generateInsertQuery(tableName, values)
    case "UPDATE":
        return generateUpdateQuery(tableName, values)
    case "DELETE":
        return generateDeleteQuery(tableName)
	case "SELECTALL":
        return generateSelectAllQuery(tableName)
	case "SELECTBYID":
        return generateSelectByIDQuery(tableName, values)
    default:
        return ""
    }
}

func generateInsertQuery(tableName string, values map[string]interface{}) string {
    columns := []string{}
    args := []string{}

    for column, value := range values {
        columns = append(columns, column)
        args = append(args, value)

		if valueType == "string" {
			value = "'" + value.(string) + "'"
		}
		if valueType == "float64" {
			value = strconv.FormatFloat(value.(float64), 'f', -1, 64)
		}
		if valueType == "int64" {
			value = strconv.FormatInt(value.(int64), 10)
		}
		if valueType == "bool" {
			value = strconv.FormatBool(value.(bool))
		}
		args = append(args, value.(string))
    }

    return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(args, ", "))
}

func generateUpdateQuery(tableName string, values map[string]interface{}) string {
    updates := []string{}
    var args []interface{}

    for column, value := range values {
        updates = append(updates, fmt.Sprintf("%s = ?", column))
        args = append(args, value)
    }

    return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(updates, ", "))
}

func generateDeleteQuery(tableName string) string {
    // This is a simple example. You may want to add conditions based on values for a DELETE query.
    return fmt.Sprintf("DELETE FROM %s", tableName)
}

func generateSelectAllQuery(tableName string) string {
    // This is a simple example. You may want to customize the SELECT query based on the values provided.
    return fmt.Sprintf("SELECT * FROM %s", tableName)
}

func generateSelectByIDQuery(tableName string, values map[string]interface{}) string {
    // This is a simple example for retrieving a record by ID.
    return fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
}

func createDatabase() error {
	// Define the database connection parameters for creating the database
	dbHost := "{{.DBHost}}"
	dbPort := "{{.DBPort}}"
	dbUser := "{{.DBUser}}"
	dbPassword := "{{.DBPassword}}"
	dbName := "{{.DBName}}"

	// Connect to the database server without specifying a database name
	dbURL := fmt.Sprintf("{{.DBURLFormat}}", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("{{.DatabaseDriverName}}", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS {{.DBName}}")
	return err
}

func ConnectToDatabase() (*sql.DB, error) {
	// Define the database connection parameters
	dbHost := "{{.DBHost}}"
	dbPort := "{{.DBPort}}"
	dbName := "{{.DBName}}"
	dbUser := "{{.DBUser}}"
	dbPassword := "{{.DBPassword}}"
	
	// Construct the database connection URL
	dbURL := fmt.Sprintf("{{.DBURLFormat}}", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open a database connection
	db, err := sql.Open("{{.DatabaseDriverName}}", dbURL)
	if err != nil {
		return nil, err
	}

	// Test the database connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	{{.DBName}} = db

	

	return db, nil
}


`

// func main() {
// 	// Create the database
// 	err := createDatabase()
// 	if err != nil {
// 		fmt.Println("Failed to create the database:", err)
// 		return
// 	}
// 	fmt.Println("Database created successfully.")

// 	// Create and connect to the database
// 	db, err := connectToDatabase()
// 	if err != nil {
// 		fmt.Println("Failed to connect to the database:", err)
// 		return
// 	}
// 	defer db.Close()

// 	fmt.Println("Connected to the database successfully")

// 	// Now, you can use 'db' to perform database operations
// }
