# goctl

go language development tools.

## Install

> go install github.com/zhufuyi/goctl@latest

<br>

the installation path is in `$GOPATH/bin`, see the command help.

```
go language development tools

Usage:
  goctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  covert      resource type conversion
  help        Help about any command
  replace     Replace fields in path files

Flags:
  -h, --help      help for goctl
  -v, --version   version for goctl

Use "goctl [command] --help" for more information about a command.
```

<br>

## Usage

### Replace command

```bash
# replace one field
goctl replace -p /tmp -o oldField -n newField

# replace multiple fields
goctl replace -p /tmp -o oldField1 -n newField1 -o oldField2 -n newField2
```

<br>

### Covert command

#### sql to gorm

```bash
# covert sql to gorm from file
goctl covert sql --file=test.sql

# covert sql to gorm from db
goctl covert sql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

# covert sql to gorm, set package name and json tag
goctl covert sql --file=test.sql --pkg-name=user --json-tag
goctl covert sql --file=test.sql --pkg-name=user --json-tag --json-named-type=1
```

<br>

#### json to struct

```bash
  # covert json to struct from data
  goctl covert json --data="json text"

  # covert json to struct from file
  goctl covert json --file=test.json

  # covert json to struct, set tag value
  goctl covert json --file=test.json --tags=gorm

  # covert yaml to struct, save to specified directory, file name is config.go
  goctl covert json --file=test.json --out=/tmp
```

<br>

#### yaml to struct

```bash
  # covert yaml to struct from data
  goctl covert yaml --data="yaml text"

  # covert yaml to struct from file
  goctl covert yaml --file=test.yaml

  # covert yaml to struct, set tag value, save to specified directory, file name is config.go
  goctl covert yaml --file=test.yaml --tags=json --out=/tmp
```
