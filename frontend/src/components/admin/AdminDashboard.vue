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
            <router-link to="/dashboard" :class="{ 'bg-secondary text-primary': $route.path === '/dashboard' || $route.path === '/dashboard/' }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <font-awesome-icon :icon="['fas', 'tachometer-alt']" class="mr-3" />
              <span>Dashboard</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/employees" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/employees') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <font-awesome-icon :icon="['fas', 'users']" class="mr-3" />
              <span>Karyawan</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/attendance" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/attendance') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <font-awesome-icon :icon="['fas', 'calendar-check']" class="mr-3" />
              <span>Absensi</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/leave-requests" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/leave-requests') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <font-awesome-icon :icon="['fas', 'calendar-alt']" class="mr-3" />
              <span>Pengajuan Cuti & Izin</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/broadcast" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/broadcast') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <font-awesome-icon :icon="['fas', 'bullhorn']" class="mr-3" />
              <span>Broadcast</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/locations" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/locations') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <font-awesome-icon :icon="['fas', 'map-marker-alt']" class="mr-3" />
              <span>Lokasi Absensi</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/settings" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/settings') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <font-awesome-icon :icon="['fas', 'cog']" class="mr-3" />
              <span>Pengaturan</span>
            </router-link>
          </li>
        </ul>
      </nav>
      <div class="p-4 border-t border-bg-muted">
        <BaseButton @click="handleLogout" class="w-full btn-danger">
          <i class="fas fa-sign-out-alt"></i> Logout
        </BaseButton>
      </div>
    </aside>

    <!-- Overlay for mobile -->
    <div v-if="isSidebarOpen" class="fixed inset-0 bg-black opacity-50 z-40 md:hidden" @click="isSidebarOpen = false"></div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Header -->
      <header class="flex justify-between items-center p-4 bg-bg-muted text-text-base shadow-md">
        <BaseButton @click="isSidebarOpen = !isSidebarOpen" class="md:hidden text-text-base focus:outline-none">
          <font-awesome-icon :icon="['fas', 'bars']" class="h-6 w-6" />
        </BaseButton>
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

      <!-- Timezone Warning Banner -->
      <div v-if="showTimezoneWarning" class="bg-yellow-300 text-yellow-800 text-center p-2">
        <span>Zona waktu perusahaan Anda belum diatur atau masih menggunakan default.</span>
        <router-link to="/dashboard/settings" class="underline font-bold ml-2">Atur Sekarang</router-link>
      </div>

      <!-- Page Content -->
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-bg-base p-6">
        <router-view v-slot="{ Component }">
          <component :is="Component" ref="routerViewComponent" />
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';
import BaseButton from '../ui/BaseButton.vue';

const router = useRouter();
const isSidebarOpen = ref(false);
const toast = useToast();
const authStore = useAuthStore();
const routerViewComponent = ref(null);

const handleLogout = () => {
  if (routerViewComponent.value && typeof routerViewComponent.value.disconnectWebSocket === 'function') {
    routerViewComponent.value.disconnectWebSocket();
  }
  authStore.clearAuth();
  axios.defaults.headers.common['Authorization'] = '';
  router.push('/');
};

const loadCompanyDetails = () => {
  if (authStore.token && authStore.companyId) {
    authStore.fetchCompanyDetails();
  }
};

onMounted(() => {
  loadCompanyDetails();
});

watch(() => authStore.companyId, (newCompanyId) => {
  if (newCompanyId) {
    loadCompanyDetails();
  }
}, { immediate: true });

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

const showTimezoneWarning = computed(() => {
  return !authStore.companyTimezone;
});
</script>
<style scoped>
/* Tailwind handles styling */
</style>


