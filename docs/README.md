# Ardi

Ardi is a command-line tool for compiling, uploading, and watching logs for
your usb connected arduino board. Ardi allows you to develop in an environment
you feel comfortable, without being forced to use arduino's web or desktop IDEs.

Ardi's `--watch` flag allows you to auto re-compile and upload on save, saving
you time and improving efficiency.

Ardi should work for all boards and platforms supported by arduino-cli.
Run `ardi init` to download all supported platforms and index files to ensure
maximum board support. To initialize only for a specific platform, run
`ardi init <platform_id>` or `ardi init <platform_id@version>`. To see a list of
supported platforms and associated IDs, run `ardi platform list`. To see a list
of all supported boards and their associated platforms and fqbns run
`ardi board list`.
(Note board fqbn will only be filled in once platform is initialized)

Once initialized run `ardi go <sketch_dir> --watch --verbose` and ardi will try
to auto detect your board, compile your sketch, upload, watch for changes in
your sketch file, and re-compile and re-upload. You can also run,
`ardi compile <sketch_directory> --fqbn <board_fqbn>` to only compile and
skip uploading.

Ardi also includes a basic library manager. Run `ardi lib init` in your project
directory to initialize it as an ardi project directory. Once initialized,
you can use `ardi lib add <lib_name>` to add libraries,
`ardi lib remove <lib_name>`, `ardi lib install` to install missing libraries
defined in ardi.json, and `ardi lib search <searchFilter>` to search existing
libraries.

Ardi stores all its platform data in `~/.ardi/` to avoid any conflicts with
existing `arduino-cli` installations.

Use "ardi [command] --help" for more information about a command.

# Installation

Ardi can be installed with golang's `go get` for versions of go >= 1.12

```bash
## go version >= 1.12

# From outside of a module directory
GO111MODULE=on go get github.com/robgonnella/ardi

# From inside of a module directory
go get github.com/robgonnella/ardi
```

You can also download and install the pre-built binaries
[here](https://github.com/robgonnella/ardi/releases)

# Usage

## Installing Platforms

Ardi requires certain binaries to be downloaded before it can properly compile
sketches, detect connected boards, or list FQBNs (fully qualified board names)
for supported boards.

```bash
# from any directory
ardi init --verbose
# or to initialize only for a specified platform
ardi init --verbose <platform_id>
# or to initialize for a specific version of a platform
ardi init --verbose <platform_id@version>
# to list the available platforms
ardi platform list <optional_search_param>
```

Note: Unfortunately, at this time ardi cannot provide a list of available
versions for each platform. However, you can find some of the arduino
platforms (aka cores) listed
[here](https://github.com/arduino?utf8=%E2%9C%93&q=core&type=&language=). Click
on your desired platform / core, then click on the "releases" tab to see
a list of versions for that platform / core.

## Remove Installed Platforms and Data

```bash
# from any directory
ardi clean
```

## List Available Platforms

```bash
# list all (alias: search)
ardi platform list
# filter output based on keyword
ardi platform list mega
```

## List Board Info

Ardi will only be able to list FQBNs for boards that have their associated
platforms installed (initialized) already. Use the board list command
bellow to find the desired platform for any board, then run
`ardi init <platform_id>`, then run the board list command again to find your
boards FQBN (fully qualified board name).

```bash
# list all
ardi board list
# filter output based on keyword
ardi board list mega
```

## Uploading Sketches

Point ardi at any absolute or relative path to a sketch directory.

    ardi go ~/<project_root>/<project_sub_dir>/blink/ --verbose

By default ardi will connect to the serial port and print logs. Ardi will read
the sketch file and attempt to auto-detect the baud rate. To manually specify
the baud rate run:

    ardi go <path_to_sketch_dir> --baud <baud_rate>

Accepts the a custom build property flag - see
[Build Properties](#Build-Properties)

## Compiling Only (no upload)

Ardi allows you to compile your sketch without uploading. To do this ardi still
needs to know what board to compile for though. Supply the compile command the
fqbn for your board. If you don't know the fqbn of your board, run the compile
command with no --fqbn argument and ardi will present you with a list of board
names and their associated fqbns for each installed platform.

    ardi compile <sketch_directory> --fqbn <board_fqbn> --verbose

Accepts the a custom build property flag - see
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

## Using Ardi's "Watch" Feature

Ardi allows you to optionally watch a specified sketch file for changes and
auto re-compile and re-upload. Just add the `--watch` flag to the `ardi go`
command.

    ardi go <path_to_sketch_dir> --watch

## Adding Libraries

Ardi has a minimal library manager. Run `ardi lib init` in your project
directory and ardi will create the necessary config files to create repeatable
versioned local library downloads and installs for that project. Or you can use
ardi's global library directory. Unfortunately, for now, a project directory
must use one or the other but not both.

Available library commands

```bash
# initialize a project directory to use local libraries
# (does not use global directory at all)
ardi lib init
# add a library locally
ardi lib add <library_name>
# add a library to the global directory
ardi lib add <library_name> --global
# remove a local library
ardi lib remove <library_name>
# remove a global library
ardi lib remove <library_name> --global
# install all local library dependencies specified in ardi.json
ardi lib install
# search for available libraries (alias find, list)
ardi lib search <search_filter>
```

[arduino-cli]: https://github.com/arduino/arduino-cli
