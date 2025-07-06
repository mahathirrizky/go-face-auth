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
    });

    const saveSettings = async () => {
      // This function will be implemented when there's a backend endpoint to update company details
      toast.info('Fungsi simpan pengaturan umum belum diimplementasikan.');
      console.log('Saving general settings:', settings.value);
    };

    return {
      settings,
      saveSettings,
    };
  },
};
</script>
