## ardi attach-and-watch

Compile, upload, attach, and watch for changes

### Synopsis


Compile, upload, watch board logs, and watch for sketch changes. Updates to the sketch file will trigger automatic recompile, reupload, and restarts the board log watcher. If the sketch argument matches a user defined build in ardi.json, the build values will be used for compilation, upload, and watch path

```
ardi attach-and-watch [sketch|build] [flags]
```

### Aliases


```
develop
dev
```

### Options

```
  -b, --baud int                 Specify baud rate
  -p, --build-prop stringArray   Specify build property to compiler
  -f, --fqbn string              Specify fully qualified board name
  -h, --help                     help for attach-and-watch
      --port string              The port your arduino board is connected to
```

### Options inherited from parent commands

```
  -q, --quiet     Silence all logs
  -v, --verbose   Print all logs
```

### SEE ALSO

* [ardi](ardi.md)	 - Ardi is a command line build manager for arduino projects.

