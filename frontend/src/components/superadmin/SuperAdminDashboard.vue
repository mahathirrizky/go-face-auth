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
        SuperAdmin Panel
      </div>
      <nav class="flex-grow p-4">
        <ul>
          <li class="mb-2">
            <router-link to="/dashboard" class="flex items-center justify-start p-2 rounded-md hover:bg-secondary hover:text-primary transition-colors duration-200">
              <font-awesome-icon :icon="['fas', 'tachometer-alt']" class="h-5 w-5 mr-3" />
              <span v-if="isSidebarOpen">Dashboard</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/companies" class="flex items-center justify-start p-2 rounded-md hover:bg-secondary hover:text-primary transition-colors duration-200">
              <font-awesome-icon :icon="['fas', 'building']" class="h-5 w-5 mr-3" />
              <span v-if="isSidebarOpen">Companies</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/subscriptions" class="flex items-center justify-start p-2 rounded-md hover:bg-secondary hover:text-primary transition-colors duration-200">
              <font-awesome-icon :icon="['fas', 'receipt']" class="h-5 w-5 mr-3" />
              <span v-if="isSidebarOpen">Subscriptions</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/revenue-chart" class="flex items-center justify-start p-2 rounded-md hover:bg-secondary hover:text-primary transition-colors duration-200">
              <font-awesome-icon :icon="['fas', 'chart-line']" class="h-5 w-5 mr-3" />
              <span v-if="isSidebarOpen">Revenue Chart</span>
            </router-link>
          </li>
          <li class="mb-2">
            <router-link to="/subscription-packages" class="flex items-center justify-start p-2 rounded-md hover:bg-secondary hover:text-primary transition-colors duration-200">
              <font-awesome-icon :icon="['fas', 'box-open']" class="h-5 w-5 mr-3" />
              <span v-if="isSidebarOpen">Manajemen Paket</span>
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
          <font-awesome-icon :icon="['fas', 'bars']" class="h-6 w-6" />
        </button>
        <h1 class="text-xl font-semibold">Selamat Datang, SuperAdmin!</h1>
        <div>
          <span class="text-text-muted">{{ authStore.user?.email }}</span>
        </div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-bg-base">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'SuperAdminDashboard',
  setup() {
    const router = useRouter();
    const isSidebarOpen = ref(false);
    const toast = useToast();
    const authStore = useAuthStore();

    const handleLogout = () => {
      authStore.clearAuth();
      // Menggunakan full page reload untuk memastikan semua state lama bersih
      window.location.href = '/auth';
    };

    const checkScreenSize = () => {
      if (window.innerWidth >= 768) { // md breakpoint
        isSidebarOpen.value = true;
      } else {
        isSidebarOpen.value = false;
      }
    };

    onMounted(() => {
      checkScreenSize(); // Initial check
      window.addEventListener('resize', checkScreenSize);
    });

    onUnmounted(() => {
      window.removeEventListener('resize', checkScreenSize);
    });

    const toggleSidebar = () => {
      if (window.innerWidth < 768) { // Only allow toggling on small screens
        isSidebarOpen.value = !isSidebarOpen.value;
      }
    };

    return {
      isSidebarOpen,
      handleLogout,
      authStore,
      toggleSidebar,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
