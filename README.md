# Aruduino Hacking

The following set of instructions deals with compiling and uploading programs
to your usb conencted arduino board without needing to use their web or
desktop IDEs. Add your sketch to the sketches directory and
run `./upload <sketch_name>`

This script should work for the following boards:
Board Name                       	FQBN
Adafruit Circuit Playground      	arduino:avr:circuitplay32u4cat
Arduino BT                       	arduino:avr:bt
Arduino Duemilanove or Diecimila 	arduino:avr:diecimila
Arduino Esplora                  	arduino:avr:esplora
Arduino Ethernet                 	arduino:avr:ethernet
Arduino Fio                      	arduino:avr:fio
Arduino Gemma                    	arduino:avr:gemma
Arduino Industrial 101           	arduino:avr:chiwawa
Arduino Leonardo                 	arduino:avr:leonardo
Arduino Leonardo ETH             	arduino:avr:leonardoeth
Arduino Mega ADK                 	arduino:avr:megaADK
Arduino Mini                     	arduino:avr:mini
Arduino NG or older              	arduino:avr:atmegang
Arduino Nano                     	arduino:avr:nano
Arduino Pro or Pro Mini          	arduino:avr:pro
Arduino Robot Control            	arduino:avr:robotControl
Arduino Robot Motor              	arduino:avr:robotMotor
Arduino Uno WiFi                 	arduino:avr:unowifi
Arduino Yún                      	arduino:avr:yun
Arduino Yún Mini                 	arduino:avr:yunmini
Arduino/Genuino Mega or Mega 2560	arduino:avr:mega
Arduino/Genuino Micro            	arduino:avr:micro
Arduino/Genuino Uno              	arduino:avr:uno
LilyPad Arduino                  	arduino:avr:lilypad
LilyPad Arduino USB              	arduino:avr:LilyPadUSB
Linino One                       	arduino:avr:one

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