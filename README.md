# Package Build & Repository Automation

This repository provides scripts for automating the following workflow:

- Building a binary from a Go application
- Creating `.deb` or `.rpm` packages from an existing binary
- Adding new packages to a repository, generating metadata, and signing releases with GPG

All scripts are executed **via Taskfile** and are designed to be used as a unified automation interface.

---

Scripts Overview

`build.sh`
Builds a binary for a Go-based application.

`package.sh`
Takes an existing binary file and creates either a .deb or .rpm package.

`package-sync.sh`
Adds a new package to a repository, generates repository metadata (Packages, Release, InRelease), and performs GPG signing.

---

## Documentation

- [Overview](Documentation/overview.md)
- [System Requirements](Documentation/requirements.md)
- [Taskfile Configuration](Documentation/taskfile.md)
- [build.sh](Documentation/build.sh.md)
- [package.sh](Documentation/package.sh.md)
- [package-sync.sh](Documentation/package-sync.sh.md)
