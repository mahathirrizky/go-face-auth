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
      <BaseInput
        id="expireDate"
        label="Berlaku Hingga Tanggal:"
        v-model="expireDate"
        type="date"
        :min="todayDate"
      />
      <BaseButton @click="sendBroadcastMessage" class="btn-primary">
        Kirim Broadcast
      </BaseButton>

      <h3 class="text-xl font-semibold text-text-base mt-8 mb-4">Riwayat Pesan Broadcast</h3>
      <div v-if="adminBroadcastStore.broadcastMessages.length > 0">
        <div v-for="(msg, index) in adminBroadcastStore.broadcastMessages" :key="index" class="bg-bg-base p-4 rounded-lg shadow-sm mb-3">
          <p class="text-text-base">{{ msg.message }}</p>
          <p class="text-text-muted text-sm">Dikirim: {{ isValidDate(msg.timestamp) ? new Date(msg.timestamp).toLocaleString() : 'N/A' }}</p>
          <p class="text-text-muted text-sm">Berlaku hingga: {{ isValidDate(msg.expire_date) ? new Date(msg.expire_date).toLocaleDateString() : 'N/A' }}</p>
        </div>
      </div>
      <div v-else>
        <p class="text-text-muted">Belum ada pesan broadcast yang dikirim.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAdminBroadcastStore } from '../../stores/adminBroadcast';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';

const toast = useToast();
const broadcastMessage = ref('');
const expireDate = ref('');

const todayDate = computed(() => {
  const now = new Date();
  const year = now.getFullYear();
  const month = (now.getMonth() + 1).toString().padStart(2, '0');
  const day = now.getDate().toString().padStart(2, '0');
  return `${year}-${month}-${day}`;
});
const adminBroadcastStore = useAdminBroadcastStore();

const isValidDate = (dateString) => {
  const d = new Date(dateString);
  return d instanceof Date && !isNaN(d);
};

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
      expire_date: expireDate.value,
    };
    await axios.post('/api/broadcast', payload);
    toast.success('Pesan broadcast berhasil dikirim!');

    adminBroadcastStore.addBroadcastMessage({
      message: broadcastMessage.value,
      expire_date: expireDate.value,
      timestamp: new Date().toISOString(),
    });

    broadcastMessage.value = '';
    expireDate.value = '';
  } catch (error) {
    console.error('Error sending broadcast message:', error);
    toast.error(error.response?.data?.meta?.message || 'Gagal mengirim pesan broadcast.');
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
