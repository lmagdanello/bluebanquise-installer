# BlueBanquise Installer CLI

A Go CLI tool to automate the installation of [BlueBanquise](https://github.com/bluebanquise/bluebanquise), a coherent collection of Ansible roles designed to deploy and manage large groups of hosts (clusters of nodes).

## Usage

### Online Installation

To install BlueBanquise by downloading collections directly from GitHub:

```bash
sudo ./bluebanquise-installer online
```

Or with custom user settings:

```bash
sudo ./bluebanquise-installer online --user myuser --home /opt/bluebanquise
```

This will install collections using:
```bash
ansible-galaxy collection install git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master -vvv --upgrade
```

### Offline Installation

You can install BlueBanquise offline using pre-installed collections, tarball files, offline Python requirements, and core variables:

#### Using tarball files:
```bash
sudo ./bluebanquise-installer offline --collections-path /path/to/collections
```

**Note**: The download command downloads collection tarballs (`.tar.gz` files) that can be used for offline installation. Use `--collections-path` to specify the collections directory.

#### Download collections and tarballs:
```bash
# Download collection tarballs for offline installation
sudo ./bluebanquise-installer download --path /tmp/offline --collections

# Transfer tarball files to target machine
scp -r /tmp/offline user@target-machine:/tmp/

# Install on target machine with offline tarball files
sudo ./bluebanquise-installer offline --collections-path /tmp/offline/collections
```

#### Using offline Python requirements:
```bash
sudo ./bluebanquise-installer offline \
  --collections-path /path/to/collections \
  --requirements-path /path/to/requirements
```

#### Download Python requirements:
```bash
# Download Python requirements for offline installation
sudo ./bluebanquise-installer download --path /tmp/offline --requirements

# Transfer requirements to target machine
scp -r /tmp/offline user@target-machine:/tmp/

# Install on target machine with offline requirements
sudo ./bluebanquise-installer offline \
  --collections-path /path/to/collections \
  --requirements-path /tmp/offline/requirements
```

#### Download core variables:
```bash
# Download core variables for offline installation
sudo ./bluebanquise-installer download --path /tmp/offline --core-vars

# Transfer core variables to target machine
scp /tmp/offline/bb_core.yml user@target-machine:/tmp/

# Install on target machine with offline core variables
sudo ./bluebanquise-installer offline \
  --collections-path /path/to/collections \
  --core-vars-path /tmp/offline/core-vars/bb_core.yml
```

#### Complete offline installation:
```bash
sudo ./bluebanquise-installer offline \
  --collections-path /tmp/offline/collections \
  --requirements-path /tmp/offline/requirements \
  --core-vars-path /tmp/offline/core-vars/bb_core.yml \
  --user myuser \
  --home /opt/bluebanquise
```

#### Command options:

- `--collections-path, -c`: Path to BlueBanquise collections (directory)
- `--requirements-path, -r`: Path to Python requirements for offline installation
- `--core-vars-path, -v`: Path to core variables (bb_core.yml) for offline installation
- `--user, -u`: BlueBanquise username (default: bluebanquise)
- `--home, -H`: User home directory (default: /var/lib/bluebanquise)
- `--skip-environment, -e`: Skip environment configuration
- `--debug, -d`: Enable debug mode

**Note**: The `--requirements-path` and `--core-vars-path` are optional and can be used with the `--collections-path` method.

### Status Check

Check the installation status:

```bash
./bluebanquise-installer status
```

Or with custom user settings:

```bash
./bluebanquise-installer status --user myuser --home /opt/bluebanquise
```

### Example usage with custom user:

```bash
# Online installation with custom user
sudo ./bluebanquise-installer online \
  --user ansible-admin \
  --home /opt/ansible

# Offline installation with custom user
sudo ./bluebanquise-installer offline \
  --collections-path /tmp/offline/collections \
  --requirements-path /tmp/offline/requirements \
  --user ansible-admin \
  --home /opt/ansible

# Check status for custom user
./bluebanquise-installer status \
  --user ansible-admin \
  --home /opt/ansible
```

## Testing

This project includes comprehensive tests to ensure reliability and functionality:

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run tests with coverage report
make test-coverage

# Run tests in verbose mode
make test-v

# Run tests with race detection
make test-race
```

### Test Structure

- **Unit Tests**: Test individual functions and components
  - `internal/system/packages_test.go` - OS detection and package definitions
  - `internal/utils/check_test.go` - System prerequisites validation
  - `internal/bootstrap/user_test.go` - User creation and management
  - `internal/bootstrap/collections_test.go` - Collections and core variables installation
  - `cmd/root_test.go` - CLI command structure

- **Integration Tests**: Test complete workflows
  - `integration_test.go` - End-to-end installation flows

### Test Requirements

- **Unit Tests**: Can run without special privileges
- **Integration Tests**: Require root privileges for user creation tests
- **Network Tests**: Some tests require internet connectivity

### Code Quality

```bash
# Run linter
make lint

# Format code
make format

# Install development tools
make install-tools
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## CI/CD Pipeline

This project features a complete CI/CD pipeline using GitHub Actions to ensure code quality and reliability:

### Continuous Integration

The CI runs on every push and pull request, including:

- **Unit Tests**: Comprehensive test coverage
- **Integration Tests**: End-to-end workflow testing
- **Linting**: Code quality checks with golangci-lint
- **Multi-Platform Testing**: Tests on different architectures
- **Offline Installation Testing**: Validates offline installation workflows with Python 3.12

### Workflows

#### Main CI (`ci.yml`)
- **Test**: Unit tests, integration tests, and linting
- **Online Installation Tests**: Tests installation on Rocky Linux 9 with Python 3.12
- **Offline Installation Tests**: Tests offline installation on Rocky Linux 9 with Python 3.12
- **Tarball Installation Tests**: Tests using tarball files on Rocky Linux 9 with Python 3.12
- **Architecture Tests**: Tests on amd64 and arm64
- **Integration Tests**: Complete workflow validation
- **Build Release**: Automated builds for multiple platforms (Linux, macOS, Windows)

#### Release (`release.yml`)
- **Automated Releases**: Creates releases when tags are pushed
- **Multi-Platform Builds**: Builds for Linux, macOS, and Windows (amd64/arm64)
- **Checksums**: Generates SHA256 checksums for all binaries

### Key Improvements

- **Simplified Testing**: Focused on Rocky Linux 9 with Python 3.12 for consistent and reliable testing
- **Enhanced Offline Support**: Automatic download of `setuptools` and `wheel` packages for complete offline Python installation
- **Core Variables Integration**: Full support for offline core variables installation
- **Modern Dependencies**: Updated to use latest GitHub Actions and Go tooling

### Local Development

```bash
# Run CI checks locally
make ci

# Run CI checks without Docker
make ci-local

# Build release binaries
make release

# Run Docker-based tests
make test-docker
```

### Code Quality Tools

- **Pre-commit Hooks**: Automated checks before commits
- **GolangCI-Lint**: Comprehensive Go linting
- **Code Coverage**: Minimum 80% coverage required

## About BlueBanquise

BlueBanquise is a generic collection that can be adapted to any type of architecture (HPC clusters, university or enterprise infrastructure, Blender render farm, K8S cluster, etc). Special focus on scalability for very large clusters.

## Features

This CLI provides:

- **Online Installation**: Downloads and installs BlueBanquise directly from GitHub
- **Offline Installation**: Installs from pre-downloaded local collections, tarballs, and Python requirements
- **Automatic OS Detection**: Supports RHEL/CentOS/Rocky/AlmaLinux, Ubuntu, Debian, OpenSUSE
- **Automatic Configuration**: Creates user, Python virtual environment and necessary configurations
- **Multi-Distribution Support**: Specific configurations for each OS version
- **Custom User Support**: Configure custom username and home directory
- **Complete Offline Support**: Download collections and Python requirements for air-gapped environments
- **Core Variables Installation**: Automatically installs BlueBanquise core variables (bb_core.yml)
- **Enhanced Python Requirements**: Automatic inclusion of `setuptools` and `wheel` for complete offline Python package installation
- **Python 3.12 Support**: Optimized for Python 3.12 across all supported distributions

## Core Variables

BlueBanquise requires core variables to be installed in your inventory at `group_vars/all/` level. The installer automatically handles this by:

- **Online Mode**: Downloads `bb_core.yml` directly from the [BlueBanquise GitHub repository](https://github.com/bluebanquise/bluebanquise/blob/master/resources/bb_core.yml)
- **Offline Mode**: Copies the provided `bb_core.yml` file to the correct location

The core variables file contains essential configuration variables that BlueBanquise needs to function properly. You can also:

- Use the vars plugin at ansible-playbook execution: `ANSIBLE_VARS_ENABLED=ansible.builtin.host_group_vars,bluebanquise.infrastructure.core`
- Add it to your `ansible.cfg` file: `vars_plugins_enabled = ansible.builtin.host_group_vars,bluebanquise.infrastructure.core`

## Supported Distributions

| OS Family | Distribution | Tested Versions | Architectures |
|-----------|--------------|-----------------|---------------|
| Red Hat   | RHEL         | 7, 8, 9         | x86_64, aarch64 |
|           | Rocky Linux  | 8, 9            | x86_64, aarch64 |
|           | AlmaLinux    | 8, 9            | x86_64, aarch64 |
|           | CentOS       | 7, 8, Stream    | x86_64, aarch64 |
| Debian    | Ubuntu       | 20.04, 22.04, 24.04 | x86_64, arm64 |
|           | Debian       | 11, 12          | x86_64, arm64 |
| SUSE      | OpenSUSE Leap| 15.5, 15.6      | x86_64, aarch64 |
|           | SLES         | 15.6            | x86_64, aarch64 |

## Installation

### Prerequisites

- Go 1.24.3 or higher
- Root access or sudo for package installation

### Compilation

```bash
git clone https://github.com/lmagdanello/bluebanquise-installer.git
cd bluebanquise-installer
go build -o bluebanquise-installer
```

## Troubleshooting

### Common Issues

1. **Permission denied errors**: Run with sudo/root
2. **Package manager not found**: The installer supports apt-get, dnf, yum, and zypper
3. **Python not found**: Make sure python3 is installed and available in PATH
4. **Internet connectivity issues**: Use offline installation methods for air-gapped environments

### Logs

The installer logs all operations to `/var/log/bluebanquise/bluebanquise-installer.log`.

### Debug Mode

Enable debug mode for more verbose output:

```bash
sudo ./bluebanquise-installer online --debug
sudo ./bluebanquise-installer offline --debug --collections-path /path/to/collections
```

## License

MIT License - see the LICENSE file for details.

## Acknowledgments

- [BlueBanquise](https://github.com/bluebanquise/bluebanquise) - Main project
- [Ansible](https://www.ansible.com/) - Automation platform

