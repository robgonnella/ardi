# Arduino Hacking

Ardi is a tool for compiling, uploading code, and watching
logs for your usb connected arduino board from command-line.
This allows you to develop in an environment you feel comfortable
in, without needing to use arduino's web or desktop IDEs.

This tools should work for all boards and platforms supported by arduino-cli:

## Prereqs
___

Install golang: https://golang.org/doc/install

run:

```bash
go get -u github.com/robgonnella/ardi
```

## Installing platforms for board detection
___

```bash
# from any directory
ardi init
```

## Remove all installed platforms and data
___

```bash
# from any directory
ardi clean
```

## Creating and uploading Sketches
___

There are two options for compiling and uploading sketches.
Both options require your sketch `.ino` file to be in a
directory that matches the `.ino` file name.</br>
e.g. `blink/blink.ino`

**Running from root of project directory:**

- create a sketches directory in your project folder
- add your sketch directory to the sketches directory</br>
  e.g. `<project>/sketches/blink/blink.ino`
- From the root of your project run
  `ardi <name_of_sketch_directory>`</br>
  e.g. `ardi blink`

**Running using an absolute or relative path to sketch:**

- point ardi at any absolute or relative path to a
  sketch directory.</br>
  e.g. `ardi ~/<project_root>/<project_sub_dir>/blink/`

By default ardi will connect to the serial port and print
logs. Ardi will read the sketch file and attempt to
auto-detect the baud rate. To manually specify the baud
rate run:</br>
`ardi <sketch_name> --baud <BAUD_RATE>`

To ignore logs and only compile and upload run:</br>
`ardi <sketch_name> --watch false`

For a list of all ardi options run: `ardi --help`

### Adding Libraries

Create a "libraries" directory at the same level as your sketch directory.
Add your libs to this directory and ardi will automatically include them.
