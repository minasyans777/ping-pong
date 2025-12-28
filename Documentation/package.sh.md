# package.sh

`package.sh` creates a `.deb` or `.rpm` package from an existing binary.

---

## Requirements

### For `.deb` packages
- dpkg-dev

#### Install on Debian / Ubuntu (apt)

```bash
sudo apt update
sudo apt install dpkg-dev
```

### For `.rpm` packages
- rpm-build

### Install on RHEL / Fedora / Alma / Rocky (dnf)

```bash
sudo dnf install rpm-build
```

---

### Options

```text
-n         Specify the directory name for the package
-s         Specify the full path of the binary file
-v         Specify the package version
--deb      Create a .deb package
--rpm      Create a .rpm package
-h         Show this help message
```

All options are required except -h

---

### Example Usage

```bash
task package -- -n myapp -s ./bin/myapp -v 1.0.0 --deb
```
