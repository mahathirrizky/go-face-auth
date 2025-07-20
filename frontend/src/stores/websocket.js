import { defineStore } from 'pinia';
import { useAuthStore } from './auth'; // Import auth store to get token

export const useWebSocketStore = defineStore('websocket', {
  state: () => ({
    ws: null,
    isConnected: false,
    reconnectAttempts: 0,
    maxReconnectAttempts: 5,
    reconnectInterval: 3000, // 3 seconds
    messageHandlers: {}, // To store callbacks for different message types
    superAdminDashboardData: null, // New state to store superadmin dashboard data
  }),
  actions: {
    initWebSocket(wsUrl) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        console.log('WebSocket already connected.');
        return;
      }

      const authStore = useAuthStore();
      if (!authStore.token) {
        console.warn('No auth token found, cannot establish WebSocket connection.');
        this.isConnected = false;
        return;
      }

      // Append token to WebSocket URL
      const urlWithToken = `${wsUrl}?token=${authStore.token}`;
      this.ws = new WebSocket(urlWithToken);

      this.ws.onopen = () => {
        console.log('WebSocket connected.');
        this.isConnected = true;
        this.reconnectAttempts = 0; // Reset reconnect attempts on successful connection
      };

      this.ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log('WebSocket message received:', data);
        if (data.type === 'superadmin_dashboard_update') {
          this.superAdminDashboardData = data.payload; // Store the data
        }
        if (data.type && this.messageHandlers[data.type]) {
          this.messageHandlers[data.type](data.payload);
        } else if (this.messageHandlers['*']) { // Fallback for generic handler
          this.messageHandlers['*'](data);
        }
      };

      this.ws.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        this.isConnected = false;
        if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) { // 1000 is normal closure
          this.reconnectAttempts++;
          console.log(`Attempting to reconnect WebSocket... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
          setTimeout(() => this.initWebSocket(wsUrl), this.reconnectInterval);
        } else if (this.reconnectAttempts >= this.maxReconnectAttempts) {
          console.error('Max reconnect attempts reached. WebSocket will not attempt to reconnect.');
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.isConnected = false;
        // Error might precede close, so close handler will manage reconnects
      };
    },

    closeWebSocket() {
      if (this.ws) {
        console.log('Closing WebSocket connection.');
        this.ws.close(1000, 'Manual closure'); // 1000 is normal closure
        this.ws = null;
        this.isConnected = false;
      }
    },

    sendMessage(message) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(message));
      } else {
        console.warn('WebSocket not connected. Message not sent:', message);
      }
    },

    // Register a handler for a specific message type or a generic handler
    onMessage(type, handler) {
      this.messageHandlers[type] = handler;
    },

    // Unregister a handler
    offMessage(type) {
      delete this.messageHandlers[type];
    },
  },
  // Pinia persist plugin will not persist WebSocket instance
});
