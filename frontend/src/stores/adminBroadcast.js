import { defineStore } from 'pinia';

interface BroadcastMessage {
  message: string;
  expire_date: string; // YYYY-MM-DD format
  timestamp: string; // ISO string of when it was sent
}

interface AdminBroadcastState {
  broadcastMessages: BroadcastMessage[];
}

export const useAdminBroadcastStore = defineStore('adminBroadcast', {
  state: (): AdminBroadcastState => ({
    broadcastMessages: [],
  }),
  actions: {
    addBroadcastMessage(message: BroadcastMessage) {
      // Add new message to the beginning of the array
      this.broadcastMessages.unshift(message);
      // Optional: Limit the number of stored messages to prevent excessive storage
      // this.broadcastMessages = this.broadcastMessages.slice(0, 50); 
    },
    // You might add an action to clean up expired messages if needed, 
    // but for admin side, it might be useful to see all sent messages.
  },
  persist: true, // Enable persistence for this store
});
