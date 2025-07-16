<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md text-center">
      <h2 class="text-2xl font-bold text-center mb-4 text-text-base">Konfirmasi Email</h2>
      <p v-if="loading" class="text-text-base">Memverifikasi email Anda...</p>
      <p v-else-if="success" class="text-green-700 font-semibold mb-4">{{ message }}</p>
      <p v-else class="text-red-700 font-semibold mb-4">{{ message }}</p>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { getBaseDomain } from '../../utils/subdomain'; // Import getBaseDomain

export default {
  name: 'ConfirmEmail',
  setup() {
    const route = useRoute();
    const router = useRouter();
    const toast = useToast();
    const loading = ref(true);
    const success = ref(false);
    const message = ref('');

    onMounted(async () => {
      const token = route.query.token;
      if (!token) {
        message.value = 'Token konfirmasi tidak ditemukan.';
        success.value = false;
        loading.value = false;
        toast.add({ severity: 'error', summary: 'Error', detail: message.value, life: 3000 });
        return;
      }

      try {
        const response = await axios.get(`/api/confirm-email?token=${token}`);
        if (response.data && response.data.status === 'success') {
          message.value = response.data.message || 'Email Anda berhasil dikonfirmasi!';
          success.value = true;
          toast.add({ severity: 'success', summary: 'Success', detail: response.data.message, life: 3000 });
          // Auto-redirect to admin login after 3 seconds
          setTimeout(() => {
            const baseDomain = getBaseDomain();
            const adminLoginURL = `${window.location.protocol}//admin.${baseDomain}${window.location.port ? ':' + window.location.port : ''}/`;
            window.location.href = adminLoginURL;
          }, 3000); // 3 second delay
        } else {
          message.value = response.data?.message || 'Gagal mengkonfirmasi email.';
          success.value = false;
          toast.add({ severity: 'error', summary: 'Error', detail: message.value, life: 3000 });
        }
      } catch (error) {
        console.error('Error confirming email:', error);
        message.value = error.response?.data?.message || 'Terjadi kesalahan saat mengkonfirmasi email.';
        success.value = false;
        toast.add({ severity: 'error', summary: 'Error', detail: message.value, life: 3000 });
      } finally {
        loading.value = false;
      }
    });

    return {
      loading,
      success,
      message,
    };
  },
};
</script>