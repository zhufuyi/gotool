package sql2code

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/zhufuyi/goctl/pkg/sql2code/parser"
)

// Args 参数
type Args struct {
	Sql string // DDL sql

	InputFile string // 读取文件的DDL sql

	DBDsn   string // 从db获取表的DDL sq
	DBTable string

	Package        string // 生成字段的包名
	JsonTag        bool   // 是否包括json tag
	JsonNamedType  int    // json命名类型，0:默认，其他值表示驼峰
	GormType       bool   // gorm type
	IsEmbed        bool   // 是否嵌入gorm.Model
	CodeType       string // 指定生成代码用途，支持4中类型，分别是 model, json, dao, handler
	ForceTableName bool
	Charset        string
	Collation      string
	TablePrefix    string
	ColumnPrefix   string
	NoNullType     bool
	NullStyle      string
}

func getSql(args *Args) (string, error) {
	if args.Sql != "" {
		return args.Sql, nil
	}

	sql := ""
	if args.InputFile != "" {
		b, err := ioutil.ReadFile(args.InputFile)
		if err != nil {
			return sql, fmt.Errorf("read %s failed, %s\n", args.InputFile, err)
		}
		return string(b), nil

	} else if args.DBDsn != "" {
		if args.DBTable == "" {
			return sql, errors.New("miss mysql table")
		}
		sqlStr, err := parser.GetCreateTableFromDB(args.DBDsn, args.DBTable)
		if err != nil {
			return sql, fmt.Errorf("get create table error: %s", err)
		}
		return sqlStr, nil
	}

	return sql, errors.New("no SQL input(-sql|-f|-db-dsn)")
}

func getOptions(args *Args) []parser.Option {
	var opts []parser.Option

	if args.Charset != "" {
		opts = append(opts, parser.WithCharset(args.Charset))
	}
	if args.Collation != "" {
		opts = append(opts, parser.WithCollation(args.Collation))
	}
	if args.JsonTag {
		opts = append(opts, parser.WithJsonTag(args.JsonNamedType))
	}
	if args.TablePrefix != "" {
		opts = append(opts, parser.WithTablePrefix(args.TablePrefix))
	}
	if args.ColumnPrefix != "" {
		opts = append(opts, parser.WithColumnPrefix(args.ColumnPrefix))
	}
	if args.NoNullType {
		opts = append(opts, parser.WithNoNullType())
	}
	if args.IsEmbed {
		opts = append(opts, parser.WithEmbed())
	}

	if args.NullStyle != "" {
		switch args.NullStyle {
		case "sql":
			opts = append(opts, parser.WithNullStyle(parser.NullInSql))
		case "ptr":
			opts = append(opts, parser.WithNullStyle(parser.NullInPointer))
		default:
			fmt.Printf("invalid null style: %s\n", args.NullStyle)
			return nil
		}
	} else {
		opts = append(opts, parser.WithNullStyle(parser.NullDisable))
	}
	if args.Package != "" {
		opts = append(opts, parser.WithPackage(args.Package))
	}
	if args.GormType {
		opts = append(opts, parser.WithGormType())
	}
	if args.ForceTableName {
		opts = append(opts, parser.WithForceTableName())
	}

	return opts
}

// GetGormCode 根据sql生成gorm代码，sql可以从参数、文件、db三种方式获取，优先从高到低
func GetGormCode(args *Args) (string, error) {
	codes, err := GetCodes(args)
	if err != nil {
		return "", err
	}

	if args.CodeType == "" {
		args.CodeType = parser.CodeTypeModel // 默认为model code
	}
	out, ok := codes[args.CodeType]
	if !ok {
		return "", fmt.Errorf("unkown code type %s", args.CodeType)
	}

	return out, nil
}

// GetCodes 生成model, json, dao, handler不同用途代码
func GetCodes(args *Args) (map[string]string, error) {
	sql, err := getSql(args)
	if err != nil {
		return nil, err
	}

	opt := getOptions(args)

	return parser.ParseSql(sql, opt...)
}
