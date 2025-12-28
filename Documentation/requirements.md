# System Requirements

This document describes required system dependencies for each script.

---

## build.sh Requirements

- Go (installed and available in PATH)

The script uses `go build` to compile the application.

---

## package.sh Requirements

### For `.deb` packages
- dpkg-dev

### For `.rpm` packages
- rpm-build

---

## package-sync.sh Requirements

### Common
- gpg
- A generated GPG key available on the system

### For `.deb` repositories
- dpkg-dev

### For `.rpm` repositories
- createrepo
- rpm-sign

