## ardi compile

Compile specified sketch

### Synopsis


Compile sketches for a specified board. You must provide the board FQBN, if left unspecified, a list of available choices will be be printed. If the sketch argument matches as user defined build in ardi.json, the values defined in build will be used to compile

```
ardi compile [sketch|build] [flags]
```

### Options

```
  -p, --build-prop stringArray   Specify build property to compiler
  -f, --fqbn string              Specify fully qualified board name
  -h, --help                     help for compile
  -s, --show-props               Show all build properties (does not compile)
  -w, --watch                    Watch sketch file for changes and recompile
```

### Options inherited from parent commands

```
  -g, --global        Use global data directory
      --port string   Set port for cli daemon
  -q, --quiet         Silence all logs
  -v, --verbose       Print all logs
```

### SEE ALSO

* [ardi](ardi.md)	 - Ardi is a command line build manager for arduino projects.

