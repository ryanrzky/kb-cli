# kb-cli

A fast, lightweight command-line interface tool written in Go for managing projects and tasks on Kanboard via its JSON-RPC API.

## Features

- **Project Management**: Create, list, and view Kanboard projects.
- **Task Management**: Create, list, view, update, and delete tasks.
- **Color Support**: Assign native Kanboard colors (red, yellow, blue, green, etc.) to tasks directly from the CLI.
- **Human-Readable Output**: Automatically maps IDs to actual column titles and assignee names for `task list`. Handles empty fields gracefully by returning `None` instead of native `<nil>` values.
- **JSON Export**: Supports the `--json` global flag on commands to output raw JSON for integration with `jq` or other automation scripts.

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed.

Clone the repository and build the binary:

```bash
git clone https://github.com/ryanrizky/kb-cli.git
cd kb-cli
go build -o kb-cli
```

You can optionally move the binary to your `$PATH`:
```bash
sudo mv kb-cli /usr/local/bin/
```

## Configuration

`kb-cli` requires a configuration file containing your Kanboard API credentials. Create a file named `.kb-cli.yaml` in your home directory (`~/.kb-cli.yaml`) with the following structure:

```yaml
url: "https://kanboard.example.com/jsonrpc.php"
user: "jsonrpc"
token: "YOUR_API_TOKEN"
```
> **Note:** Ensure your URL explicitly points to the JSON-RPC API endpoint (`/jsonrpc.php`), not just the base web interface URL.

## Usage

### Projects

**List all projects:**
```bash
./kb-cli project list
```

**Create a new project:**
```bash
./kb-cli project create --name "My New Project" --description "Optional description"
```

**Get project details by ID:**
```bash
./kb-cli project get 1
```

### Tasks

**List all active tasks for a specific project:**
```bash
./kb-cli task list --project-id 1
```

**Create a task:**
*(Includes optional column ID and color flags)*
```bash
./kb-cli task create -p 1 -t "Fix Login Bug" -d "Investigate auth error" -c 2 --color red
```

**Get task details:**
```bash
./kb-cli task get 123
```

**Update a task:**
```bash
./kb-cli task update 123 --title "Fix Login Bug - Resolved" --color green
```

**Delete a task:**
```bash
./kb-cli task delete 123
```

### Global Flags

- `--json`: Output the direct JSON API response.
- `--config`: Specify a custom config file path (defaults to `~/.kb-cli.yaml`).
