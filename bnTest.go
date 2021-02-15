package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kpeu3i/bno055"
)

func main() {
	sensor, err := bno055.NewSensor(0x28, 1)
	if err != nil {
		panic(err)
	}

	err = sensor.UseExternalCrystal(true)
	if err != nil {
		panic(err)
	}

	status, err := sensor.Status()
	if err != nil {
		panic(err)
	}

	fmt.Printf("*** Status: system=%v, system_error=%v, self_test=%v\n", status.System, status.SystemError, status.SelfTest)

	revision, err := sensor.Revision()
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"*** Revision: software=%v, bootloader=%v, accelerometer=%v, gyroscope=%v, magnetometer=%v\n",
		revision.Software,
		revision.Bootloader,
		revision.Accelerometer,
		revision.Gyroscope,
		revision.Magnetometer,
	)

	axisConfig, err := sensor.AxisConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"*** Axis: x=%v, y=%v, z=%v, sign_x=%v, sign_y=%v, sign_z=%v\n",
		axisConfig.X,
		axisConfig.Y,
		axisConfig.Z,
		axisConfig.SignX,
		axisConfig.SignY,
		axisConfig.SignZ,
	)

	temperature, err := sensor.Temperature()
	if err != nil {
		panic(err)
	}

	fmt.Printf("*** Temperature: t=%v\n", temperature)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			err := sensor.Close()
			if err != nil {
				panic(err)
			}
		default:
			vector, err := sensor.Euler()
			if err != nil {
				panic(err)
			}

			fmt.Printf("\r*** Euler angles: x=%5.3f, y=%5.3f, z=%5.3f\n", vector.X, vector.Y, vector.Z)

			temperature, err := sensor.Temperature()
			if err != nil {
				panic(err)
			}
			fmt.Printf("*** Temperature: t=%d\n", temperature) // temp is int8

			magnetometer, err := sensor.Magnetometer()
			if err != nil {
				panic(err)
			}

			fmt.Printf("*** magnetometer: t=%v\n", magnetometer)

			bearing := ParseMagnetometer(magnetometer)

			fmt.Printf("*** bearing: t=%v\n", bearing)

			time.Sleep(2 * time.Second)

		}

		time.Sleep(100 * time.Millisecond)
	}

	// Output:
	// *** Status: system=133, system_error=0, self_test=15
	// *** Revision: software=785, bootloader=21, accelerometer=251, gyroscope=15, magnetometer=50
	// *** Axis: x=0, y=1, z=2, sign_x=0, sign_y=0, sign_z=0
	// *** Temperature: t=27
	// *** Euler angles: x=2.312, y=2.000, z=91.688
}

// func ParseDegrees(value string, direction string) (string, error) {

func ParseMagnetometer(magVector *bno055.Vector) float64 {

	// angle = atan2(Y, X);

	// invalid operation: (*magVector)[0] (type bno055.Vector does not support indexing)

	xData := float64((*magVector).X)
	yData := float64((*magVector).Y)

	angle := math.Atan2(xData, yData)

	return angle
}

// https://github.com/kpeu3i/bno055/blob/master/sensor.go

// type Vector struct {
// 	X float32
// 	Y float32
// 	Z float32
// }

// func (s *Sensor) Magnetometer() (*Vector, error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	x, y, z, err := s.readVector(bno055MagDataXLsb)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 1uT = 16 LSB
// 	vector := &Vector{
// 		X: float32(x) / 16,
// 		Y: float32(y) / 16,
// 		Z: float32(z) / 16,
// 	}

// 	return vector, nil
// }
