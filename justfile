set shell := ["cmd", "/c"]

set dotenv-load := true

root_dir := justfile_directory()

default_cmd := "qbot"
run cmd=default_cmd:
    go run {{root_dir}}/cmd/{{cmd}}/main.go

repomix:
    repomix
    powershell Set-Clipboard -Path '{{root_dir}}/.repomix'
