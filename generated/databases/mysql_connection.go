
package databases

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Import the database driver package
)

var Vpat *sql.DB


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
    placeholders := []string{}
    var args []interface{}

    for column, value := range values {
        columns = append(columns, column)
        placeholders = append(placeholders, "?")
        args = append(args, value)
    }

    return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
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
	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPassword := "root"
	dbName := "Vpat"

	// Connect to the database server without specifying a database name
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS Vpat")
	return err
}

func ConnectToDatabase() (*sql.DB, error) {
	// Define the database connection parameters
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "Vpat"
	dbUser := "root"
	dbPassword := "root"
	
	// Construct the database connection URL
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open a database connection
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	// Test the database connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}


