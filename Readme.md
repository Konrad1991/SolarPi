<p align="justify">

# SolarPi

![SolarPi Logo](./logo/SolarPi_Logo.png)

🌱 **Project Status: Recently Initiated, Work in Progress**

## Overview 

The SolarPi project focuses on harnessing solar energy to power a Raspberry Pi while incorporating it into a versatile setup as both a backend server and a weather station. The following aspects encapsulate the project's hardware and software components:

## Hardware Implementation

We provide a comprehensive list of components, accompanied by a schematic plan and photos showcasing the final product. Additionally, the Raspberry Pi serves as both a backend server and a weather station in this setup.

### Hardware list:

- **Raspberry Pi 4** 
    * The Raspberry Pi 4 serves as the core computing unit for the SolarPi project. It handles various tasks, including operating as a backend server and weather station.
    * ([More Information](https://www.raspberrypi.com/products/raspberry-pi-4-model-b/))
- **Solar panel**
    * [Amazon Link](https://www.amazon.de/dp/B075X49XJS?tag=idealode-am-pk-21&ascsubtag=2024-02-28_45582a61d1533d4954e668a287c172221c40feb7d104ff73492b0fa5fe674e25&th=1)
    * 80 W, 12 V; monocrystalline solar panel, providing a sustainable and renewable power source for the Raspberry Pi setup.
- **Solar regulator**
    * The solar regulator, also known as a charge controller, manages the power flow from the solar panel to prevent overcharging the battery. It ensures optimal charging and extends the battery life.
- **Battery**
    * A bike battery is used and serves as the energy storage solution for the SolarPi project.
- **Stepdown converter**
    * The stepdown converter regulates the voltage output from the battery or solar panel to match the requirements of specific components.
       It ensures a consistent and appropriate voltage level for the Raspberry Pi.
- **Sensors**
    * GY-68 BMP180 barometric air pressure and temperature sensor. Ordered from Az Delivery 
    * DHT20 Digital temperature sensor and air humidity sensor with I2C interface 2.5V to 5.5V Compatible with Raspberry Pi board for DIY microelectronics projects. Ordered from Az Delivery

## Software Components

- Firstly, there's the code essential for Raspberry Pi operations, which includes functionalities like safely shutting down the Pi in case of low battery status, particularly crucial during the colder months, such as December and January, in Germany.
- Secondly, we delve into the code responsible for controlling various sensors measuring light, humidity, temperature, and more.
- Finally, we address the software dedicated to utilizing the Raspberry Pi as a backend server. This involves crafting code for both the Pi itself and our personal computers and smartphones. 

 ## Language Choice

Throughout the entire project, the language of choice for coding is Rust, ensuring robust and efficient implementation.


</p>
