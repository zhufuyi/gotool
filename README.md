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
  gen         Generate web service code
  help        Help about any command
  replace     Replace fields in path files

Flags:
  -h, --help      help for goctl
  -v, --version   version for goctl

Use "goctl [command] --help" for more information about a command.
```

<br>

## Usage

### Generate command

#### generate web service code

> goctl gen web -p yourProjectName -a yourApiName

<br>

#### generate api code

> goctl gen api -p yourProjectName -a yourApiName

copy the generated code folder to your project folder and fill in the business logic code to complete an api interface add, delete and query function

<br>

#### generate user code

> goctl gen user -p yourProjectName

The user code includes registration, login and logout api interfaces, including authentication and ip rate limits functions, and can be modified to suit the actual business.

<br>

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

    # covert json to struct, set tag value and subStruct flag
    goctl covert json --file=test.sql --tags=gorm --sub-struct=false
```

<br>

#### yaml to struct

```bash
    # covert yaml to struct from data
    goctl covert yaml --data="yaml text"

    # covert yaml to struct from file
    goctl covert yaml --file=test.json

    # covert yaml to struct, set tag value and subStruct flag
    goctl covert yaml --file=test.sql --tags=gorm --sub-struct=false
```
