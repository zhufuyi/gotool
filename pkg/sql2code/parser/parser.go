package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"

	"github.com/blastrain/vitess-sqlparser/tidbparser/ast"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/mysql"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/types"
	"github.com/blastrain/vitess-sqlparser/tidbparser/parser"
	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
)

const (
	// TableName table name
	TableName = "__table_name__"
	// CodeTypeModel model code
	CodeTypeModel = "model"
	// CodeTypeJSON json code
	CodeTypeJSON = "json"
	// CodeTypeDAO update fields code
	CodeTypeDAO = "dao"
	// CodeTypeHandler handler request and respond code
	CodeTypeHandler = "handler"
)

// Codes 生成的代码
type Codes struct {
	Model         []string // model code
	UpdateFields  []string // update fields code
	ModelJSON     []string // model json code
	HandlerStruct []string // handler request and respond code
}

// modelCodes 生成的代码
type modelCodes struct {
	Package    string
	ImportPath []string
	StructCode []string
}

// ParseSql 根据sql生成不用用途代码
func ParseSql(sql string, options ...Option) (map[string]string, error) {
	initTemplate()
	opt := parseOption(options)

	stmts, err := parser.New().Parse(sql, opt.Charset, opt.Collation)
	if err != nil {
		return nil, err
	}
	modelStructCodes := make([]string, 0, len(stmts))
	updateFieldsCodes := make([]string, 0, len(stmts))
	handlerStructCodes := make([]string, 0, len(stmts))
	modelJSONCodes := make([]string, 0, len(stmts))
	importPath := make(map[string]struct{})
	tableNames := make([]string, 0, len(stmts))
	for _, stmt := range stmts {
		if ct, ok := stmt.(*ast.CreateTableStmt); ok {
			code, err := makeCode(ct, opt)
			if err != nil {
				return nil, err
			}
			modelStructCodes = append(modelStructCodes, code.modelStruct)
			updateFieldsCodes = append(updateFieldsCodes, code.updateFields)
			handlerStructCodes = append(handlerStructCodes, code.handlerStruct)
			modelJSONCodes = append(modelJSONCodes, code.modelJSON)
			tableNames = append(tableNames, toCamel(ct.Table.Name.String()))
			for _, s := range code.importPaths {
				importPath[s] = struct{}{}
			}
		}
	}

	importPathArr := make([]string, 0, len(importPath))
	for s := range importPath {
		importPathArr = append(importPathArr, s)
	}
	sort.Strings(importPathArr)

	mc := modelCodes{
		Package:    opt.Package,
		ImportPath: importPathArr,
		StructCode: modelStructCodes,
	}
	modelCode, err := getModelCode(mc)
	if err != nil {
		return nil, err
	}

	var codesMap = map[string]string{
		CodeTypeModel:   modelCode,
		CodeTypeJSON:    strings.Join(modelJSONCodes, "\n\n"),
		CodeTypeDAO:     strings.Join(updateFieldsCodes, "\n\n"),
		CodeTypeHandler: strings.Join(handlerStructCodes, "\n\n"),
		TableName:       strings.Join(tableNames, ", "),
	}

	return codesMap, nil
}

// ConfigureAcronym config acronym
func ConfigureAcronym(words []string) {
	acronym = make(map[string]struct{}, len(words))
	for _, w := range words {
		acronym[strings.ToUpper(w)] = struct{}{}
	}
}

type tmplData struct {
	TableName    string
	NameFunc     bool
	RawTableName string
	Fields       []tmplField
	Comment      string
}

type tmplField struct {
	Name    string
	ColName string
	GoType  string
	Tag     string
	Comment string
}

// ConditionZero type of condition 0
func (t tmplField) ConditionZero() string {
	switch t.GoType {
	case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint", "float64", "float32":
		return `!= 0`
	case "string":
		return `!= ""`
	case "time.Time":
		return `.IsZero() == false`
	}

	return `!= ` + t.GoType
}

// GoZero type of 0 逗号分隔
func (t tmplField) GoZero() string {
	switch t.GoType {
	case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint", "float64", "float32":
		return `= 0`
	case "string":
		return `= "string"`
	case "time.Time":
		return `= "0000-01-00T00:00:00.000+08:00"`
	}

	return `= ` + t.GoType
}

const (
	__mysqlModel__ = "__mysqlModel__"
	__type__       = "__type__"
)

var replaceFields = map[string]string{
	__mysqlModel__: "mysql.Model",
	__type__:       "",
}

const (
	columnID         = "id"
	columnCreatedAt  = "created_at"
	columnUpdatedAt  = "updated_at"
	columnDeletedAt  = "deleted_at"
	columnMysqlModel = __mysqlModel__
)

