#include "display_module.h"
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SH110X.h>

#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64
#define i2c_Address 0x3c // SH1106 default address

Adafruit_SH1106G display = Adafruit_SH1106G(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, -1);

void initDisplay() {
  display.begin(i2c_Address, true); // true = reset
  display.clearDisplay();
  display.setTextColor(SH110X_WHITE);
  display.display();
}

void updateDisplay(float temp, int bpm, int spo2) {
  display.clearDisplay();
  
  display.setTextSize(1);
  display.setCursor(0, 0);
  display.println("Health Data:");
  
  display.setCursor(0, 20);
  display.print("Temp: ");
  display.print(temp);
  display.println(" C");

  display.setCursor(0, 35);
  display.print("BPM: ");
  display.println(bpm);

  display.setCursor(0, 50);
  display.print("SpO2: ");
  display.print(spo2);
  display.println(" %");

  display.display();
}
