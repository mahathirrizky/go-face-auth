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

  // Use window.base_url (which is https://api.4commander.my.id) for the WebSocket host
  const apiBaseUrl = window.base_url; 
  const wsProtocol = apiBaseUrl.startsWith('https:') ? 'wss:' : 'ws:';
  const wsHost = apiBaseUrl.replace(/^(http|https):\/\//, ''); // Remove protocol from base URL

  const wsUrl = `${wsProtocol}//${wsHost}${currentWsPath}?token=${authStore.token}`;
  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    isConnected.value = true;
    if (reconnectInterval) {
      clearInterval(reconnectInterval);
      reconnectInterval = null;
    }
  };

  ws.onmessage = (event) => {
    // Handle incoming messages (e.g., dashboard updates)
  };

  ws.onclose = (event) => {
    isConnected.value = false;
    // Attempt to reconnect after a delay
    if (!reconnectInterval) {
      reconnectInterval = setInterval(() => {
        connectWebSocket(currentWsPath); // Pass the stored path during reconnection
      }, 5000); // Try to reconnect every 5 seconds
    };

  ws.onerror = (error) => {
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
