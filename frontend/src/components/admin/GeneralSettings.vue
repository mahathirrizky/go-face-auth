<template>
  <div class="bg-bg-muted p-6 rounded-lg shadow-md">
    <h3 class="text-xl font-semibold text-text-base mb-4">Pengaturan Umum Perusahaan</h3>
    <div class="mb-4">
      <label for="companyName" class="block text-text-muted text-sm font-bold mb-2">Nama Perusahaan:</label>
      <input
        type="text"
        id="companyName"
        v-model="settings.companyName"
        class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
      />
    </div>
    <div class="mb-4">
      <label for="companyAddress" class="block text-text-muted text-sm font-bold mb-2">Alamat Perusahaan:</label>
      <input
        type="text"
        id="companyAddress"
        v-model="settings.companyAddress"
        class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
      />
    </div>
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
    <button @click="saveSettings" class="btn btn-secondary mt-4">
      Simpan Pengaturan
    </button>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'GeneralSettings',
  setup() {
    const authStore = useAuthStore();
    const toast = useToast();

    const settings = ref({
      companyName: authStore.companyName || '',
      companyAddress: authStore.companyAddress || '',
      timezone: authStore.companyTimezone || 'Asia/Jakarta', // Initialize with current timezone or default
    });

    const timezones = ref([
      { value: 'Asia/Jakarta', label: 'Asia/Jakarta (WIB)' },
      { value: 'Asia/Makassar', label: 'Asia/Makassar (WITA)' },
      { value: 'Asia/Jayapura', label: 'Asia/Jayapura (WIT)' },
      { value: 'UTC', label: 'UTC' },
      // Add more timezones as needed
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
          // Update auth store with new values
          authStore.companyName = settings.value.companyName;
          authStore.companyAddress = settings.value.companyAddress;
          authStore.companyTimezone = settings.value.timezone;
        } else {
          toast.error(response.data.message || 'Gagal menyimpan pengaturan umum.');
        }
      } catch (error) {
        console.error('Error saving general settings:', error);
        toast.error(error.response?.data?.message || 'Terjadi kesalahan saat menyimpan pengaturan.');
      }
    };

    return {
      settings,
      timezones,
      saveSettings,
    };
  },
};
</script>
