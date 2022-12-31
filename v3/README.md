![V3](https://github.com/robgonnella/ardi/workflows/V3/badge.svg)
[![codecov](https://codecov.io/gh/robgonnella/ardi/branch/main/graph/badge.svg?token=YMMFDIKX3A)](https://codecov.io/gh/robgonnella/ardi)

# Ardi

Ardi is a command-line tool built around ardiuno-cli that enables you to
properly version and manage project dependencies, and provides a simple
interface, similar to npm, for creating repeatable builds

Things ardi can fo for you:

- Manage versioned platforms and libraries on a per-project basis
- Store user defined build configurations with a mechanism for easily running
  consistent and repeatable builds.
- Enable running your builds in a CI pipeline
- All features supported by arudino-cli

Ardi should work for all boards and platforms supported by [arduino-cli].

Use "ardi help [command]" for more information about a command.

# Installation

Ardi can be installed with golang's `go get` for versions of go >= 1.12

```bash
## go version >= 1.12 < 1.18
GO111MODULE=on go get github.com/robgonnella/ardi/v3@latest

## go version >= 1.18
go install github.com/robgonnella/ardi/v3@latest
```

You can also download and install the pre-built binaries
[here](https://github.com/robgonnella/ardi/releases)

**Note: Linux users may need to add their user to the `dialout` group**

```
sudo usermod -aG dialout $(whoami)
```

# Usage

Ardi requires certain packages to be downloaded before it can properly compile
sketches, detect connected boards, and perform other tasks. These packages are
stored in a project data directory to isolate version specific dependencies
for multiple projects.

To initialize an ardi project directory run:

```bash
ardi init
```

## Storing Builds in ardi.json

Ardi enables you to store custom build details in ardi.json which you can
then easily run via the `ardi compile`, and `ardi upload`
commands.

To add a build either manually modify ardi.json or use `ardi add build`

```bash
ardi add build \
--name <name> \
--sketch <path_to_sketch_or_directory> \
--fqbn <fqbn> \
--build-prop build.extra_flags="-DSOME_OPTION" \
--build-prop compiler.cpp.extra_flags="-std=c++11"
```

To run stored builds

```bash
# Compile single build in ardi.json
ardi build <build_name>
# Compile multiple builds in ardi.json
ardi build <build1> <build2>
# Compile all
ardi build --all
```

Documentation for all commands can be found in [docs directory][docs]

[arduino-cli]: https://github.com/arduino/arduino-cli
[docs]: ./v2/docs/ardi.md