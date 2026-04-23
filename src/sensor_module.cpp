#include "sensor_module.h"
#include <Wire.h>
#include <Adafruit_MLX90614.h>
#include "MAX30105.h"

Adafruit_MLX90614 mlx;
MAX30105 particleSensor;

void initSensors() {
  Wire.begin(21, 22);
  mlx.begin();
  particleSensor.begin(Wire, I2C_SPEED_FAST);
  particleSensor.setup();
}

float getTemperature() {
  return mlx.readObjectTempC();
}

int getBPM() {
  // Placeholder for BPM calculation logic
  return 0;
}

int getSpO2() {
  // Placeholder for SpO2 calculation logic
  return 0;
}
