<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <h2 class="text-2xl font-bold text-text-base mb-6">Kirim Pesan Broadcast ke Karyawan</h2>

    <Card class="bg-bg-muted shadow-md">
      <template #title>
        <h3 class="text-xl font-semibold text-text-base">Pesan Broadcast Baru</h3>
      </template>
      <template #content>
        <div class="p-fluid">
          <div class="field mb-4">
            <label for="broadcastMessage">Pesan</label>
            <Textarea
              id="broadcastMessage"
              v-model="broadcastMessage"
              placeholder="Tulis pesan broadcast Anda di sini..."
              rows="4"
              :autoResize="true"
              fluid
            />
          </div>
          <div class="field mb-4">
            <label for="expireDate">Berlaku Hingga Tanggal (Opsional)</label>
            <DatePicker
              id="expireDate"
              v-model="expireDate"
              dateFormat="yy-mm-dd"
              :minDate="new Date()"
              showIcon
              placeholder="Pilih tanggal kedaluwarsa"
          fluid  
              />
          </div>
          <Button @click="sendBroadcastMessage" :loading="isSending" :label="isSending ? 'Mengirim...' : 'Kirim Broadcast'" icon="pi pi-send" />
        </div>
      </template>
    </Card>

    <Card class="bg-bg-muted shadow-md mt-8">
        <template #title>
            <h3 class="text-xl font-semibold text-text-base">Riwayat Pesan Broadcast</h3>
        </template>
        <template #content>
            <div v-if="isLoadingHistory" class="flex items-center justify-center py-4">
                <ProgressSpinner />
                <span class="ml-2">Memuat riwayat...</span>
            </div>
            <div v-else-if="adminBroadcastStore.broadcastMessages.length > 0">
                <div v-for="(msg, index) in adminBroadcastStore.broadcastMessages" :key="msg.id || index" class="bg-bg-base p-4 rounded-lg shadow-sm mb-3 border border-surface-border">
                <p class="text-text-base">{{ msg.message }}</p>
                <p class="text-text-muted text-sm mt-2">Dikirim: {{ isValidDate(msg.created_at) ? new Date(msg.created_at).toLocaleString('id-ID') : 'N/A' }}</p>
                <p class="text-text-muted text-sm">Berlaku hingga: {{ isValidDate(msg.expire_date) ? new Date(msg.expire_date).toLocaleDateString('id-ID') : 'Tidak ada tanggal kedaluwarsa' }}</p>
                </div>
            </div>
            <div v-else>
                <p class="text-text-muted">Belum ada pesan broadcast yang dikirim.</p>
            </div>
        </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAdminBroadcastStore } from '../../stores/adminBroadcast';
import Button from 'primevue/button';
import Textarea from 'primevue/textarea';
import DatePicker from 'primevue/datepicker';
import Card from 'primevue/card';
import ProgressSpinner from 'primevue/progressspinner';
import Toast from 'primevue/toast';

const toast = useToast();
const broadcastMessage = ref('');
const expireDate = ref(null);
const isSending = ref(false);
const isLoadingHistory = ref(false);

const adminBroadcastStore = useAdminBroadcastStore();

const isValidDate = (dateString) => {
  if (!dateString) return false;
  const d = new Date(dateString);
  return d instanceof Date && !isNaN(d);
};

onMounted(() => {
  isLoadingHistory.value = true;
  adminBroadcastStore.fetchBroadcasts().finally(() => {
    isLoadingHistory.value = false;
  });
});

const sendBroadcastMessage = async () => {
  if (!broadcastMessage.value.trim()) {
    toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Pesan broadcast tidak boleh kosong.', life: 3000 });
    return;
  }

  isSending.value = true;
  try {
    const payload = {
      message: broadcastMessage.value,
      expire_date: expireDate.value ? expireDate.value.toISOString().slice(0, 10) : null,
    };
    await axios.post('/api/broadcasts', payload);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Pesan broadcast berhasil dikirim!', life: 3000 });

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
.field > label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
}
</style>