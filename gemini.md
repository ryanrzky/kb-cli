# Project Overview: kb-cli

`kb-cli` is a command-line interface tool written in Go that interacts with the Kanboard JSON-RPC API. It allows users to manage projects and tasks in Kanboard directly from their terminal.

## Tech Stack
- **Language:** Go
- **CLI Framework:** Cobra (`github.com/spf13/cobra`)
- **Configuration Management:** Viper (`github.com/spf13/viper`)
- **API Communication:** Kanboard JSON-RPC API over HTTP POST

## Project Structure
- `main.go`: Entry point of the application.
- `client/kanboard.go`: Contains the HTTP client (`KanboardClient`) that handles formatting and sending JSON-RPC requests to the Kanboard server.
- `cmd/root.go`: Initializes the root Cobra command and handles Viper configuration setup (flags and `.kb-cli.yaml`).
- `cmd/project.go`: Contains project-related commands (`project list`, `project create`, `project get`).
- `cmd/task.go`: Contains task-related commands (`task list`, `task create`, `task get`, `task update`, `task delete`). 

## Configuration
The CLI reads its configuration from `~/.kb-cli.yaml`.
Required fields in the config:
```yaml
url: "https://kanboard.example.com/jsonrpc.php"
user: "jsonrpc"
token: "YOUR_API_TOKEN"
```

## Key Features & Quirks
- **Task List Relational Mapping**: The `task list` command makes additional background requests (`getColumns`, `getProjectUsers`) so it can display human-readable column titles and assignee names instead of just raw IDs, as `getAllTasks` only returns IDs.
- **Nil Handling**: Kanboard often returns `null` for empty fields (like descriptions, dates, external uris). The `get` commands handle this gracefully by printing `None` instead of Go's native `<nil>`.
- **Color Support**: Tasks support Kanboard color IDs (`red`, `yellow`, `blue`, `green`, `purple`, etc.) via the `--color` flag in both `create` and `update` commands.
- **Output Formats**: Supports standard tab-separated/key-value output, as well as raw JSON output via the `--json` global flag.

## Build Instructions
To compile the binary locally:
```bash
go build -o kb-cli
```
