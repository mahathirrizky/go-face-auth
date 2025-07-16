<template>
  <div class="bg-bg-muted p-6 rounded-lg shadow-md">
    <h3 class="text-xl font-semibold text-text-base mb-4">Pengaturan Umum Perusahaan</h3>
    <div class="mb-4">
      <label for="companyName" class="block text-text-muted text-sm font-bold mb-2">Nama Perusahaan:</label>
      <InputText
        id="companyName"
        v-model="settings.companyName"
        class="w-full"
      />
    </div>
    <div class="mb-4">
      <label for="companyAddress" class="block text-text-muted text-sm font-bold mb-2">Alamat Perusahaan:</label>
      <InputText
        id="companyAddress"
        v-model="settings.companyAddress"
        class="w-full"
      />
    </div>
    <div class="mb-4">
      <label for="timezone" class="block text-text-muted text-sm font-bold mb-2">Zona Waktu:</label>
      <Select
        id="timezone"
        v-model="settings.timezone"
        :options="timezones"
        optionLabel="label"
        optionValue="value"
        placeholder="Pilih Zona Waktu"
        class="w-full"
      />
    </div>
    <BaseButton @click="saveSettings" class="mt-4">
      <i class="pi pi-save"></i> Simpan Pengaturan
    </BaseButton>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import axios from 'axios';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
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
      toast.add({ severity: 'success', summary: 'Success', detail: 'Pengaturan umum berhasil disimpan!', life: 3000 });
      authStore.companyName = settings.value.companyName;
      authStore.companyAddress = settings.value.companyAddress;
      authStore.companyTimezone = settings.value.timezone;
      authStore.hasConfiguredTimezone = true;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Gagal menyimpan pengaturan umum.', life: 3000 });
    }
  } catch (error) {
    console.error('Error saving general settings:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Terjadi kesalahan saat menyimpan pengaturan.', life: 3000 });
  }
};
</script>
