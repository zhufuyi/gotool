# goctl

go language development tools.

```
$ goctl -h
go language development tools

Usage:
  goctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gen         Generate gin api code
  help        Help about any command
  resources   List of supported resources

Flags:
  -h, --help      help for goctl
  -v, --version   version for goctl

Use "goctl [command] --help" for more information about a command.
```

<br>

## usage

(1) generate web service code

> goctl gen web -p yourProjectName -a yourApiName

<br>

(2) generate api code

> goctl gen api -p yourProjectName -a yourApiName

copy the generated code folder to your project folder and fill in the business logic code to complete an api interface add, delete and query function
