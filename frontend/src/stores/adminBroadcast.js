import { defineStore } from 'pinia';
import axios from 'axios';

export const useAdminBroadcastStore = defineStore('adminBroadcast', {
  state: () => ({
    broadcastMessages: [], // This will now hold messages fetched from the backend
  }),
  actions: {
    async fetchBroadcasts() {
      try {
        const response = await axios.get('/api/broadcasts');
        this.broadcastMessages = response.data.data; // Assuming data is in response.data.data
      } catch (error) {
        console.error('Error fetching broadcast messages:', error);
        // Handle error, e.g., show a toast notification
      }
    },
  },
});