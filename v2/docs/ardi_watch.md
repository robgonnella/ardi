## ardi watch

Compile, upload, and watch

### Synopsis


Compile and upload code to an arduino board. Simply pass the directory containing the .ino file as the first argument. Ardi will automatically watch your sketch file for changes and auto re-compile & re-upload for you. Baud will be automatically be detected from sketch file.

```
ardi watch [sketch] [flags]
```

### Options

```
  -p, --build-prop stringArray   Specify build property to compiler
  -h, --help                     help for watch
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

