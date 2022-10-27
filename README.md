## gotool

go language development tools.

## Install

> go install github.com/zhufuyi/gotool@latest

<br>

the installation path is in `$GOPATH/bin`, see the command help.

```
go language development tools

Usage:
  gotool [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  covert      resource type conversion
  help        Help about any command
  replace     Replace fields in path files

Flags:
  -h, --help      help for gotool
  -v, --version   version for gotool

Use "gotool [command] --help" for more information about a command.
```

<br>

## Usage

### Replace command

```bash
# replace one field
gotool replace -p /tmp -o oldField -n newField

# replace multiple fields
gotool replace -p /tmp -o oldField1 -n newField1 -o oldField2 -n newField2
```

<br>

### Covert command

#### sql to gorm

```bash
# covert sql to gorm from file
gotool covert sql --file=test.sql

# covert sql to gorm from db
gotool covert sql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

# covert sql to gorm, set package name and json tag
gotool covert sql --file=test.sql --pkg-name=user --json-tag
gotool covert sql --file=test.sql --pkg-name=user --json-tag --json-named-type=1
```

<br>

#### json to go struct

```bash
  # covert json to struct from data
  gotool covert json --data="json text"

  # covert json to struct from file
  gotool covert json --file=test.json

  # covert json to struct, set tag value
  gotool covert json --file=test.json --tags=gorm

  # covert yaml to struct, save to specified directory, file name is config.go
  gotool covert json --file=test.json --out=/tmp
```

<br>

#### yaml to go struct

```bash
  # covert yaml to struct from data
  gotool covert yaml --data="yaml text"

  # covert yaml to struct from file
  gotool covert yaml --file=test.yaml

  # covert yaml to struct, set tag value, save to specified directory, file name is config.go
  gotool covert yaml --file=test.yaml --tags=json --out=/tmp
```
