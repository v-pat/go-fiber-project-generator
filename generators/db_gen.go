package generators

import (
	"fmt"
	"os"
	"text/template"
	tmpl "vpat_codegen/templates"
)

// DatabaseConnectionParams represents the parameters required for generating database connection code.
type DatabaseConnectionParams struct {
	DatabaseDriver     string
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	DBURLFormat        string
	DatabaseDriverName string
}

// GenerateDatabaseConnectionCode generates code for connecting to a database (PostgreSQL or MySQL) and writes it to a file.
func generateDatabaseConnectionCode(params DatabaseConnectionParams, fileName string) error {
	// Create a new template
	tmpl, err := template.New("databaseConnection").Parse(tmpl.DatabaseConnectionTemplate)
	if err != nil {
		return err
	}

	// Create or open the file
	file, err := os.Create("generated/databases/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute the template and write the generated code to the file
	if err := tmpl.Execute(file, params); err != nil {
		return err
	}

	return nil
}

func CreateDatabase(database string, appName string) error {
	params := DatabaseConnectionParams{
		DatabaseDriver:     "github.com/lib/pq",
		DBHost:             "localhost",
		DBPort:             "5432",
		DBName:             appName,
		DBUser:             "myuser",
		DBPassword:         "mypassword",
		DBURLFormat:        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DatabaseDriverName: "postgres",
	}

	if database == "postgres" {
		// Example usage: Generate code for PostgreSQL database connection and write to a file
		err := generateDatabaseConnectionCode(params, "postgres_connection.go")
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	} else if database == "mysql" {
		// Example usage: Generate code for MySQL database connection and write to a file
		params.DatabaseDriver = "github.com/go-sql-driver/mysql"
		params.DBPort = "3306"
		params.DatabaseDriverName = "mysql"
		params.DBUser = "root"
		params.DBPassword = "root"
		params.DBURLFormat = "%s:%s@tcp(%s:%s)/%s"
		params.DBName = appName

		err := generateDatabaseConnectionCode(params, "mysql_connection.go")
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

	}

	fmt.Println("Database connection code generated and written to files successfully.")

	return nil
}
