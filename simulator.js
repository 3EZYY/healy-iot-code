const WebSocket = require('ws');

// Konfigurasi URL WebSocket Backend di VPS (Sesuaikan IP jika perlu)
const WS_URL = 'ws://10.35.96.208:8080/ws/device';

console.log(`Connecting to HEALY Backend Simulator at ${WS_URL}...`);

// Sambungkan ke WebSocket sebagai "Device"
const ws = new WebSocket(WS_URL, {
  headers: {
    'device_id': 'healy-001'
  }
});

ws.on('open', () => {
  console.log('✅ Terhubung ke Backend!');
  console.log('📡 Mulai mengirim data dummy realtime setiap 2 detik...');

  // Kirim data setiap 2 detik
  setInterval(() => {
    // Generate data random (mirip dengan mock data sebelumnya)
    const payload = {
      device_id: "healy-001",
      temperature: 36.5 + (Math.random() * 2 - 1), // antara 35.5 - 37.5
      bpm: Math.floor(70 + Math.random() * 30),    // antara 70 - 100
      spo2: Math.floor(95 + Math.random() * 5),    // antara 95 - 99
      status: "ONLINE",
      timestamp: new Date().toISOString()
    };

    console.log(`[TX] Mengirim: T=${payload.temperature.toFixed(1)}°C, HR=${payload.bpm} bpm, SpO2=${payload.spo2}%`);
    ws.send(JSON.stringify(payload));
  }, 2000);
});

ws.on('error', (err) => {
  console.error('❌ Error WebSocket:', err.message);
});

ws.on('close', () => {
  console.log('🔌 Koneksi terputus dari Backend.');
});
