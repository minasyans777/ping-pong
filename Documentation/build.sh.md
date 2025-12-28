# build.sh

`build.sh` builds a binary from a Go application.

---

## Requirements

- **Go** must be installed and available in your PATH.

### Installing Go

#### On RHEL / Fedora / Alma / Rocky (dnf):

```bash
sudo dnf install golang
```

#### On Debian / Ubuntu (apt):

```bash
sudo apt update
sudo apt install golang
```

The script uses go build to compile the application.

---

## Options

```text
-s    Specify the full path of the program to build
-o    Specify the directory where to save the built program
-h    Show this help message

All options are required except -h
```
