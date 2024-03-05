use embedded_hal::i2c::{I2c, Error};

const ADDR: u8 = 0x38;
pub struct TemperatureSensorDriver<I2C> {
    i2c: I2C,
}

impl<I2C: I2c> TemperatureSensorDriver<I2C> {
    pub fn new(i2c: I2C) -> Self {
        Self { i2c }
    }

    pub fn read_temperature(&mut self) -> Result<u8, I2C::Error> {
        let mut temp = [0];
        let TEMP_REGISTER = 0x01;
        self.i2c.write_read(ADDR, &[TEMP_REGISTER], &mut temp)?;
        Ok(temp[0])
    }
}
