# ardpifi

# johnusher/ardpifi Summary

Mesh network with Raspis to control Arduino controlled LED strips using GPS and accelerometers.

A webserver on the pi provides a bi-directional UI (currently only via ethernet SSH to the headless pi).

LED code can be updated, compiled and flashed locally via USB from Raspi to Arduino using the duino CLI.

Keyboard inputs on each Raspi will send LED pattern info and sync across all Raspis in the network.




## Hardware set-up

See shopping list section for details on components.

Connect Arduino Uno (also Nano clone tested) to Raspi 3 via USB. Programmable NeoPixelLED strips connected to Arduino and should be powered separately:


NEO-6M GPS module connect with GPIO serial (uses serial port /dev/ttyS0. In raspi-config settings, you need to disable serial console output, and disable bluetooth.). I have used an external antenna and a small ceramic antenna: both seem to work.

Bosch BNO055:  triaxial accelerometer, gyroscope, geomagnetic sensor.

OLED 128*64.

Both BNo055 and OLED connect via I2C. ie SDA and SCL are shared lines (pin 3 and 5), so connect in parallel. I powered  BNo055 with +3.3 V (pin 1), and OLED with +5 V.

WLAN dongle: not tested yet!  This messes up BATMAN as the dongle defaults to WLAN0.
To change this and make WLAN0 the built-in: see https://github.com/RaspAP/raspap-webgui/issues/335

Notes on I2c:

From https://github.com/kpeu3i/bno055/:

"it seems all versions of Raspberry Pi have an I²C bus hardware problem preventing them from working correctly with Bosch BNO055. The problem has been variously diagnosed as being due to the Pi’s inability to handle clock stretching in arbitrary parts of the I²C transaction and the BNO055 chip’s exquisite sensitivity to I²C bus levels."

Solutions:

Configuring software I²C driver

Raspbian has a software I²C driver that can be enabled by adding the following line to /boot/config.txt:

dtoverlay=i2c-gpio,bus=3
This will create an I²C bus called /dev/i2c-3. SDA will be on GPIO23 and SCL will be on GPIO24 which are pins 16 and 18 on the GPIO header respectively.

Slowing the I²C bus transactions

The solution require slowing the I²C bus transactions to 25 kb/s, by inserting a line in the /boot/config.txt file:

dtparam=i2c_arm_baudrate=25000



sudo i2cdetect -l

You will now see that i2c bus 3 also listed. Also run:

sudo i2cdetect -y 1

This should show OLED on 3c


sudo i2cdetect -y 3

This should show the BNo055 on 28


There are "no hardware" options for running the main Go file without these HW modules.

Run the go script, and keyboard numbers will send a command to the mesh network, eg controlling LED sequence on the LED strip.

LED sequences can be programmed on the Raspi, compiled using arduino-cli, and flashed from the Raspi.



## Hardware shopping list

Rasperry Pi 3 Model B+ testesd.

Arduino clone:
https://www.christians-shop.de/Nano-V3-developer-board-for-Arduino-IDE-ATMEL-ATmega328P-AVR-Microcontroller-CH340-Chip-Christians-Technikshop

https://www.amazon.de/-/en/gp/product/B078SBBST6/ref=ppx_yo_dt_b_asin_title_o02_s00?ie=UTF8&psc=1

NEO-6M GPS:
https://www.amazon.de/dp/B088LR3488/ref=pe_3044161_185740101_TE_item

OLED display:
https://www.amazon.de/-/en/gp/product/B01L9GC470/ref=ppx_od_dt_b_asin_title_s00?ie=UTF8&psc=1

Bosch BNO055:
https://www.amazon.de/-/en/gp/product/B072NLTPTJ/ref=ppx_yo_dt_b_asin_title_o01_s00?ie=UTF8&psc=1

BROgrammable LED strip WS2812:
https://www.mouser.de/ProductDetail/digilent/122-000/?qs=WbxR7jUW5e9xhU9oZFzZgA==&countrycode=DE&currencycode=EUR

WLAN dongle TP-Link TL-WN725N:
https://www.amazon.de/gp/product/B008IFXQFU/ref=ppx_yo_dt_b_asin_title_o03_s00?ie=UTF8&psc=1


