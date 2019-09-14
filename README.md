# Ardi

Ardi is a command-line tool for compiling, uploading, and watching logs for
your usb connected arduino board. Ardi allows you to develop in an environment
you feel comfortable, without being forced to use arduino's web or desktop IDEs.

Ardi's `--watch` flag allows you to auto re-compile and upload on save, saving
you time and improving efficiency.

Ardi should work for all boards and platforms supported by arduino-cli.
Run `ardi init` to download all supported platforms and indexes to ensure
maximum board support.

Once initialized run `ardi go <sketch_dir> --watch --verbose` and ardi will try
to auto detect your board, compile your sketch, upload, watch for changes in
your sketch file, and re-compile and re-upload.

Ardi stores all its data in a `.ardi` directory in the users home directory
to avoid any conflicts with existing `arduino-cli` installations.

Use "ardi [command] --help" for more information about a command.
___

## Prereqs

Install golang: https://golang.org/doc/install

run:

```bash
go get -v github.com/robgonnella/ardi
```
___
## Installing platforms for board detection

```bash
# from any directory
ardi init --verbose
```
___
## Remove all installed platforms and data

```bash
# from any directory
ardi clean
```
___
## Creating and uploading Sketches

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
___
## Using arid's "watch" feature

Ardi allows you to optionally watch a specified sketch file for changes and
auto re-compile and re-upload. Just add the `--watch` flag to the `ardi go`
command.

```bash
ardi go blink --watch
#or
ardi go <path_to_sketch_dir> --watch

```
___
### Adding Libraries

Create a "libraries" directory at the same level as your sketch directory.
Add your libs to this directory and ardi will automatically include them.
