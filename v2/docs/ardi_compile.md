## ardi compile

Compile specified sketch or build(s)

### Synopsis


Compile sketches and builds for specified boards. When compiling for a sketch, you must provide the board FQBN. If left unspecified, a list of available choices will be be printed. If the sketch argument matches a user defined build in ardi.json, the values defined in build will be used to compile

```
ardi compile [sketch|build(s)] [flags]
```

### Aliases


```
build
```

### Options

```
  -a, --all                      Compile all builds specified in ardi.json
  -p, --build-prop stringArray   Specify build property to compiler
  -f, --fqbn string              Specify fully qualified board name
  -h, --help                     help for compile
  -s, --show-props               Show all build properties (does not compile)
  -w, --watch                    Watch sketch file for changes and recompile
```

### Options inherited from parent commands

```
  -q, --quiet     Silence all logs
  -v, --verbose   Print all logs
```

### SEE ALSO

* [ardi](ardi.md)	 - Ardi is a command line build manager for arduino projects.