var ignoreColumns = map[string]struct{}{
	columnID:         struct{}{},
	columnCreatedAt:  struct{}{},
	columnUpdatedAt:  struct{}{},
	columnDeletedAt:  struct{}{},
	columnMysqlModel: struct{}{},
}

func isIgnoreFields(colName string, falseColumn ...string) bool {
	for _, v := range falseColumn {
		if colName == v {
			return false
		}
	}

	_, ok := ignoreColumns[colName]
	return ok
}

type codeText struct {
	importPaths   []string
	modelStruct   string
	modelJSON     string
	updateFields  string
	handlerStruct string
}

func makeCode(stmt *ast.CreateTableStmt, opt options) (*codeText, error) {
	importPath := make([]string, 0, 1)
	data := tmplData{
		TableName:    stmt.Table.Name.String(),
		RawTableName: stmt.Table.Name.String(),
		Fields:       make([]tmplField, 0, 1),
	}
	tablePrefix := opt.TablePrefix
	if tablePrefix != "" && strings.HasPrefix(data.TableName, tablePrefix) {
		data.NameFunc = true
		data.TableName = data.TableName[len(tablePrefix):]
	}
	if opt.ForceTableName || data.RawTableName != inflection.Plural(data.RawTableName) {
		data.NameFunc = true
	}

	data.TableName = toCamel(data.TableName)

	// find table comment
	for _, opt := range stmt.Options {
		if opt.Tp == ast.TableOptionComment {
			data.Comment = opt.StrValue
			break
		}
	}

	isPrimaryKey := make(map[string]bool)
	for _, con := range stmt.Constraints {
		if con.Tp == ast.ConstraintPrimaryKey {
			isPrimaryKey[con.Keys[0].Column.String()] = true
		}
	}

	columnPrefix := opt.ColumnPrefix
	for _, col := range stmt.Cols {
		colName := col.Name.Name.String()
		goFieldName := colName
		if columnPrefix != "" && strings.HasPrefix(goFieldName, columnPrefix) {
			goFieldName = goFieldName[len(columnPrefix):]
		}

		field := tmplField{
			Name:    toCamel(goFieldName),
			ColName: colName,
		}

		tags := make([]string, 0, 4)
		// make GORM's tag
		gormTag := strings.Builder{}
		gormTag.WriteString("column:")
		gormTag.WriteString(colName)
		if opt.GormType {
			gormTag.WriteString(";type:")
			gormTag.WriteString(col.Tp.InfoSchemaStr())
		}
		if isPrimaryKey[colName] {
			gormTag.WriteString(";primary_key")
		}
		isNotNull := false
		canNull := false
		for _, o := range col.Options {
			switch o.Tp {
			case ast.ColumnOptionPrimaryKey:
				if !isPrimaryKey[colName] {
					gormTag.WriteString(";primary_key")
					isPrimaryKey[colName] = true
				}
			case ast.ColumnOptionNotNull:
				isNotNull = true
			case ast.ColumnOptionAutoIncrement:
				gormTag.WriteString(";AUTO_INCREMENT")
			case ast.ColumnOptionDefaultValue:
				if value := getDefaultValue(o.Expr); value != "" {
					gormTag.WriteString(";default:")
					gormTag.WriteString(value)
				}
			case ast.ColumnOptionUniqKey:
				gormTag.WriteString(";unique")
			case ast.ColumnOptionNull:
				//gormTag.WriteString(";NULL")
				canNull = true
			case ast.ColumnOptionOnUpdate: // For Timestamp and Datetime only.
			case ast.ColumnOptionFulltext:
			case ast.ColumnOptionComment:
				field.Comment = o.Expr.GetDatum().GetString()
			default:
				//return "", nil, errors.Errorf(" unsupport option %d\n", o.Tp)
			}
		}
		if !isPrimaryKey[colName] && isNotNull {
			gormTag.WriteString(";NOT NULL")
		}
		tags = append(tags, "gorm", gormTag.String())

		if opt.JsonTag {
			if opt.JsonNamedType != 0 {
				colName = xstrings.FirstRuneToLower(xstrings.ToCamelCase(colName)) // 使用驼峰类型json名称
			}
			tags = append(tags, "json", colName)
		}

		field.Tag = makeTagStr(tags)

		// get type in golang
		nullStyle := opt.NullStyle
		if !canNull {
			nullStyle = NullDisable
		}
		goType, pkg := mysqlToGoType(col.Tp, nullStyle)
		if pkg != "" {
			importPath = append(importPath, pkg)
		}
		field.GoType = goType

		data.Fields = append(data.Fields, field)
	}

	updateFieldsCode, err := getUpdateFieldsCode(data, opt.IsEmbed)
	if err != nil {
		return nil, err
	}

	handlerStructCode, err := getHandlerStructCodes(data)
	if err != nil {
		return nil, err
	}

	modelStructCode, importPaths, err := getModelStructCode(data, importPath, opt.IsEmbed)
	if err != nil {
		return nil, err
	}

	modelJSONCode, err := getModelJSONCode(data)
	if err != nil {
		return nil, err
	}

	//return modelStructCode, importPaths, nil
	return &codeText{
		importPaths:   importPaths,
		modelStruct:   modelStructCode,
		modelJSON:     modelJSONCode,
		updateFields:  updateFieldsCode,
		handlerStruct: handlerStructCode,
	}, nil
}

