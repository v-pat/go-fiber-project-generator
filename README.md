# go-fiber-project-generator

go-fiber-project-generator is a template-based code generator for creating simple Go Fiber backend applications configured for MySQL.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)

## Features

- Quickly generate a Go Fiber backend application with MySQL integration.
- Customizable template structure to fit your project needs.
- Generates boilerplate code for models, routes, middleware, and database configurations.


## Installation
- Go (at least Go 1.16)
- MySQL 

### Installation Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/go-fiber-project-generator.git
    ```

2. Open terminal and run generate.go:
   
   ```bash
   go run generate.go
   ```
  

## Usage

- Generate new project by sending post request to ```http://localhost:3000/generate?database=mysql```. Request's body should be in following format:
  ```
  {
    "appName":"<APP_NAME>",
    "tables":[
        {
            "name":"<TABLE_NAME>",
            "columns":{
                "FIRST_COLUMN_NAME":<VALUE>,
                  ...
                "NTH_COLUMN_NAME":<VALUE>,
            },
            "endpoint":"<ENDPOINT_FOR_THIS_TABLE>"
        },
    ]
  }
  ```
- In columns object ```VALUE``` denotes a data expected to be in the row of the column or what is type of that column.
  
- For example if a column would have int type values in it, then columns object would be like :
  ```
  "columns":{
      "id":1,
  },
  ```
  
- Similarly if columns would have string or boolean type values in it, then columns object would be like:
  ```
  "columns":{
      "id":1,
      "name":"vpat",
      "is_developer":true,
  },
  ```
  
- Please note that the values given in front of column names, will be used to only determine type of column and will not be added in database as first row of table. 



