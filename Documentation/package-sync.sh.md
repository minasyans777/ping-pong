# package-sync.sh

`package-sync.sh` adds a new package to a repository and generates signed repository metadata.

---

## Requirements

### Common
- gpg
- A generated GPG key available on the system

### For `.deb` repositories
- dpkg-dev

#### Install on Debian / Ubuntu (apt)

```bash
sudo apt update
sudo apt install dpkg-dev gpg
```

### For `.rpm` repositories

- createrepo
- rpm-sign
- gpg


### Install on RHEL / Fedora / Alma / Rocky (dnf)

```bash
sudo dnf install createrepo rpm-sign gpg
```

---

### Options

```text
-r                Specify the full path of the repository
--arch            Specify the architecture
-s                Specify the full path of the new package
-d                Specify the full paths of the packages inside the repository
-m                Specify where the Packages file should be generated
-v                Specify the package version
-R                Specify where the Release file should be generated
-k                Specify the GPG key ID
-I                Specify where the InRelease file should be generated
-h|--help         Show this help message
```

---

### Required Options by Package Type

### RPM repositories

Required options:

```text
-r
-s
-k
-d
```

### DEB repositories

All options are required.

