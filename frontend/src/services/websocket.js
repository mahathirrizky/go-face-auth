import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';

const isConnected = ref(false);
let ws = null;
let reconnectInterval = null;
let currentWsPath = null; // Store the path for reconnection

const connectWebSocket = (wsPath) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    return; // Already connected
  }

  currentWsPath = wsPath; // Store the path

  const authStore = useAuthStore();
  if (!authStore.token) {
    console.warn('WebSocket: No token available, skipping connection.');
    isConnected.value = false;
    return;
  }

  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${currentWsPath}?token=${authStore.token}`;
  // console.log('WebSocket URL components:', { host: window.location.host, path: currentWsPath, fullUrl: wsUrl }); // Remove after debugging
  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    // console.log('WebSocket connected.'); // Remove after debugging
    isConnected.value = true;
    if (reconnectInterval) {
      clearInterval(reconnectInterval);
      reconnectInterval = null;
    }
  };

  ws.onmessage = (event) => {
    // Handle incoming messages (e.g., dashboard updates)
    // console.log('WebSocket message received:', event.data); // Remove after debugging
  };

  ws.onclose = (event) => {
    // console.warn(`WebSocket disconnected: ${event.code} ${event.reason}`); // Remove after debugging
    isConnected.value = false;
    // Attempt to reconnect after a delay
    if (!reconnectInterval) {
      // console.log('Attempting to reconnect WebSocket...'); // Remove after debugging
      reconnectInterval = setInterval(() => {
        connectWebSocket(currentWsPath); // Pass the stored path during reconnection
      }, 5000); // Try to reconnect every 5 seconds
    }
  };

  ws.onerror = (error) => {
    // console.error('WebSocket error:', error); // Remove after debugging
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
