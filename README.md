![V1](https://github.com/robgonnella/ardi/workflows/V1/badge.svg)<br/>
![V2](https://github.com/robgonnella/ardi/workflows/V2/badge.svg)

# Ardi

Ardi is a command-line tool for ardiuno that enables you to properly version and
manage project builds, and provides tools to help facilitate the development
process.

Things ardi can fo for you:

- Manage versioned platforms and libraries on a per-project basis
- Store user defined build configurations with a mechanism for easily running
  consistent and repeatable builds.
- Enable running your builds in a CI pipeline
- Compile and upload to an auto discovered connected board
- Watche a sketch for changes and auto recompil / reupload to a connected board
- Print various info about platforms and boards
- Search and print available libraries and versions

Ardi should work for all boards and platforms supported by [arduino-cli].

Use "ardi help [command]" for more information about a command.

# Installation

Ardi can be installed with golang's `go get` for versions of go >= 1.12

```bash
## go version >= 1.12

# From outside of a module directory
GO111MODULE=on go get github.com/robgonnella/ardi/v2@latest

# From inside of a module directory
go get github.com/robgonnella/ardi/v2@latest
```

You can also download and install the pre-built binaries
[here](https://github.com/robgonnella/ardi/releases)

# Usage

Ardi requires certain packages to be downloaded before it can properly compile
sketches, detect connected boards, and perform other tasks. These packages can
be stored in a project data directory to isolate version specific dependencies
for multiple projects, or in a global data directory available to all projects.
Ardi defaults to a "project level" data directory, use the "--global" flag
to use the global data directory.

To initialize an ardi project directory run:

```bash
ardi project-init
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
ardi compile <build_name>
# Compile multiple builds in ardi.json
ardi compile <build_name1> <build_name2> <build_name3>
# Compile all builds in ardi.json
ardi compile --all
# Upload only (skips building/compiling)
ardi upload <name>
```

Documentation for all commands can be found in [docs directory][docs]

Documentaion for V1 can be found [here](./docs/README.md)

[arduino-cli]: https://github.com/arduino/arduino-cli
[docs]: ./v2/docs/ardi.md
