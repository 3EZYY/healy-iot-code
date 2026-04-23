#include <Arduino.h>
#include <ArduinoJson.h>
#include "sensor_module.h"
#include "display_module.h"
#include "network_module.h"
#include "audio_module.h"

unsigned long previousMillis = 0;
const long interval = 1000;

void setup() {
  Serial.begin(115200);
  
  // Initialize subsystems
  initSensors();
  initDisplay();
  initI2S();
  
  // WiFi and WebSocket initialization (placeholders)
  connectWiFi("HEALY_AP", "healy123");
  initWebSocket("healy-server.local", 8080, "/ws");
}

void loop() {
  // Continuous non-blocking tasks
  networkLoop();
  loopbackTest();
  
  unsigned long currentMillis = millis();
  
  if (currentMillis - previousMillis >= interval) {
    previousMillis = currentMillis;
    
    // 1) Fetch live sensor readings
    float temp = getTemperature();
    int bpm = getBPM();
    int spo2 = getSpO2();
    
    // 2) Update local OLED display
    updateDisplay(temp, bpm, spo2);
    
    // 3) Create JSON payload
    JsonDocument doc;
    doc["temp"] = temp;
    doc["bpm"] = bpm;
    doc["spo2"] = spo2;
    doc["status"] = "online";
    
    // 4) Serialize and send telemetry via WebSocket
    String jsonString;
    serializeJson(doc, jsonString);
    sendTelemetry(jsonString);
  }
}