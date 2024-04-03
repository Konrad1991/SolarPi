/*
 * dhtXX_sensor.c
 *
 * program to read DHTXX temperature and humidity sensors using pigpio library
 *
 * written by:  knute johnson
 * date:        11 april 2015
 * version:     0.94.2
 *
 * compile:     gcc -o sensor dhtXX_sensor.c -lrt -lpigpio -lpthread
 * usage:       sudo ./sensor [help] [DHTXX] [loop] [nn]
 *                  help    -   prints help message and exits
 *                  DHTXX   -   sensor type DHT11 or DHT22, default DHT22
 *                  loop    -   causes the program to run in a loop
 *                  nn      -   GPIO the sensor is attached to, default 4
 */

#include <stdio.h>
#include <strings.h>
#include <pigpio.h>

#define VERSION         "0.94.2"
#define DEFAULT_GPIO    4
#define DHT11           0
#define DHT22           1
#define DEFAULT_DEVICE  DHT22
#define LOOP            1
#define NO_LOOP         0

int bitCount;
int data[5];
int startTick;

void callback(int gpio, int level, uint32_t tick);

int main(int argc, char* argv[]) {
    int i;
    int device = DEFAULT_DEVICE;
    int gpio = DEFAULT_GPIO;
    int loopFlag = NO_LOOP;

    printf("dhtXX_sensor.c - v%s\n",VERSION);

    // parse command line arguments
    for (i=1; i<argc; i++) {
        if (!strcasecmp(argv[i],"help")) {
            puts(
             "Usage: sudo ./dhtXX_sensor [help] [DHTXX] [loop] [nnn] ");
            puts("    help  -   prints this message and exits");
            puts("    DHTXX -   sensor type DHT11 or DHT22, default DHT11");
            puts("    loop  -   causes the program to run in loop");
            puts("    nnn   -   GPIO the sensor is attached to (default 4)");
            return 0;
        }
        if (!strcasecmp(argv[i],"DHT11"))
            device = DHT11;
        if (!strcasecmp(argv[i],"DHT22"))
            device = DHT22;
        if (!strcasecmp(argv[i],"loop"))
            loopFlag = LOOP;
        int temp_gpio = atoi(argv[i]);
        if (temp_gpio > 1 && temp_gpio <= 31)
           gpio = temp_gpio; 
    }

    // initialize pigpio
    int version = gpioInitialise();
    if (version == PI_INIT_FAILED) {
        puts("PI_INIT_FAILED");
        return -1;
    } else {
        printf("pigpio - V%d\n",version);
    }

    // set callback function on gpio transition
    if (gpioSetAlertFunc(gpio,callback) != 0) {
        puts("PI_BAD_USER_GPIO");
        return -1;
    }

    // number of attempts to query sensor
    int tries = 3;

    // attempt to read sensor tries times or loop forever
    do {
        // start bit count less 3 non-data bits
        bitCount = -3;

        // zero the data array
        for (i=0; i<5; i++)
            data[i] = 0;

        // set start time of first low signal
        startTick = gpioTick();

        // set pin as output and make high for 50ms so we can detect first low
        gpioSetMode(gpio,PI_OUTPUT);
        gpioWrite(gpio,1);
        gpioDelay(50000);

        // send start signal
        gpioWrite(gpio,0);
        // wait for 18ms
        gpioDelay(18000);
        // return bus to high for 20us
        gpioWrite(gpio,1);
        gpioDelay(20);
        // change to input mode
        gpioSetMode(gpio,PI_INPUT);

        // wait 50ms for data input
        gpioDelay(50000);

        // if we received 40 data bits and the checksum is valid
        if (bitCount == 40 &&
         data[4] == ((data[0] + data[1] + data[2] + data[3]) & 0xff)) {
            float tempC;
            float humidity;

            if (device == DHT11) {
                humidity = data[0] + data[1] / 10.0f;
                tempC = data[2] + data[3] / 10.0f;
            } else if (device == DHT22) {
                humidity = (data[0] * 256 + data[1]) / 10.0f;
                tempC = ((data[2] & 0x7f) * 256 + data[3]) / 10.0f;
                // check for negative temp bit
                if (data[2] & 0x80)
                    tempC *= -1.0f;
            }
            float tempF = 9.0f * tempC / 5.0f + 32.0f;
            printf("Temperature: %.1fC %.1fF  Humidity: %.1f%%\n",
             tempC,tempF,humidity);
            // we're done
            tries = 0;
        } else {
            puts("Data Invalid!");
            --tries;
        }

        // minimum device reset time, 2 seconds
        gpioSleep(PI_TIME_RELATIVE,2,0);
    } while (loopFlag || tries) ;

    // shutdown pigpio
    gpioTerminate();

    return 0;
}

// level change call back function
void callback(int gpio, int level, uint32_t tick) {
    // if the level has gone low
    if (level == 0) {
        // duration is the elapsed time between lows
        int duration = tick - startTick;
        // set the timer start point to this low
        startTick = tick;

        // if we have seen the first three lows which aren't data
        if (++bitCount > 0) {
            // point into data structure, eight bits per array element
            // shift the data one bit left
            data[(bitCount-1)/8] <<= 1;
            // set data bit high if elapsed time greater than 100us
            data[(bitCount-1)/8] |= (duration > 100 ? 1 : 0);
        }
    }
}
