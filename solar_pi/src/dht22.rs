/*
analgous code to:
https://github.com/cjee21/RPi-DHT20/blob/main/DFRobot_DHT20.py
*/

use linux_embedded_hal as hal;
use std::error::Error;
use std::thread;
use std::thread::sleep;
use std::time::Duration;
use rppal::i2c::I2c;

pub struct DHT20 {
    i2c: I2c,
}

impl DHT20 {
    pub fn new(bus: i32, address: u16) -> Result<Self, rppal::i2c::Error> {
        let mut i2c = I2c::new()?;
        i2c.set_slave_address(address)?;
        Ok(DHT20 { i2c })
    }

    pub fn begin(&self) -> bool {
        sleep(Duration::from_millis(500)); // wait at least 100 ms
        let data = self.read_reg(0x71, 1).unwrap();
        (data[0] & 0x18) == 0x18
    }

    pub fn get_temperature_and_humidity(&self) -> (f32, f32, bool) {
        self.write_reg(0xac, &[0x33, 0x00]).unwrap();
        loop {
            sleep(Duration::from_millis(80));
            let data = self.read_reg(0x71, 1).unwrap();
            if (data[0] & 0x80) == 0 {
                break;
            }
        }
        let data = self.read_reg(0x71, 7).unwrap();
        let temperature_raw_data = ((data[3] as i32 & 0xf) << 16) + ((data[4] as i32) << 8) + data[5] as i32;
        let humidity_raw_data = (((data[3] & 0xf0) as u32) >> 4) + ((data[1] as u32) << 12) + ((data[2] as u32) << 4);
        let temperature = (temperature_raw_data as f32) / 5242.88 - 50.0;
        let humidity = (humidity_raw_data as f32) / 0x100000 as f32 * 100.0;
        let crc_error = self.calc_crc8(&data) != data[6];
        (temperature, humidity, crc_error)
    }

    fn calc_crc8(&self, data: &[u8]) -> u8 {
        let mut crc = 0xFF;
        for &byte in &data[..data.len() - 1] {
            crc ^= byte;
            for _ in 0..8 {
                if (crc & 0x80) != 0 {
                    crc = (crc << 1) ^ 0x31;
                } else {
                    crc <<= 1;
                }
            }
        }
        crc & 0xFF
    }

    fn write_reg(&self, reg: u8, data: &[u8]) -> Result<(), rppal::i2c::Error> {
        sleep(Duration::from_millis(10));
        self.i2c.block_write(reg, data)
    }

    fn read_reg(&self, reg: u8, len: usize) -> Result<Vec<u8>, rppal::i2c::Error> {
        sleep(Duration::from_millis(10));
        let mut buffer = vec![0; len];
        self.i2c.block_read(reg, &mut buffer)?;
        Ok(buffer)
    }
}
