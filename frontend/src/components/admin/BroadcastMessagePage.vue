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
      <button @click="sendBroadcastMessage" class="btn btn-primary">
        Kirim Broadcast
      </button>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';

export default {
  name: 'BroadcastMessagePage',
  setup() {
    const toast = useToast();
    const broadcastMessage = ref('');

    const sendBroadcastMessage = async () => {
      if (!broadcastMessage.value.trim()) {
        toast.error('Pesan broadcast tidak boleh kosong.');
        return;
      }
      try {
        await axios.post('/api/broadcast', { message: broadcastMessage.value });
        toast.success('Pesan broadcast berhasil dikirim!');
        broadcastMessage.value = ''; // Clear the message input
      } catch (error) {
        console.error('Error sending broadcast message:', error);
        toast.error(error.response?.data?.meta?.message || 'Gagal mengirim pesan broadcast.');
      }
    };

    return {
      broadcastMessage,
      sendBroadcastMessage,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
