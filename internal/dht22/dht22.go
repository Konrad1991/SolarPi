// sshfs konrad@192.168.1.222:/home/ /home/konrad/Documents/GitHub/SolarPi/remote
// fusermount -u /home/konrad/Documents/GitHub/SolarPi/remote
// ssh konrad@192.168.1.222
// gcc i2c.c -li2c

package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

const (
	I2C_SLAVE = 0x0703
)

func ioctl(fd, cmd, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, arg)
	if errno != 0 {
		return fmt.Errorf("ioctl failed: %v", errno)
	}
	return nil
}

type i2c struct {
	file *os.File
	addr int
	open bool
}

func newI2C(file_path string, addr int) (*i2c, error) {
	file, err := os.OpenFile(file_path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	if err := ioctl(file.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return nil, errors.New("Failed to acquire bus access and/or talk to slave: %v\n")
	}
	device := i2c{file: file, addr: addr, open: false}
	return &device, nil
}

func (device i2c) read(register byte) ([]byte, error) {
	if err := device.write([]byte{register}); err != nil {
		return nil, err
	}
	buf := make([]byte, 1)
	if _, err := device.file.Read(buf); err != nil {
		return buf, errors.New("Failed to read from the I2C device: %v\n")
	}
	if (buf[0] & 0x18) != 0x18 {
		return buf, errors.New("Bitmask 0x18 test was not successful")
	}

	return buf, nil
}

func (device i2c) read_buffer(register byte, len int) ([]byte, error) {
	if err := device.write([]byte{register}); err != nil {
		return nil, err
	}
	buf := make([]byte, len)
	if _, err := device.file.Read(buf); err != nil {
		return buf, errors.New("Failed to read from the I2C device: %v\n")
	}
	return buf, nil
}

func (device i2c) read_reg(reg byte, len int) ([]byte, error) {
	time.Sleep(10 * time.Millisecond)
	if err := device.write([]byte{reg}); err != nil {
		return nil, fmt.Errorf("failed to write register address 0x%x: %v", reg, err)
	}
	data, err := device.read_buffer(reg, len)
	if err != nil {
		return nil, fmt.Errorf("failed to read from register 0x%x: %v", reg, err)
	}
	return data, nil
}

func (device i2c) write(buffer []byte) error {
	if _, err := device.file.Write(buffer); err != nil {
		return errors.New("Failed to write to the I2C device: %v\n")
	}
	return nil
}

func (device *i2c) write_reg(reg byte, data []byte) error {
	if _, err := device.file.Write([]byte{reg}); err != nil {
		return fmt.Errorf("Failed to write register address: %v", err)
	}
	if _, err := device.file.Write(data); err != nil {
		return fmt.Errorf("Failed to write data: %v", err)
	}
	return nil
}

func (device i2c) calc_crc8(buffer []byte) byte {
	var crc byte = 0xFF
	for _, byte := range buffer {
		crc ^= byte
		for i := 0; i < 8; i++ {
			if (crc & 0x80) != 0 {
				crc = (crc << 1) ^ 0x31
			} else {
				crc <<= 1
			}
		}
	}
	return crc & 0xFF
}

func (device *i2c) get_T_H() (float32, float32, bool, error) {
	register := byte(0x71)
	addr := byte(0xac)
	content := make([]byte, 2)
	content[0] = 0x33
	content[1] = 0x00
	device.write_reg(addr, content)
	for {
		time.Sleep(80 * time.Millisecond)
		data, err := device.read_reg(register, 1)
		if err != nil {
			return 0.0, 0.0, false, err
		}
		if (data[0] & 0x80) == 0 {
			break
		}
	}
	data, err := device.read_reg(register, 7)
	if err != nil {
		return 0, 0, false, err
	}
	t_raw := ((int(data[3]) & 0xf) << 16) + (int(data[4]) << 8) + int(data[5])
	h_raw := ((int(data[3]) & 0xf0) >> 4) + (int(data[1]) << 12) + (int(data[2]) << 4)
	temperature := float32(t_raw)/5242.88 - 50.0
	humidity := float32(h_raw) / float32(0x100000) * 100.0
	crc_error := device.calc_crc8(data) != data[6]
	return temperature, humidity, crc_error, nil
}

func (device i2c) close() {
	device.file.Close()
}

func main() {
	device, err := newI2C("/dev/i2c-1", 0x38)
	defer device.close()
	if err != nil {
		fmt.Printf("Failed to open the bus: %v\n", err)
		return
	}

	time.Sleep(500 * time.Millisecond)
	register := byte(0x71)
	res, err := device.read(register)
	if err != nil {
		fmt.Printf("Failed to open the bus: %v\n", err)
		return
	}
	fmt.Printf("Read from I2C device: %v\n", res)

	for {
		t, h, _, err := device.get_T_H()
		if err != nil {
			fmt.Printf("Failed to get data: %v\n", err)
			return
		}
		fmt.Printf("Temperature = %f°C; Humidity = %f%% \n", t, h)
		time.Sleep(1 * time.Second)
	}
}