## Upcoming attractions:
-send GPS on the mesh so devices can work out relative baring (using magnetometer compass).
-accelerometer/ gyro control using I2C bus.


# thanks @siggy!

### OS

Install Raspi lite from onlione or try this:

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


if we have an additional wlan usb:

sudo tee --append /etc/dhcpcd.conf > /dev/null << 'EOF'

# set static ip

interface wlan1
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

Run Bootstrap from https://github.com/johnusher/ardpifi/tree/master/bin
./bootstrap
```




### Mesh network

On raspi #1 run

./runBATMAN.sh

On raspi #2 run

./runBATMAN2.sh



## Code

Note
All dependencies managed in `go.mod` now,
just add an import directive for any new depedency in your `*.go` files, and
`go run/build` should just handle it.


All code (and go) is installed via bootstrap.



# Arduino CLI install

follow instructions here:
https://siytek.com/arduino-cli-raspberry-pi/

this additional command was needed:
arduino-cli core install arduino:avr

Note the directory for the Duino project must have the same name as the main file .ino.

Tested with Duino Uno and Duino clone: "Nano V3 | ATMEL ATmega328P AVR Microcontroller | CH340-Chip".
The Uno shows on port ttyACM0 and the clone on ttyUSB.


## add libraries:

arduino-cli lib search Adafruit_NeoPixel


in duino_src:

compile and flash:
Uno:
arduino-cli compile --fqbn arduino:avr:uno duino_src

arduino-cli upload -p /dev/ttyACM0 --fqbn arduino:avr:uno duino_src

Clone

arduino-cli compile --fqbn arduino:avr:diecimila:cpu=atmega328 duino_src

arduino-cli upload -p /dev/ttyUSB0 --fqbn arduino:avr:diecimila:cpu=atmega328 duino_src

arduino-cli upload -p /dev/ttyACM0 --fqbn arduino:avr:diecimila:cpu=atmega328 duino_src


## Run


On raspi #1 run

```bash
./runBATMAN.sh
go run JU_led_mesh.go -rasp-id=me --web-addr :8080 -log-level debug
```

On raspi #2 run

```bash
./runBATMAN2.sh
go run JU_led_mesh.go -rasp-id=poo --web-addr :8081 -no-lcd -log-level debug.
```

Press any key, sent to mesh, and if it is a 0 or 1, we change led pattern.

To exit, press "q" to exit termbox, and then ctrl-c to exit the program.

There are various flags if certain HW modules are not attached to the raspi:


### Flags

```bash
$ go run JU_led_mesh.go -h
  -log-level string
    	log level, must be one of: panic, fatal, error, warn, info, debug, trace (default "info")
  -no-batman
    	run without batman network
  -no-duino
    	run without arduino
  -no-gps
    	run without gps
  -no-oled
    	run without oled display
  -rasp-id string
    	unique raspberry pi ID (default "raspi 1")
  -web-addr string
    	address to serve web on (default ":8080")
```

### Run without hardware

Run with hardware (serial, network) API calls mocked out:

```bash
go run JU_led_mesh.go --rasp-id "raspi 1" --web-addr :8080 --no-batman --no-duino --no-gps --no-lcd --log-level debug
```

# Set up port forwarding for web server

## First-time ssh config setup

On laptop:


```bash
PI_IP=192.168.1.164
USER=pi

cat << EOF >> ~/.ssh/config

Host pi
  HostName $PI_IP
  User $USER
  Port 22
  BatchMode yes
  ServerAliveInterval 60
  ServerAliveCountMax 5
  ForwardAgent yes
    LocalForward 127.0.0.1:8080 127.0.0.1:8080
EOF
```

## Tunnel

```bash
ssh $USER@pi
```

Open a browser window to localhost:8080


## Other useful scripts
```bash
run BNo055_save_data.go
```
This saves accelerometer and gyroscope data to a txt file.

```bash
run OLEDtest.go
```
This displays a real-time clock on the OLED.

convAccData.m

MATLAB script to process accelerometer data.
