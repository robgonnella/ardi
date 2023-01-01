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

```bash
sudo usermod -aG dialout $(whoami)
```

# Usage

Ardi requires certain packages to be downloaded before it can properly compile
sketches. These packages are stored in a project data directory to isolate
version specific dependencies for multiple projects.

To initialize an ardi project directory run:

```bash
ardi init
```

## Adding Project Platforms

```bash
ardi add platform arduino:avr
# add specific version
ardi add platform arduino:avr@<version>
```

## Adding Project Libraries

```bash
ardi add library "Adfruit Pixie"
# add specific version
ardi add library "Adafruit Pixie"@<version>
```

## Storing Builds in ardi.json

Ardi enables you to store custom build details in ardi.json which you can
then easily run via the `ardi build` command.

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

## Executing arduino-cli commands

Ardi wraps arduino-cli via an "exec" command. This allows you to run any
arduino-cli command using the specific libraries and version defined in
`ardi.json`

```bash
ardi exec -- arduino-cli upload <build_dir>

# for help on all arduino-cli commands
ardi exec -- arduino-cli help
# help on specific arduino-cli command
ardi exec -- arduino-cli upload --help
```

Documentation for all commands can be found in [docs directory][docs]


# A Note about V2

Ardi V2 supported more features than V3; however, some of those features have
been, or are being, incorported in aruidno-cli. In an effort to reduce the
scope of this project, and complement arduino-cli without conflict, V3 has
been paired down to only focus on management of platorm and library versions,
and build configurations. All other features are supported by proxing
arduino-cli via `ardi exec -- arduino-cli ...`

I will do my best to continue to support V2 in terms of bug fixes but I do not
intend to add any additional features to V2.

All V2 docs and info can be found [here][docsV2].

[arduino-cli]: https://github.com/arduino/arduino-cli
[docs]: ./v3/docs/ardi.md
[docsV2]: ./v2/docs/ardi.md
