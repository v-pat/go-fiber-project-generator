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

Generate new project by sending post request to ```http://localhost:3000/generate?database=mysql```.
Request's body should be in following format:
```
{
    "appName":"go-fiber-app",
    "tables":[
        {
            "name":"table_1",
            "columns":{
                "id":1,
                "name":"vpat",
                "role":"SDE"
            },
            "endpoint":"table_1"
        },
    ]
}
```