func getModelStructCode(data tmplData, importPaths []string, isEmbed bool) (string, []string, error) {
	// 过滤忽略字段字段
	var newFields = []tmplField{}
	var newImportPaths = []string{}
	if isEmbed {
		// 嵌入字段
		newFields = append(newFields, tmplField{
			Name:    __mysqlModel__,
			ColName: __mysqlModel__,
			GoType:  __type__,
			Tag:     `gorm:"embedded"`,
			Comment: "embed id and time\n",
		})
		for _, field := range data.Fields {
			if isIgnoreFields(field.ColName) {
				continue
			}
			newFields = append(newFields, field)
		}
		data.Fields = newFields

		// 过滤time包
		isOk := false
		for _, field := range newFields {
			if !strings.Contains(field.GoType, "time.Time") {
				isOk = true
			}
		}
		for _, path := range importPaths {
			if isOk && path == "time" {
				continue
			}
			newImportPaths = append(newImportPaths, path)
		}
		newImportPaths = append(newImportPaths, "github.com/zhufuyi/pkg/mysql")
	} else {
		newImportPaths = importPaths
	}

	builder := strings.Builder{}
	err := structTmpl.Execute(&builder, data)
	if err != nil {
		return "", nil, fmt.Errorf("structTmpl.Execute error: %v", err)
	}
	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", nil, fmt.Errorf("structTmpl format.Source error: %v", err)
	}
	structCode := string(code)
	// 还原真实的嵌入字段
	if isEmbed {
		structCode = strings.ReplaceAll(structCode, __mysqlModel__, replaceFields[__mysqlModel__])
		structCode = strings.ReplaceAll(structCode, __type__, replaceFields[__type__])
	}

	return structCode, newImportPaths, nil
}

func getModelCode(data modelCodes) (string, error) {
	builder := strings.Builder{}
	err := fileTmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}

	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", fmt.Errorf("format.Source error: %v", err)
	}

	return string(code), nil
}

func getUpdateFieldsCode(data tmplData, isEmbed bool) (string, error) {
	// 过滤字段
	var newFields = []tmplField{}
	for _, field := range data.Fields {
		falseColumns := []string{}
		if !isEmbed {
			falseColumns = append(falseColumns, columnCreatedAt, columnUpdatedAt, columnDeletedAt)
		}
		if isIgnoreFields(field.ColName, falseColumns...) {
			continue
		}
		newFields = append(newFields, field)
	}
	data.Fields = newFields

	builder := strings.Builder{}
	err := updateFieldTmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}

	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", err
	}

	return string(code), nil
}

func getHandlerStructCodes(data tmplData) (string, error) {
	postStructCode, err := getHandlerStructCode(data, handlerPostStructTmpl)
	if err != nil {
		return "", fmt.Errorf("handlerPostStructTmpl error: %v", err)
	}

	putStructCode, err := getHandlerStructCode(data, handlerPutStructTmpl, columnID)
	if err != nil {
		return "", fmt.Errorf("handlerPutStructTmpl error: %v", err)
	}

	getStructCode, err := getHandlerStructCode(data, handlerGetStructTmpl, columnID, columnCreatedAt, columnUpdatedAt)
	if err != nil {
		return "", fmt.Errorf("handlerGetStructTmpl error: %v", err)
	}

	return postStructCode + putStructCode + getStructCode, nil
}

func getHandlerStructCode(data tmplData, tmpl *template.Template, reservedColumns ...string) (string, error) {
	// post过滤字段
	var newFields = []tmplField{}
	for _, field := range data.Fields {
		if isIgnoreFields(field.ColName, reservedColumns...) {
			continue
		}
		newFields = append(newFields, field)
	}
	data.Fields = newFields

	builder := strings.Builder{}
	err := tmpl.Execute(&builder, data)
	if err != nil {
		return "", fmt.Errorf("tmpl.Execute error: %v", err)
	}
	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", fmt.Errorf("format.Source error: %v", err)
	}

	return string(code), nil
}

