[![GoDoc](https://godoc.org/github.com/robgonnella/ardi?status.svg)](https://godoc.org/github.com/robgonnella/ardi)

# Ardi

Ardi is a command-line tool for compiling, uploading, and watching logs for
your usb connected arduino board. Ardi allows you to develop in an environment
you feel comfortable, without being forced to use arduino's web or desktop IDEs.

Ardi's `--watch` flag allows you to auto re-compile and upload on save, saving
you time and improving efficiency.

Ardi should work for all boards and platforms supported by arduino-cli.
Run `ardi init` to download all supported platforms and indexes to ensure
maximum board support. To initialize only for a specific platform, run
`ardi init <platformID>`. To see a list of supported platforms and associated
IDs, run `ardi platform list`. To see a list of all supported boards and their
associated platforms and fqbns run `ardi board list`.

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

  Install golang: https://golang.org/doc/install

  run:

    GO111MODULE=on go get github.com/robgonnella/ardi

  Note:<br/>
  This tool is based directly on a specific commit of [arduino-cli]. The exact
  commit can be found in the [go.mod](./go.mod) file of ardi. When installing
  via "go get" be sure to omit the "-u" flag to prevent updating dependencies
  as arduino-cli may have changed and could behave unpredictably with ardi.
  Also because ardi is written using go modules, you should include
  `GO111MODULE=on` when installing to ensure the proper versions of dependencies
  are used.

# Installing platforms for board detection

```bash
# from any directory
ardi init --verbose
# or to initialize only for a specified platform
ardi init --verbose <platform_id>
# to list the available platforms
ardi platform list <optional_search_param>
```

# Remove all installed platforms and data

```bash
# from any directory
ardi clean
```

# List all available platforms

```bash
# list all (alias: search)
ardi platform list
# filter output based on keyword
ardi platform list mega
```

# List all available boards and associated platforms and FQBNs

```bash
# list all
ardi board list
# filter output based on keyword
ardi board list mega
```

# Creating and uploading Sketches

There are two options for compiling and uploading sketches.
Both options require your sketch `.ino` file to be in a
directory that matches the `.ino` file name.</br>
e.g. `blink/blink.ino`

**Using a project level "sketches" directory:**

- create a sketches directory in your project folder
- add your sketch directory to the sketches directory</br>
  e.g. `<project>/sketches/blink/blink.ino`
- From the root of your project run
  `ardi go <name_of_sketch_directory> --verbose`</br>
  e.g. `ardi go blink --verbose`

**Using an absolute or relative path to sketch directory:**

- point ardi at any absolute or relative path to a
  sketch directory.</br>
  e.g. `ardi go ~/<project_root>/<project_sub_dir>/blink/ --verbose`

By default ardi will connect to the serial port and print
logs. Ardi will read the sketch file and attempt to
auto-detect the baud rate. To manually specify the baud
rate run:</br>
`ardi go <sketch_name> --baud <BAUD_RATE>`

For a list of all ardi options run: `ardi --help` or `ardi [command] --help`.

# Using ardi's "watch" feature

Ardi allows you to optionally watch a specified sketch file for changes and
auto re-compile and re-upload. Just add the `--watch` flag to the `ardi go`
command.

```bash
ardi go blink --watch
#or
ardi go <path_to_sketch_dir> --watch
```

# Compiling only (no upload)

Ardi allow you to choose to just compile your sketch and not upload anywhere.
To do this ardi still needs to know what board to compile for though. Supply
the compile command the fqbn for your board. If you don't know the fqbn of your
board, run the compile command with no --fqbn argument and ardi will present
you with a list of board names and their associated fqbns for each installed
platform.

    ardi compile <sketch_directory> --fqbn <board_fqbn>

# Adding Libraries

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
