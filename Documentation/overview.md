# Overview

This repository automates the full lifecycle of package management:

1. Build a binary from a Go application
2. Package the binary as `.deb` or `.rpm`
3. Add the package to a repository and generate signed metadata

All steps are executed using Taskfile to provide a consistent and simple interface.
