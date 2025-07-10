import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';

const isConnected = ref(false);
let ws = null;
let reconnectInterval = null;

const connectWebSocket = (wsPath) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    return; // Already connected
  }

  const authStore = useAuthStore();
  if (!authStore.token) {
    console.warn('WebSocket: No token available, skipping connection.');
    isConnected.value = false;
    return;
  }

  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${wsPath}?token=${authStore.token}`;
  console.log('WebSocket URL components:', { host: window.location.host, path: wsPath, fullUrl: wsUrl });
  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log('WebSocket connected.');
    isConnected.value = true;
    if (reconnectInterval) {
      clearInterval(reconnectInterval);
      reconnectInterval = null;
    }
  };

  ws.onmessage = (event) => {
    // Handle incoming messages (e.g., dashboard updates)
    console.log('WebSocket message received:', event.data);
    // You might want to parse event.data and update Pinia store or other reactive data
  };

  ws.onclose = (event) => {
    console.warn(`WebSocket disconnected: ${event.code} ${event.reason}`);
    isConnected.value = false;
    // Attempt to reconnect after a delay
    if (!reconnectInterval) {
      console.log('Attempting to reconnect WebSocket...');
      reconnectInterval = setInterval(() => {
        connectWebSocket();
      }, 5000); // Try to reconnect every 5 seconds
    }
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    isConnected.value = false;
    ws.close(); // Close the connection to trigger onclose and reconnection logic
  };
};

const disconnectWebSocket = () => {
  if (ws) {
    ws.close();
    ws = null;
  }
  if (reconnectInterval) {
    clearInterval(reconnectInterval);
    reconnectInterval = null;
  }
  isConnected.value = false;
};

export { isConnected, connectWebSocket, disconnectWebSocket };
