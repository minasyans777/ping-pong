# Taskfile Configuration

Taskfile is used as the main execution interface for all scripts.

Each task simply forwards CLI arguments to the underlying script using `{{.CLI_ARGS}}`.
This allows users to pass all required options directly to the scripts without modifying the Taskfile.

---

## Example Taskfile

```yaml

version: "3"

tasks:
  build:
    cmds:
      - ./scripts/build.sh {{.CLI_ARGS}}

  package:
    cmds:
      - ./scripts/package.sh {{.CLI_ARGS}}

  package-sync:
    cmds:
      - ./scripts/package-sync.sh {{.CLI_ARGS}}
```

### Example Usage

```bash
task build -- -s /path/to/go/project -o ./bin
```

