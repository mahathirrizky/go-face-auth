<template>
  <div class="bg-bg-muted p-6 rounded-lg shadow-md">
    <h3 class="text-xl font-semibold text-text-base mb-4">Pengaturan Umum Perusahaan</h3>
    <BaseInput
      id="companyName"
      label="Nama Perusahaan:"
      v-model="settings.companyName"
    />
    <BaseInput
      id="companyAddress"
      label="Alamat Perusahaan:"
      v-model="settings.companyAddress"
    />
    <div class="mb-4">
      <label for="timezone" class="block text-text-muted text-sm font-bold mb-2">Zona Waktu:</label>
      <select
        id="timezone"
        v-model="settings.timezone"
        class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
      >
        <option v-for="tz in timezones" :key="tz.value" :value="tz.value">{{ tz.label }}</option>
      </select>
    </div>
    <BaseButton @click="saveSettings" class="mt-4">
      Simpan Pengaturan
    </BaseButton>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';
import axios from 'axios';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';

const authStore = useAuthStore();
const toast = useToast();

const settings = ref({
  companyName: authStore.companyName || '',
  companyAddress: authStore.companyAddress || '',
  timezone: authStore.companyTimezone || 'Asia/Jakarta',
});

const timezones = ref([
  { value: 'Asia/Jakarta', label: 'Asia/Jakarta (WIB)' },
  { value: 'Asia/Makassar', label: 'Asia/Makassar (WITA)' },
  { value: 'Asia/Jayapura', label: 'Asia/Jayapura (WIT)' },
  { value: 'UTC', label: 'UTC' },
]);

const saveSettings = async () => {
  try {
    const payload = {
      name: settings.value.companyName,
      address: settings.value.companyAddress,
      timezone: settings.value.timezone,
    };
    const response = await axios.put('/api/company-details', payload, {
      headers: { Authorization: `Bearer ${authStore.token}` },
    });
    if (response.data && response.data.status === 'success') {
      toast.success('Pengaturan umum berhasil disimpan!');
      authStore.companyName = settings.value.companyName;
      authStore.companyAddress = settings.value.companyAddress;
      authStore.companyTimezone = settings.value.timezone;
      authStore.hasConfiguredTimezone = true;
    } else {
      toast.error(response.data.message || 'Gagal menyimpan pengaturan umum.');
    }
  } catch (error) {
    console.error('Error saving general settings:', error);
    toast.error(error.response?.data?.message || 'Terjadi kesalahan saat menyimpan pengaturan.');
  }
};
</script>
