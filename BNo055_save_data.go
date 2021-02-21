package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kpeu3i/bno055"
)

const (
	Pi = 3.14159265358979323846264338327950288419716939937510582097494459 // pi https://oeis.org/A000796
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

	f1, err := os.Create("acc_data.txt")
	if err != nil {
		panic(err)
	}

	f2, err := os.Create("gyro_data.txt")
	if err != nil {
		panic(err)
	}

	defer f1.Close()
	defer f2.Close()

	for {
		select {
		case <-signals:
			err := sensor.Close()
			if err != nil {
				panic(err)
			}
		default:

			// vector, err := sensor.Euler()
			// if err != nil {
			// 	panic(err)
			// }

			// fmt.Printf("\r*** Bearing =%5.3f, roll=%5.3f, tilt=%5.3f\n", vector.X, vector.Y, vector.Z)

			acc, err := sensor.Accelerometer()
			if err != nil {
				panic(err)
			}
			// fmt.Printf("\r*** Acc x =%5.3f, Acc y =%5.3f, Acc z=%5.3f\n", acc.X, acc.Y, acc.Z)

			// write accelerometer data to file:

			// s := strconv.FormatFloat(acc.X)

			sx := strconv.FormatFloat(float64(acc.X), 'f', -1, 32)
			sy := strconv.FormatFloat(float64(acc.Y), 'f', -1, 32)
			sz := strconv.FormatFloat(float64(acc.Z), 'f', -1, 32)

			_, err = f1.WriteString(sx + " " + sy + " " + sz + "\n")

			gyro, err := sensor.Gyroscope()
			if err != nil {
				panic(err)
			}

			sx = strconv.FormatFloat(float64(gyro.X), 'f', -1, 32)
			sy = strconv.FormatFloat(float64(gyro.Y), 'f', -1, 32)
			sz = strconv.FormatFloat(float64(gyro.Z), 'f', -1, 32)

			_, err = f2.WriteString(sx + " " + sy + " " + sz + "\n")

			// d2 := []byte{fmt.Sprintf("%f", acc.X), 111, 109, 101, 10}

			if err != nil {
				panic(err)
			}
			// fmt.Printf("%.6f %.6f %.6f\n", acc.X, acc.Y, acc.Z)

			// temperature, err := sensor.Temperature()
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Printf("*** Temperature: t=%d\n", temperature) // temp is int8

			// time.Sleep(100 * time.Millisecond)

		}

		time.Sleep(10 * time.Millisecond)
	}

	// Output:
	// *** Status: system=133, system_error=0, self_test=15
	// *** Revision: software=785, bootloader=21, accelerometer=251, gyroscope=15, magnetometer=50
	// *** Axis: x=0, y=1, z=2, sign_x=0, sign_y=0, sign_z=0
	// *** Temperature: t=27
	// *** Euler angles: x=2.312, y=2.000, z=91.688
}

// ParseMagnetometer converts mag vector int angle. ignores z
// func ParseMagnetometer(magVector *bno055.Vector) float64 {

// 	// angle = atan2(Y, X);

// 	xData := float64((*magVector).X)
// 	yData := float64((*magVector).Y)

// 	angle := math.Atan2(xData, yData)

// 	if angle >= 0 {
// 		angle = angle * (180.0 / Pi)
// 	} else {
// 		angle = (angle + 2.0*Pi) * (180.0 / Pi)
// 	}

// 	return angle
// }

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