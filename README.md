# ardpifi

# johnusher/ardpifi Summary

Mesh network with Raspis to control Arduino controlled LED strips.

LED code can be updated, compiled and flashed locally via USB from Raspi to Arduino.

Keyboard inputs on each Raspi will send LED pattern info and sync across all Raspis in the network.

We install everything now with /bin/bootstrap (not tested!!)

Mesh code from: https://github.com/siggy/ledmesh

## Hardware set-up
Connect Arduino Uno (also Nano clone tested) to Raspi 3 via USB
Run the go script, and keyboard numbers will control LED sequence on the NeoPixel strip.
LED sequence can be programmed on the Raspi, compiled using arduino-cli, and flashed from the Raspi.


## Upcoming attractions:
-integrate the mesh network to allow multiple Raspis to communicate and change the LED show, sync'd on all devices.
-accelerometer/ gyro control using I2C bus.


# thanks @siggy!

### OS

1. Download Raspian Lite: https://downloads.raspberrypi.org/raspbian_lite_latest
2. Flash `20XX-XX-XX-raspbian-stretch-lite.zip` using Etcher
3. Remove/reinsert flash drive
4. Add `ssh` and `bootstrap` files:
    ```bash
    touch /Volumes/boot/ssh
    cp bin/bootstrap /Volumes/boot/
    chmod a+x /Volumes/boot/bootstrap
    diskutil umount /Volumes/boot
    ```

### First Boot

```bash
ssh pi@raspberrypi.local
# password: raspberry

# change default password
passwd



# set quiet boot
sudo sed -i '${s/$/ quiet loglevel=1/}' /boot/cmdline.txt

# install packages
sudo apt-get update
sudo apt-get install -y git tmux vim dnsmasq hostapd

# set up wifi (note leading space to avoid bash history)
sudo tee --append /etc/wpa_supplicant/wpa_supplicant.conf > /dev/null << 'EOF'
network={
    ssid="<WIFI_SSID>"
    psk="<WIFI_PASSWORD>"
}
EOF

# set static IP address
sudo tee --append /etc/dhcpcd.conf > /dev/null << 'EOF'

# set static ip

interface eth0
static ip_address=192.168.1.164/24
static routers=192.168.1.1
static domain_name_servers=192.168.1.1

interface wlan0
static ip_address=192.168.1.164/24
static routers=192.168.1.1
static domain_name_servers=192.168.1.1
EOF

# reboot to connect over wifi
sudo shutdown -r now
```

```bash
# configure git
git config --global push.default simple
git config --global core.editor "vim"
git config --global user.email "you@example.com"
git config --global user.name "Your Name"

# disable services
sudo systemctl disable hciuart
sudo systemctl disable bluetooth
sudo systemctl disable plymouth

# remove unnecessary packages
sudo apt-get -y purge libx11-6 libgtk-3-common xkb-data lxde-icon-theme raspberrypi-artwork penguinspuzzle ntp plymouth*
sudo apt-get -y autoremove

sudo raspi-config nonint do_boot_behaviour B2 0
sudo raspi-config nonint do_boot_wait 1
sudo raspi-config nonint do_serial 1
```

# now you are setup! run bootstrap to download packages

/bootstrap
```

## Build and run (in ~/code/go/src/github.com/siggy/ledmesh)

```bash
go run JU_led_mesh.go
```


### Mesh network

On raspi #1 run
./runBATMAN.sh

On raspi #2 run
./runBATMAN2.sh

See file main.go, from https://github.com/siggy/ledmesh

Based on:

https://www.reddit.com/r/darknetplan/comments/68s6jp/how_to_configure_batmanadv_on_the_raspberry_pi_3/

see instructions here:

https://github.com/siggy/ledmesh/blob/master/bin/bootstrap

eg

sudo ifconfig bat0 172.27.0.1/16

and sudo ifconfig bat0 172.27.0.2/16


## Code

Note
All dependencies managed in `go.mod` now,
just add an import directive for any new depedency in your `*.go` files, and
`go run/build` should just handle it.


All code (and go) is installed via bootstrap.
Repos we install:
https://github.com/siggy/ledmesh.git
https://github.com/johnusher/ardpifi.git

https://github.com/d2r2/go-i2c.git
https://github.com/nsf/termbox-go.git

# Arduino CLI install

folow instructions here:
https://siytek.com/arduino-cli-raspberry-pi/

this additional command was needed:
arduino-cli core install arduino:avr

Note the directory for the Arudion project must have the same name as the main file ()

Tested with Aurdion Uno and Aurdion clone: "Nano V3 | ATMEL ATmega328P AVR Microcontroller | CH340-Chip".
The Uno shows on port ttyACM0 and the clone on ttyUSB.
NB only 1 from 2 clones works for me!

<del>
## add libraries:
arduino-cli lib search Adafruit_NeoPixel
</del>

in duino_src:

compile and flash:
Uno:
arduino-cli compile --fqbn arduino:avr:uno duino_src
arduino-cli upload -p /dev/ttyACM0 --fqbn arduino:avr:uno duino_src
Clone
arduino-cli compile --fqbn arduino:avr:diecimila:cpu=atmega328 duino_src
arduino-cli upload -p /dev/ttyUSB0 --fqbn arduino:avr:diecimila:cpu=atmega328 duino_src




## Run

```bash
go run jumain.go
```

Press any key to print to screen (and eventually send to arduino).

To exit, press "q" to exit termbox, and then ctrl-c to exit the program.

### Run without hardware

Run with hardware (serial, network) API calls mocked out:

```bash
go run JU_led_mesh.go -no-hardware -no-lcd
```
