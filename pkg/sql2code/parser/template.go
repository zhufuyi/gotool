package parser

import (
	"sync"
	"text/template"
)

var (
	structTmpl    *template.Template
	structTmplRaw = `
{{- if .Comment -}}
// {{.TableName}} {{.Comment}}
{{end -}}
type {{.TableName}} struct {
{{- range .Fields}}
	{{.Name}} {{.GoType}} {{if .Tag}}` + "`{{.Tag}}`" + `{{end}}{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
{{if .NameFunc}}
// TableName table name
func (m *{{.TableName}}) TableName() string {
	return "{{.RawTableName}}"
}
{{end}}
`

	fileTmpl    *template.Template
	fileTmplRaw = `package {{.Package}}
{{if .ImportPath}}
import (
	{{- range .ImportPath}}
	"{{.}}"
	{{- end}}
)
{{- end}}
{{range .StructCode}}
{{.}}
{{end}}
`

	updateFieldTmpl    *template.Template
	updateFieldTmplRaw = `
{{- range .Fields}}
	if table.{{.Name}} {{.ConditionZero}} {
		update["{{.ColName}}"] = table.{{.Name}}
	}
{{- end}}
`
	handlerPostStructTmpl    *template.Template
	handlerPostStructTmplRaw = `
// Create{{.TableName}}Request request form
type Create{{.TableName}}Request struct {
// todo fill in the binding rules https://github.com/go-playground/validator
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`
	handlerPutStructTmpl    *template.Template
	handlerPutStructTmplRaw = `
// Update{{.TableName}}ByIDRequest update form
type Update{{.TableName}}ByIDRequest struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerGetStructTmpl    *template.Template
	handlerGetStructTmplRaw = `
// Get{{.TableName}}ByIDRespond respond data
type Get{{.TableName}}ByIDRespond struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	modelJSONTmpl    *template.Template
	modelJSONTmplRaw = `{
{{- range .Fields}}
	"{{.ColName}}" {{.GoZero}}
{{- end}}
}
`

	tmplParseOnce sync.Once
)

func initTemplate() {
	tmplParseOnce.Do(func() {
		var err error
		structTmpl, err = template.New("goStruct").Parse(structTmplRaw)
		if err != nil {
			panic(err)
		}
		fileTmpl, err = template.New("goFile").Parse(fileTmplRaw)
		if err != nil {
			panic(err)
		}
		updateFieldTmpl, err = template.New("goUpdateField").Parse(updateFieldTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerPostStructTmpl, err = template.New("goPostStruct").Parse(handlerPostStructTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerPutStructTmpl, err = template.New("goPutStruct").Parse(handlerPutStructTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerGetStructTmpl, err = template.New("goGetStruct").Parse(handlerGetStructTmplRaw)
		if err != nil {
			panic(err)
		}
		modelJSONTmpl, err = template.New("modelJSON").Parse(modelJSONTmplRaw)
		if err != nil {
			panic(err)
		}
	})
}
