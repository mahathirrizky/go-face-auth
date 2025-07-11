<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Kirim Pesan Broadcast ke Karyawan</h2>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pesan Broadcast</h3>
      <div class="mb-4">
        <textarea
          v-model="broadcastMessage"
          placeholder="Tulis pesan broadcast Anda di sini..."
          rows="4"
          class="form-input w-full p-2 border border-gray-300 rounded-md bg-bg-base text-text-base focus:outline-none focus:ring-secondary focus:border-secondary"
        ></textarea>
      </div>
      <div class="mb-4">
        <label for="expireDate" class="block text-text-muted text-sm font-bold mb-2">Berlaku Hingga Tanggal:</label>
        <input
          type="date"
          id="expireDate"
          v-model="expireDate"
          class="form-input w-full p-2 border border-gray-300 rounded-md bg-bg-base text-text-base focus:outline-none focus:ring-secondary focus:border-secondary"
        />
      </div>
      <button @click="sendBroadcastMessage" class="btn btn-primary">
        Kirim Broadcast
      </button>

      <h3 class="text-xl font-semibold text-text-base mt-8 mb-4">Riwayat Pesan Broadcast</h3>
      <div v-if="adminBroadcastStore.broadcastMessages.length > 0">
        <div v-for="(msg, index) in adminBroadcastStore.broadcastMessages" :key="index" class="bg-bg-base p-4 rounded-lg shadow-sm mb-3">
          <p class="text-text-base">{{ msg.message }}</p>
          <p class="text-text-muted text-sm">Dikirim: {{ new Date(msg.timestamp).toLocaleString() }}</p>
          <p class="text-text-muted text-sm">Berlaku hingga: {{ new Date(msg.expire_date).toLocaleDateString() }}</p>
        </div>
      </div>
      <div v-else>
        <p class="text-text-muted">Belum ada pesan broadcast yang dikirim.</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAdminBroadcastStore } from '../../stores/adminBroadcast'; // Import the new store

export default {
  name: 'BroadcastMessagePage',
  setup() {
    const toast = useToast();
    const broadcastMessage = ref('');
    const expireDate = ref(''); // New ref for expire date
    const adminBroadcastStore = useAdminBroadcastStore(); // Initialize the new store

    const sendBroadcastMessage = async () => {
      if (!broadcastMessage.value.trim()) {
        toast.error('Pesan broadcast tidak boleh kosong.');
        return;
      }
      if (!expireDate.value) {
        toast.error('Tanggal berlaku hingga harus diisi.');
        return;
      }

      try {
        const payload = {
          message: broadcastMessage.value,
          expire_date: expireDate.value, // Include expire date in payload
        };
        await axios.post('/api/broadcast', payload);
        toast.success('Pesan broadcast berhasil dikirim!');

        // Save to Pinia store
        adminBroadcastStore.addBroadcastMessage({
          message: broadcastMessage.value,
          expire_date: expireDate.value,
          timestamp: new Date().toISOString(),
        });

        broadcastMessage.value = ''; // Clear the message input
        expireDate.value = ''; // Clear expire date input
      } catch (error) {
        console.error('Error sending broadcast message:', error);
        toast.error(error.response?.data?.meta?.message || 'Gagal mengirim pesan broadcast.');
      }
    };

    return {
      broadcastMessage,
      expireDate,
      sendBroadcastMessage,
      adminBroadcastStore, // Expose the store
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>