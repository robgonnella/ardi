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

if your board isn't in the above list you'll need to find which
arduino core supports that board and change line 13 in the Dockerfile
to install that core instead.

### Prereqs

Install docker

>follow instructions here: https://docs.docker.com/v17.12/install/

Install docker-compose

>follow instructions here: https://github.com/Yelp/docker-compose/blob/master/docs/install.md

Hook up your arduino board via USB

### Creating and uploading Sketches

- create a directory in "sketches" directory
- add a ".ino" file with your code in this directory
- ./upload <name_of_sketch_directory>
e.g. ./upload blink

### Without Docker AND without IDE

If you really don't want to use docker but you don't want to use their
IDE either, you can set up your local environment with arduino-cli and run all
the commands yourself.

- install golang
>https://golang.org/doc/install

- install arduino-cli tool - run this from outside of a go "mod" directory
```
go get -u github.com/arduino/arduino-cli
```

These next command update our cli tool with the necessary packages to
compile and upload to our board
```
arduino-cli core update-index
arduino-cli core install arduino:avr
# If all has gone well, the following should display the name
# of your board. If it shows "unknown" we likely didn't install
# the correct core package
arduino-cli board list
```

To compile and upload your sketch:
```
export BOARD="$(arduino-cli board list | awk 'FNR == 2 {print $1}')"
arduino-cli compile --fqbn $BOARD sketches/<sketch_dir>
arduino-cli upload -p /dev/ttyACM0 --fqbn $BOARD sketches/<sketch_dir>
```

To watch the program logs in terminal:
```
stty -F /dev/ttyACM0 9600 raw -clocal -echo
cat /dev/ttyACM0
```

You may need to set permissions on /dev/ttyACM0
```
sudo chmod a+rw /dev/ttyACM0
```

### Adding Libraries

Create a "libraries" directory in sketches. Add your libs to this directory
and arduino-cli will automatically include them.
