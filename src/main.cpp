#include <Arduino.h>
#include "sensor_module.h"
#include "display_module.h"
#include "network_module.h"

unsigned long previousMillis = 0;
const long interval = 1000;

void setup() {
  Serial.begin(115200);
  
  initSensors();
  initDisplay();
  
  // WiFi and WebSocket initialization (placeholders)
  connectWiFi("HEALY_AP", "healy123");
  initWebSocket("healy-server.local", 8080, "/ws");
}

void loop() {
  networkLoop();
  
  unsigned long currentMillis = millis();
  
  if (currentMillis - previousMillis >= interval) {
    previousMillis = currentMillis;
    
    float temp = getTemperature();
    int bpm = getBPM();
    int spo2 = getSpO2();
    
    // Update local display
    updateDisplay(temp, bpm, spo2);
    
    // Send telemetry via WebSocket
    String jsonTelemetry = "{\"temp\":" + String(temp) + ",\"bpm\":" + String(bpm) + ",\"spo2\":" + String(spo2) + "}";
    sendTelemetry(jsonTelemetry);
  }
}