<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Kirim Pesan Broadcast ke Karyawan</h2>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pesan Broadcast</h3>
      <div class="mb-4">
        <Textarea
          v-model="broadcastMessage"
          placeholder="Tulis pesan broadcast Anda di sini..."
          rows="4"
          class="w-full"
        />
      </div>
      <div class="mb-4">
        <label for="expireDate" class="block text-text-muted text-sm font-bold mb-2">Berlaku Hingga Tanggal:</label>
        <DatePicker
          id="expireDate"
          v-model="expireDate"
          dateFormat="yy-mm-dd"
          :minDate="new Date()"
          showIcon
          class="w-full"
        />
      </div>
      <BaseButton @click="sendBroadcastMessage" class="btn-primary" :disabled="isSending">
        <i v-if="!isSending" class="pi pi-send"></i>
        <i v-else class="pi pi-spin pi-spinner"></i>
        {{ isSending ? 'Mengirim...' : 'Kirim Broadcast' }}
      </BaseButton>

      <h3 class="text-xl font-semibold text-text-base mt-8 mb-4">Riwayat Pesan Broadcast</h3>
      <div v-if="isLoadingHistory" class="flex items-center justify-center py-4">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        <span class="ml-2">Memuat riwayat...</span>
      </div>
      <div v-else-if="adminBroadcastStore.broadcastMessages.length > 0">
        <div v-for="(msg, index) in adminBroadcastStore.broadcastMessages" :key="msg.id || index" class="bg-bg-base p-4 rounded-lg shadow-sm mb-3">
          <p class="text-text-base">{{ msg.message }}</p>
          <p class="text-text-muted text-sm">Dikirim: {{ isValidDate(msg.created_at) ? new Date(msg.created_at).toLocaleString() : 'N/A' }}</p>
          <p class="text-text-muted text-sm">Berlaku hingga: {{ isValidDate(msg.expire_date) ? new Date(msg.expire_date).toLocaleDateString() : 'Tidak ada tanggal kedaluwarsa' }}</p>
        </div>
      </div>
      <div v-else>
        <p class="text-text-muted">Belum ada pesan broadcast yang dikirim.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAdminBroadcastStore } from '../../stores/adminBroadcast';
import BaseButton from '../ui/BaseButton.vue';
import Textarea from 'primevue/textarea';
import DatePicker from 'primevue/datepicker';

const toast = useToast();
const broadcastMessage = ref('');
const expireDate = ref(null);
const isSending = ref(false);
const isLoadingHistory = ref(false);

const todayDate = computed(() => {
  const now = new Date();
  const year = now.getFullYear();
  const month = (now.getMonth() + 1).toString().padStart(2, '0');
  const day = now.getDate().toString().padStart(2, '0');
  return `${year}-${month}-${day}`;
});
const adminBroadcastStore = useAdminBroadcastStore();

const isValidDate = (dateString) => {
  if (!dateString) return false; // Handle null or empty date strings
  const d = new Date(dateString);
  return d instanceof Date && !isNaN(d);
};

// Fetch messages when the component is mounted
onMounted(() => {
  isLoadingHistory.value = true;
  adminBroadcastStore.fetchBroadcasts().finally(() => {
    isLoadingHistory.value = false;
  });
});

const sendBroadcastMessage = async () => {
  if (!broadcastMessage.value.trim()) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Pesan broadcast tidak boleh kosong.', life: 3000 });
    return;
  }
  // expireDate can be empty, meaning no expiration

  isSending.value = true;
  try {
    const payload = {
      message: broadcastMessage.value,
      expire_date: expireDate.value ? expireDate.value.toISOString().slice(0, 10) : null, // Send null if empty
    };
    await axios.post('/api/broadcasts', payload); // Changed endpoint to /api/broadcasts
    toast.add({ severity: 'success', summary: 'Success', detail: 'Pesan broadcast berhasil dikirim!', life: 3000 });

    // Refresh the list of messages from the backend
    isLoadingHistory.value = true;
    await adminBroadcastStore.fetchBroadcasts().finally(() => {
      isLoadingHistory.value = false;
    });

    broadcastMessage.value = '';
    expireDate.value = null;
  } catch (error) {
    console.error('Error sending broadcast message:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.meta?.message || 'Gagal mengirim pesan broadcast.', life: 3000 });
  } finally {
    isSending.value = false;
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
