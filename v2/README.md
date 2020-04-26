[![GoDoc](https://godoc.org/github.com/robgonnella/ardi?status.svg)](https://godoc.org/github.com/robgonnella/ardi) ![Build](https://codebuild.us-east-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiclFVTnRkYjY1QXlnTjdueC9rZFhhRjNXVzdaUlZhL0tpdW9wMlA4TU94MC9POU5lcGJkTE4rSDJkSTdhYWJianBPWDhXdzR4a2x3U1lZL1h1NEYzSzBBPSIsIml2UGFyYW1ldGVyU3BlYyI6ImUxUXhjTjhhZFlHb0Q5b3AiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)

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

Ardi should work for all boards and platforms supported by arduino-cli.

Use "ardi [command] --help" for more information about a command.

# Installation

Ardi can be installed with golang's `go get` for versions of go >= 1.12

```bash
## go version >= 1.12

# From outside of a module directory
GO111MODULE=on go get github.com/robgonnella/ardi@latest

# From inside of a module directory
go get github.com/robgonnella/ardi@latest
```

You can also download and install the pre-built binaries
[here](https://github.com/robgonnella/ardi/releases)

# Usage

Ardi requires certain packages to be downloaded before it can properly compile
sketches, detect connected boards, and perform other tasks. Ardi will check for
a project level data directory for these packages and if found it will use this
directory, othewise ardi will use a global data directory located in the
uses home directory, `~/.ardi`. To initialize a directory as an ardi project
directory run `ardi project init`.

Any command run from within an "arid initialized" project directory will
use the local data directory for that project.

To initialize an ardi project directory run:

```bash
ardi project init
```

## Platform & Board Commands

The following "platform" commands will operate on the global data directory if
run from outside of an ardi project directory and on the the local project
data directory if run from within an ardi project directory.

```bash
# list all available platforms and their ids
ardi platform list --all
# from any directory
ardi platform add <platform_ids>
# list all installed platforms
ardi platform list --installed
# remove installed platform
ardi platform remove <platfor_ids>
```

The above platform commands are also duplicated via the project
command

```bash
# add project level platform
ardi project add platform <platform_ids>
# remove project level platform
ardi project remove platform <platform_ids>
# list project level installed platforms
ardi project list platforms
```

To see a list of boards and their associated platforms run

```bash
ardi board platforms
```

After installing your desired platform, to see a list of boards and their
associated fqbns run:

```bash
ardi board fqbns
```

## Remove All Installed Platforms and Data

```bash
# from inside initialized project directory
ardi clean # removes only project level data

# from outside of project directory
ardi clean # removes global data directory
```

## Uploading Sketches

Point ardi at any absolute or relative path to a sketch directory.

    ardi go ~/<project_root>/<project_sub_dir>/blink/ --verbose

By default ardi will read the sketch file, auto-detect the baud rate, connect
to the serial port, and print logs.

Accepts custom build properties flag - see
[Build Properties](#Build-Properties)

## Using the "watch" Feature

Ardi allows you to optionally watch a specified sketch file for changes and
auto re-compile and re-upload. Just add the `--watch` flag to the `ardi go`
command.

    ardi go <path_to_sketch_dir> --watch

## Compiling Only (no upload)

Ardi allows you to compile your sketch without uploading. To do this ardi still
needs to know what board to compile for though. Supply the compile command the
fqbn for your board. If you don't know the fqbn of your board, run the compile
command with no --fqbn argument and ardi will present you with a list of board
names and their associated fqbns for each installed platform.

    ardi compile <sketch_directory> --fqbn <board_fqbn> --verbose

Accepts custom build properties flag - see
[Build Properties](#Build-Properties)

The compile command also allows you print all build properties by using the
`--show-props` or `-s` flag. When this flag is specified ardi will ONLY
print the build properties, it will not actually compile your sketch.

## Build Properties

You can specify custom build properties to the `compile` and `go` commands by
using the `--build-prop` or `-p` flag followed by the build property and value.
To specify multiple build properties just precede each property with the
`--build-prop` or `-p` flag.

    ardi compile <sketch_dir> -v -f <fqbn> \
    -p build.extra_flags="-DSOME_OPTION" \
    -p compiler.cpp.extra_flags="-std=c++11"

## Storing Builds in ardi.json

Ardi enables you to store custom build details in ardi.json which you can
then easily run via the `ardi project build` command.

To add a build either manually modify ardi.json or use `ardi project add build`

```bash
ardi project add build \
--name <name> \
--platform <platform_id> \
--fqbn <fqbn> \
--build-prop build.extra_flags="-DSOME_OPTION" \
--build-prop compiler.cpp.extra_flags="-std=c++11"

# To see list of all options
ardi project add build --help
```

To run stored builds

```bash
# Run a single build
ardi project build <name>
# Run multiple builds
ardi project build <name1> <name2> <name3>
# Run all builds
ardi project build
```

## Adding Libraries

Ardi has a built in library manager allowing you to use repeatable
versioned local libraries for your project.

```bash
# add a library
ardi lib add <library_name1> <libaray_name2>
# remove a local library
ardi lib remove <library_name1> <libaray_name2>
# install all local library dependencies specified in ardi.json
ardi lib install
# search for available libraries (alias find, list)
ardi lib search <search_filter>
```

The above commands are duplicated via the `ardi project` command too

```bash
ardi project add lib <library_name> ...
ardi project remove lib <library_name> ...
# lists installed libraries
ardi project list lib
```

[arduino-cli]: https://github.com/arduino/arduino-cli
