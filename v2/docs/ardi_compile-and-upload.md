## ardi compile-and-upload

Compiles then uploads to connected arduino board

### Synopsis


Compiles and uploads sketches for connected boards. If a connected board cannot be detected, you can provide the fqbn and port via command flags. If the sketch argument matches a user defined build in ardi.json, the values defined in build will be used to compile and upload.

```
ardi compile-and-upload [build|sketch] [flags]
```

### Aliases


```
build-and-upload
deploy
```

### Options

```
  -b, --baud int                 Specify baud rate when using "attach" flag
  -p, --build-prop stringArray   Specify build property to compiler
  -f, --fqbn string              Specify fully qualified board name
  -h, --help                     help for compile-and-upload
      --port string              The port your arduino board is connected to
  -s, --show-props               Show all build properties (does not compile)
```

### Options inherited from parent commands

```
  -q, --quiet     Silence all logs
  -v, --verbose   Print all logs
```

### SEE ALSO

* [ardi](ardi.md)	 - Ardi is a command line build manager for arduino projects.

