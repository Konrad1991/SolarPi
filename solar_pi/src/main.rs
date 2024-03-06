#![allow(warnings)]

use core::time;
use std::thread;
use std::thread::sleep;
use std::time::Duration;

mod logging;
mod dht22;
const DHT20_ADDR: u16 = 0x38;

fn main() {
    let dht20 = dht22::DHT20::new(1, DHT20_ADDR).expect("Failed to initialize DHT20 sensor");

    if dht20.begin() {
        loop {
            let (temperature, humidity, crc_error) = dht20.get_temperature_and_humidity();
            sleep(time::Duration::from_millis(1500));
            if crc_error {
                println!("CRC Error: Data may be corrupted");
            } else {
                println!("Temperature: {:.1}°C", temperature);
                println!("Humidity: {:.1}%", humidity);
            }
        }
    } else {
        println!("Failed to initialize DHT20 sensor");
    }
}


