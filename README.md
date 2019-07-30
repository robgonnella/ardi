# Arduino Hacking

Ardi is a tool for compiling, uploading code, and watching 
logs for your usb conencted arduino board from command-
line. This allows you to develop in an environment you feel
comfortable in, without needing to use arduino's web or 
desktop IDEs.

This tools should work for the following boards:
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

run:

```bash
go get github.com/robgonnella/ardi
```

Hook up your arduino board via USB

### Creating and uploading Sketches

- create a sketches directory in your project folder
- add the name of your sketch as a another directory in 
  the sketches directory
- add an ".ino" file with the same name in this directory
  e.g. <project>/sketeches/blink/blink.ino
- run `ardi <name_of_sketch_directory>`
  (baud rate defaults to 9600 - e.g. ardi blink)
  
To use a different baud rate run:
`ardi <sketch_name> --baud <BAUD_RATE>`

By default ardi will connect to the serial port and print
logs. To ignore logs and only compile and upload run:
`ardi <sketch_name> --watch false`

For a list of ardi options run: `ardi --help`

### Adding Libraries

Create a "libraries" directory in sketches. Add your libs to this directory
and arduino-cli will automatically include them.
