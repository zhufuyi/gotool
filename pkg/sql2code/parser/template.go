package parser

import (
	"sync"
	"text/template"
)

var (
	modelStructTmpl    *template.Template
	modelStructTmplRaw = `
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

	modelTmpl    *template.Template
	modelTmplRaw = `package {{.Package}}
{{if .ImportPath}}
import (
	{{- range .ImportPath}}
	"{{.}}"
	{{- end}}
)
{{- end}}
{{range .StructCode}}
{{.}}
{{end}}`

	updateFieldTmpl    *template.Template
	updateFieldTmplRaw = `
{{- range .Fields}}
	if table.{{.Name}} {{.ConditionZero}} {
		update["{{.ColName}}"] = table.{{.Name}}
	}
{{- end}}`

	handlerCreateStructTmpl    *template.Template
	handlerCreateStructTmplRaw = `
// Create{{.TableName}}Request create params
type Create{{.TableName}}Request struct {
// todo fill in the binding rules https://github.com/go-playground/validator
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerUpdateStructTmpl    *template.Template
	handlerUpdateStructTmplRaw = `
// Update{{.TableName}}ByIDRequest update params
type Update{{.TableName}}ByIDRequest struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerDetailStructTmpl    *template.Template
	handlerDetailStructTmplRaw = `
// Get{{.TableName}}ByIDRespond respond detail
type Get{{.TableName}}ByIDRespond struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}`

	modelJSONTmpl    *template.Template
	modelJSONTmplRaw = `{
{{- range .Fields}}
	"{{.ColName}}" {{.GoZero}}
{{- end}}
}
`

	protoFileTmpl    *template.Template
	protoFileTmplRaw = `syntax = "proto3";

package api.{{.TName}}.v1;

import "api/{{.TName}}/v1/query.proto";
// import "gogo/protobuf/gogoproto/gogo.proto";
// import "validate/validate.proto";

option go_package = "./pb;pb"; // 生成的文件在proto文件同级别pb目录下

service {{.TName}}Service {
  rpc Create(Create{{.TableName}}Request) returns (Create{{.TableName}}Reply) {}
  rpc DeleteByID(Delete{{.TableName}}ByIDRequest) returns (Delete{{.TableName}}ByIDReply) {}
  rpc UpdateByID(Update{{.TableName}}ByIDRequest) returns (Update{{.TableName}}ByIDReply) {}
  rpc GetByID(Get{{.TableName}}ByIDRequest) returns (Get{{.TableName}}ByIDReply) {}
  rpc List(List{{.TableName}}Request) returns (List{{.TableName}}Reply) {}
}

// protoMessageCreateCode

message Create{{.TableName}}Reply {
  uint64   id =1;
}

message Delete{{.TableName}}ByIDRequest {
  uint64   id =1;
}

message Delete{{.TableName}}ByIDReply {

}

// protoMessageUpdateCode

message Update{{.TableName}}ByIDReply {

}

// protoMessageDetailCode

message Get{{.TableName}}ByIDRequest {
  uint64   id =1;
}

message Get{{.TableName}}ByIDReply {
  {{.TableName}} {{.TName}} = 1;
}

message List{{.TableName}}Request {
  Params params = 1;
}

message List{{.TableName}}Reply {
  int64 total =1;
  repeated {{.TableName}} {{.TName}}s = 2;
}
`

	protoMessageCreateTmpl    *template.Template
	protoMessageCreateTmplRaw = `message Create{{.TableName}}Request {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.ColName}} = {{$v.AddOne $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	protoMessageUpdateTmpl    *template.Template
	protoMessageUpdateTmplRaw = `message Update{{.TableName}}ByIDRequest {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.ColName}} = {{$v.AddOne $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	protoMessageDetailTmpl    *template.Template
	protoMessageDetailTmplRaw = `message {{.TableName}} {
{{- range $i, $v := .Fields}}
	{{$v.GoType}} {{$v.ColName}} = {{$v.AddOne $i}}; {{if $v.Comment}} // {{$v.Comment}}{{end}}
{{- end}}
}`

	tmplParseOnce sync.Once
)

func initTemplate() {
	tmplParseOnce.Do(func() {
		var err error
		modelStructTmpl, err = template.New("goStruct").Parse(modelStructTmplRaw)
		if err != nil {
			panic(err)
		}
		modelTmpl, err = template.New("goFile").Parse(modelTmplRaw)
		if err != nil {
			panic(err)
		}
		updateFieldTmpl, err = template.New("goUpdateField").Parse(updateFieldTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerCreateStructTmpl, err = template.New("goPostStruct").Parse(handlerCreateStructTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerUpdateStructTmpl, err = template.New("goPutStruct").Parse(handlerUpdateStructTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerDetailStructTmpl, err = template.New("goGetStruct").Parse(handlerDetailStructTmplRaw)
		if err != nil {
			panic(err)
		}
		modelJSONTmpl, err = template.New("modelJSON").Parse(modelJSONTmplRaw)
		if err != nil {
			panic(err)
		}
		protoFileTmpl, err = template.New("protoFile").Parse(protoFileTmplRaw)
		if err != nil {
			panic(err)
		}
		protoMessageCreateTmpl, err = template.New("protoMessageCreate").Parse(protoMessageCreateTmplRaw)
		if err != nil {
			panic(err)
		}
		protoMessageUpdateTmpl, err = template.New("protoMessageUpdate").Parse(protoMessageUpdateTmplRaw)
		if err != nil {
			panic(err)
		}
		protoMessageDetailTmpl, err = template.New("protoMessageDetail").Parse(protoMessageDetailTmplRaw)
		if err != nil {
			panic(err)
		}
	})
}
