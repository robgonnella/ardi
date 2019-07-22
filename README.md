# Arduino Hacking

The following set of instructions deals with compiling and uploading programs
to your usb conencted arduino board without needing to use their web or
desktop IDEs. Add your sketch to the sketches directory and
run `./upload <sketch_name>`

This script should work for the following boards:
<table>
  <th>
    <tr>
      <td>Board Name</td>
      <td>FQBN</td>
    </tr>
  </th>
  <tbody>
    <tr>
      <td>Adafruit Circuit Playground</td>
      <td>arduino:avr:circuitplay32u4cat</td>
    <tr>
    <tr>
      <td>Arduino BT</td>
      <td>arduino:avr:bt</td>
    <tr>
    <tr>
      <td>Arduino Duemilanove or Diecimila</td>
      <td>arduino:avr:diecimila</td>
    <tr>
    <tr>
      <td>Arduino Esplora</td>
      <td>arduino:avr:esplora</td>
    <tr>
    <tr>
      <td>Arduino Ethernet</td>
      <td>arduino:avr:ethernet</td>
    <tr>
    <tr>
      <td>Arduino Fio</td>
      <td>arduino:avr:fio</td>
    <tr>
    <tr>
      <td>Arduino Gemma</td>
      <td>arduino:avr:gemma</td>
    <tr>
    <tr>
      <td>Arduino Industrial 101</td>
      <td>arduino:avr:chiwawa</td>
    <tr>
    <tr>
      <td>Arduino Leonardo</td>
      <td>arduino:avr:leonardo</td>
    <tr>
    <tr>
      <td>Arduino Mega ADK</td>
      <td>arduino:avr:megaADK</td>
    <tr>
    <tr>
      <td>Arduino Mini</td>
      <td>arduino:avr:mini</td>
    <tr>
    <tr>
      <td>Arduino NG or older</td>
      <td>arduino:avr:atmegang</td>
    <tr>
    <tr>
      <td>Arduino Nano</td>
      <td>arduino:avr:nano</td>
    <tr>
    <tr>
      <td>Arduino Pro or Pro Mini</td>
      <td>arduino:avr:pro</td>
    <tr>
    <tr>
      <td>Arduino Robot Control</td>
      <td>arduino:avr:robotControl</td>
    <tr>
    <tr>
      <td>Arduino Robot Motor</td>
      <td>arduino:avr:robotMotor</td>
    <tr>
    <tr>
      <td>Arduino Uno WiFi</td>
      <td>arduino:avr:unowifi</td>
    <tr>
    <tr>
      <td>Arduino Yún</td>
      <td>arduino:avr:yun</td>
    <tr>
    <tr>
      <td>Arduino Yún Mini</td>
      <td>arduino:avr:yunmini</td>
    <tr>
    <tr>
      <td>Arduino/Genuino Mega or Mega 2560</td>
      <td>arduino:avr:mega</td>
    <tr>
    <tr>
      <td>Arduino/Genuino Micro</td>
      <td>arduino:avr:micro</td>
    <tr>
    <tr>
      <td>Arduino/Genuino Uno</td>
      <td>arduino:avr:uno</td>
    <tr>
    <tr>
      <td>LilyPad Arduino</td>
      <td>arduino:avr:lilypad</td>
    <tr>
    <tr>
      <td>LilyPad Arduino USB</td>
      <td>arduino:avr:LilyPadUSB</td>
    <tr>
    <tr>
      <td>Linino One</td>
      <td>arduino:avr:one</td>
    <tr>
  </tbody>
</table>

### Prereqs

Install golang: https://golang.org/doc/install
Set GOPATH environment variable: https://github.com/golang/go/wiki/SettingGOPATH

run:

```bash
go get github.com/robgonnella/ardie
```

Hook up your arduino board via USB

### Creating and uploading Sketches

- create a directory in "sketches" directory
- add a ".ino" file with your code in this directory
- ardie <name_of_sketch_directory> <baud_rate>
- baud rate defaults to 9600
e.g. ardie blink

### Adding Libraries

Create a "libraries" directory in sketches. Add your libs to this directory
and arduino-cli will automatically include them.
