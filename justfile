set shell := ["/bin/bash", "-cu"]
set windows-shell := ["cmd", "/c"]
set dotenv-load := true

root_dir := justfile_directory()
default_cmd := "qbot"

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
    # Copy the resulting .repomix file to clipboard for ease of sharing with an LLM. This is different from the --copy argument since it copies the file itself, not just the contents.
    powershell Set-Clipboard -Path '{{ root_dir }}/.repomix'
