package templates

// Define a model template for generating the struct code in the model file.
const ModelTemplate = `package model

// {{.StructName}} represents the {{.StructName}} struct.
type {{.StructName}} struct {
{{range .Fields}}
	{{.Name}} {{.Type}} ` + "`json:\"{{.Name}}\"`" +
	`{{end}}
}`
