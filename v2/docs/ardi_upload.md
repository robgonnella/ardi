## ardi upload

Upload pre-compiled sketch build to a connected board. If the sketch argument matches a user defined build in ardi.json, the build values will be used to find the appropraite build to upload

### Synopsis


Upload pre-compiled sketch build to a connected board

```
ardi upload [sketch-dir|build] [flags]
```

### Options

```
  -h, --help   help for upload
  -l, --log    Watch board logs after uploading
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

