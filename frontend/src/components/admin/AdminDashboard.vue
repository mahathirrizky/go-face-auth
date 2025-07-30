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
              <i class="pi pi-th-large mr-3 text-blue-300"></i>
              <span>Dashboard</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/employees" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/employees') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-users mr-3 text-green-300"></i>
              <span>Karyawan</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/attendance" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/attendance') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-calendar-clock mr-3 text-yellow-300"></i>
              <span>Absensi</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/leave-requests" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/leave-requests') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-calendar mr-3 text-red-300"></i>
              <span>Pengajuan Cuti & Izin</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/broadcast" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/broadcast') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-megaphone mr-3 text-purple-300"></i>
              <span>Broadcast</span>
            </router-link>
          </li>
          
          <li class="mb-2">
            <router-link to="/dashboard/divisions" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/divisions') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-sitemap mr-3 text-orange-300"></i>
              <span>Divisi</span>
            </router-link>
          </li>
          
          
          <li class="mb-2">
            <router-link to="/dashboard/subscribe" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/subscribe') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-credit-card mr-3 text-teal-300"></i>
              <span>Langganan</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/billing-history" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/billing-history') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-history mr-3 text-pink-300"></i>
              <span>Riwayat Tagihan</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/dashboard/settings" :class="{ 'bg-secondary text-primary': $route.path.startsWith('/dashboard/settings') }" class="flex items-center py-2 px-4 rounded hover:bg-secondary hover:text-primary transition-colors duration-200" @click="isSidebarOpen = false">
              <i class="pi pi-cog mr-3 text-gray-300"></i>
              <span>Pengaturan</span>
            </router-link>
          </li>
        </ul>
      </nav>
      <div class="p-4 border-t border-bg-muted">
        <Button @click="handleLogout" class="w-full p-button-danger" icon="pi pi-sign-out" label="Logout" />
      </div>
    </aside>

    <!-- Overlay for mobile -->
    <div v-if="isSidebarOpen" class="fixed inset-0 bg-black opacity-50 z-40 md:hidden" @click="isSidebarOpen = false"></div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Header -->
      <header class="flex justify-between items-center p-4 bg-bg-muted text-text-base shadow-md">
        <Button @click="isSidebarOpen = !isSidebarOpen" class="md:hidden p-button-text text-text-base" icon="pi pi-bars" />
        <h1 class="text-xl font-semibold">Selamat Datang, Admin <span v-if="authStore.companyName"> {{ authStore.companyName }}</span>!</h1>
        <div>
          <span class="text-text-muted">{{ authStore.adminEmail }}</span>
        </div>
      </header>

      <!-- Banners -->
      <Message v-if="isTrial" severity="warn" :closable="false">Anda dalam masa coba gratis. Sisa waktu Anda: {{ trialDaysRemaining }} hari. <router-link to="/dashboard/subscribe" class="underline font-bold ml-2">Berlangganan Sekarang</router-link></Message>
      <Message v-if="isExpiringSoon" severity="warn" :closable="false">Langganan Anda akan berakhir dalam {{ subscriptionDaysRemaining }} hari. <router-link to="/dashboard/subscribe" class="underline font-bold ml-2">Perpanjang Sekarang</router-link></Message>
      <Message v-if="isExpired" severity="error" :closable="false">Langganan Anda telah kedaluwarsa. Beberapa fitur mungkin tidak dapat diakses. <router-link to="/dashboard/subscribe" class="underline font-bold ml-2">Perpanjang Sekarang</router-link></Message>
      <Message v-if="showTimezoneWarning" severity="info" :closable="false">Zona waktu perusahaan Anda belum diatur atau masih menggunakan default. <router-link to="/dashboard/settings" class="underline font-bold ml-2">Atur Sekarang</router-link></Message>

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
import { useAuthStore } from '../../stores/auth';
import Button from 'primevue/button';
import Message from 'primevue/message';

const router = useRouter();
const isSidebarOpen = ref(false);

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