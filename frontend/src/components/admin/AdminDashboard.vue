<template>
  <div class="flex h-screen bg-bg-base">
    <!-- Sidebar -->
    <aside :class="[
      'w-64 bg-primary text-white flex-shrink-0 flex flex-col',
      isSidebarOpen ? 'translate-x-0 ease-out' : '-translate-x-full ease-in',
      'md:translate-x-0 md:static',
      'fixed inset-y-0 left-0 z-50 transform transition-transform duration-300',
    ]">
      <div class="p-4 text-2xl font-bold border-b border-bg-muted">
        Admin Panel
      </div>
      <nav class="flex-grow p-4">
        <ul>
        
          <li class="mb-2">
            <router-link to="/dashboard" class="block py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              Dashboard
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/employees" class="block py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              Karyawan
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/attendance" class="block py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              Absensi
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/settings" class="block py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              Pengaturan
            </router-link>
          </li>
        </ul>
      </nav>
      <div class="p-4 border-t border-bg-muted">
        <button @click="handleLogout" class="w-full btn btn-danger">
          Logout
        </button>
      </div>
    </aside>

    <!-- Overlay for mobile -->
    <div v-if="isSidebarOpen" class="fixed inset-0 bg-black opacity-50 z-40 md:hidden" @click="isSidebarOpen = false"></div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Header -->
      <header class="flex justify-between items-center p-4 bg-bg-muted text-text-base shadow-md">
        <button @click="isSidebarOpen = !isSidebarOpen" class="md:hidden text-text-base focus:outline-none">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <h1 class="text-xl font-semibold">Selamat Datang, Admin <span v-if="authStore.companyName"> {{ authStore.companyName }}</span>!</h1>
        <div>
          <span class="text-text-muted">{{ authStore.adminEmail }}</span>
        </div>
      </header>

      <!-- Trial Banner -->
      <div v-if="isTrial" class="bg-yellow-400 text-yellow-900 text-center p-2">
        <span>Anda dalam masa coba gratis. Sisa waktu Anda: {{ trialDaysRemaining }} hari.</span>
        <router-link to="/dashboard/subscribe" class="underline font-bold ml-2">Berlangganan Sekarang</router-link>
      </div>

      <!-- Subscription Expiring Soon Banner -->
      <div v-if="isExpiringSoon" class="bg-orange-400 text-orange-900 text-center p-2">
        <span>Langganan Anda akan berakhir dalam {{ subscriptionDaysRemaining }} hari.</span>
        <router-link to="/dashboard/subscribe" class="underline font-bold ml-2">Perpanjang Sekarang</router-link>
      </div>

      <!-- Subscription Expired Banner -->
      <div v-if="isExpired" class="bg-red-400 text-red-900 text-center p-2">
        <span>Langganan Anda telah kedaluwarsa. Beberapa fitur mungkin tidak dapat diakses.</span>
        <router-link to="/dashboard/subscribe" class="underline font-bold ml-2">Perpanjang Sekarang</router-link>
      </div>

      <!-- Page Content -->
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-bg-base">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'AdminDashboard',
  setup() {
    const router = useRouter();
    const isSidebarOpen = ref(false);
    const toast = useToast();
    const authStore = useAuthStore();

    const handleLogout = () => {
      authStore.clearAuth();
      axios.defaults.headers.common['Authorization'] = '';
      router.push('/');
    };

    onMounted(() => {
      authStore.fetchCompanyDetails();
    });

    const isTrial = computed(() => authStore.subscriptionStatus === 'trial');

    const trialDaysRemaining = computed(() => {
      if (!authStore.trialEndDate) return 0;
      const endDate = new Date(authStore.trialEndDate);
      const now = new Date();
      const diffTime = endDate.getTime() - now.getTime();
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      return diffDays > 0 ? diffDays : 0;
    });

    const subscriptionDaysRemaining = computed(() => {
      if (!authStore.subscriptionEndDate) return 0;
      const endDate = new Date(authStore.subscriptionEndDate);
      const now = new Date();
      const diffTime = endDate.getTime() - now.getTime();
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      return diffDays > 0 ? diffDays : 0;
    });

    const isExpiringSoon = computed(() => {
      return authStore.subscriptionStatus === 'active' && subscriptionDaysRemaining.value > 0 && subscriptionDaysRemaining.value <= 7;
    });

    const isExpired = computed(() => {
      return authStore.subscriptionStatus === 'expired' || authStore.subscriptionStatus === 'expired_trial';
    });

    return {
      isSidebarOpen,
      handleLogout,
      authStore,
      isTrial,
      trialDaysRemaining,
      subscriptionDaysRemaining,
      isExpiringSoon,
      isExpired,
    };
  },
};
</script>
<style scoped>
/* Tailwind handles styling */
</style>

