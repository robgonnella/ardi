## ardi upload

Upload pre-compiled sketch build to a connected board

### Synopsis


Upload pre-compiled sketch build to a connected board. If the sketch argument matches a user defined build in ardi.json, the build values will be used to find the appropraite build to upload

```
ardi upload [sketch-dir|build] [flags]
```

### Options

```
  -a, --attach        Attach to board port and print logs
  -f, --fqbn string   The FQBN of the board you want to upload to
  -h, --help          help for upload
  -p, --port string   The port your arduino board is connected to
```

### Options inherited from parent commands

```
  -q, --quiet     Silence all logs
  -v, --verbose   Print all logs
```

### SEE ALSO

* [ardi](ardi.md)	 - Ardi is a command line build manager for arduino projects.