func getModelJSONCode(data tmplData) (string, error) {
	builder := strings.Builder{}
	err := modelJSONTmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}

	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", fmt.Errorf("format.Source error: %v", err)
	}

	modelJSONCode := strings.ReplaceAll(string(code), " =", ":")
	modelJSONCode = addCommaToJSON(modelJSONCode)

	return modelJSONCode, nil
}

func addCommaToJSON(modelJSONCode string) string {
	r := strings.NewReader(modelJSONCode)
	buf := bufio.NewReader(r)

	lines := []string{}
	count := 0
	for {
		line, err := buf.ReadString(byte('\n'))
		if err != nil {
			break
		}
		lines = append(lines, line)
		if len(line) > 5 {
			count++
		}
	}

	out := ""
	for _, line := range lines {
		if len(line) < 5 && (strings.Contains(line, "{") || strings.Contains(line, "}")) {
			out += line
			continue
		}
		count--
		if count == 0 {
			out += line
			continue
		}
		index := bytes.IndexByte([]byte(line), '\n')
		out += line[:index] + "," + line[index:]
	}
	return out
}

func mysqlToGoType(colTp *types.FieldType, style NullStyle) (name string, path string) {
	if style == NullInSql {
		path = "database/sql"
		switch colTp.Tp {
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong:
			name = "sql.NullInt32"
		case mysql.TypeLonglong:
			name = "sql.NullInt64"
		case mysql.TypeFloat, mysql.TypeDouble:
			name = "sql.NullFloat64"
		case mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString,
			mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
			name = "sql.NullString"
		case mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate:
			name = "sql.NullTime"
		case mysql.TypeDecimal, mysql.TypeNewDecimal:
			name = "sql.NullString"
		case mysql.TypeJSON:
			name = "sql.NullString"
		default:
			return "UnSupport", ""
		}
	} else {
		switch colTp.Tp {
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint"
			} else {
				name = "int"
			}
		case mysql.TypeLonglong:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint64"
			} else {
				name = "int64"
			}
		case mysql.TypeFloat, mysql.TypeDouble:
			name = "float64"
		case mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString,
			mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
			name = "string"
		case mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate:
			path = "time"
			name = "time.Time"
		case mysql.TypeDecimal, mysql.TypeNewDecimal:
			name = "string"
		case mysql.TypeJSON:
			name = "string"
		default:
			return "UnSupport", ""
		}
		if style == NullInPointer {
			name = "*" + name
		}
	}
	return
}

func makeTagStr(tags []string) string {
	builder := strings.Builder{}
	for i := 0; i < len(tags)/2; i++ {
		builder.WriteString(tags[i*2])
		builder.WriteString(`:"`)
		builder.WriteString(tags[i*2+1])
		builder.WriteString(`" `)
	}
	if builder.Len() > 0 {
		return builder.String()[:builder.Len()-1]
	}
	return builder.String()
}

func getDefaultValue(expr ast.ExprNode) (value string) {
	if expr.GetDatum().Kind() != types.KindNull {
		value = fmt.Sprintf("%v", expr.GetDatum().GetValue())
	} else if expr.GetFlag() != ast.FlagConstant {
		if expr.GetFlag() == ast.FlagHasFunc {
			if funcExpr, ok := expr.(*ast.FuncCallExpr); ok {
				value = funcExpr.FnName.O
			}
		}
	}
	return
}

var acronym = map[string]struct{}{
	"ID":  {},
	"IP":  {},
	"RPC": {},
}

func toCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	s += "."

	n := strings.Builder{}
	n.Grow(len(s))
	temp := strings.Builder{}
	temp.Grow(len(s))
	wordFirst := true
	for _, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if wordFirst && vIsLow {
			v -= 'a' - 'A'
		}

		if vIsCap || vIsLow {
			temp.WriteByte(v)
			wordFirst = false
		} else {
			isNum := v >= '0' && v <= '9'
			wordFirst = isNum || v == '_' || v == ' ' || v == '-' || v == '.'
			if temp.Len() > 0 && wordFirst {
				word := temp.String()
				upper := strings.ToUpper(word)
				if _, ok := acronym[upper]; ok {
					n.WriteString(upper)
				} else {
					n.WriteString(word)
				}
				temp.Reset()
			}
			if isNum {
				n.WriteByte(v)
			}
		}
	}
	return n.String()
}
