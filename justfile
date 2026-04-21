set shell := ["/bin/bash", "-cu"]
set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]
set dotenv-load := true

root_dir := justfile_directory()
default_cmd := "qbot"

# Recipe setup will enable githooks.
setup:
    git config core.hooksPath .githooks

# Recipe run will run the given cmd which defaults to qbot.
run cmd=default_cmd:
    go run {{ root_dir }}/cmd/{{ cmd }}/main.go

# Recipe repomix, for macos, will run repomix and copy the resulting file (not the contents but the file itself) to the clipboard.
[macos]
repomix:
    repomix {{ root_dir }}
    # Copy the resulting .repomix file to clipboard for ease of sharing with an LLM. This is different from the --copy argument since it copies the file itself, not just the contents.
    osascript -e 'tell application "Finder" to set the clipboard to (POSIX file "{{ root_dir }}/.repomix")'

# Recipe repomix, for windows, will run repomix and copy the resulting file (not the contents but the file itself) to the clipboard.
[windows]
repomix:
    repomix {{ root_dir }}
    : Copy the resulting .repomix file to clipboard for ease of sharing with an LLM. This is different from the --copy argument since it copies the file itself, not just the contents.
    powershell Set-Clipboard -Path '{{ root_dir }}/.repomix'

# Recipe compile will compile the entire codebase, which is useful for detecting any compilation errors across the project.
compile:
    go build {{ root_dir }}/...

# Recipe test will run all unit tests across the codebase.
test:
    go test {{ root_dir }}/...

# Recipe ssh will ssh into the qbot app on fly.io.
ssh:
    echo For running SQL queries do:
    echo sqlite3 /var/lib/litefs/qbot.db
    echo For editing config do:
    echo "apt-get update && apt-get install vim -y && vim /var/lib/litefs/qbot/config.json"
    fly ssh console -a qbot

# Recipe restart will restart the qbot app on fly.io.
restart:
    fly machines restart -a qbot

# Recipe push-config will upload the local config.json to the deployed bot on fly.io.
push-config:
    fly -a qbot ssh console --command "rm /var/lib/litefs/qbot/config.json"
    fly -a qbot sftp put "./config.json" "/var/lib/litefs/qbot/config.json"

# Recipe backup will run backup.bat which copies the sqlite db file from fly to the local machine.
backup:
    backup.bat
