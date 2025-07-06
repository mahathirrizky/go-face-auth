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
        <h1 class="text-xl font-semibold">Selamat Datang, Admin <span v-if="companyName"> {{ companyName }}</span>!</h1>
        <div>
          <!-- User Dropdown or other header elements -->
          <span class="text-text-muted">admin@example.com</span>
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
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'vue-toastification'; // Import useToast
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'AdminDashboard',
  setup() {
    const router = useRouter();
    const companyName = ref('Nama Perusahaan Anda'); // Placeholder for company name
    const isSidebarOpen = ref(false); // State for sidebar visibility
    const toast = useToast(); // Initialize toast
    const authStore = useAuthStore();

    const handleLogout = () => {
      console.log('Logging out...');
      authStore.clearAuth(); // Clear all auth data from Pinia store
      axios.defaults.headers.common['Authorization'] = ''; // Clear Authorization header
      router.push('/'); // Redirect to root page
    };

    const fetchCompanyDetails = async () => {
      try {
        const response = await axios.get(`/api/company-details`);
        console.log('Full API response object:', response); // Added for debugging
        if (response.data && response.data.data && response.data.data.name) {
          companyName.value = response.data.data.name;
          authStore.setCompanyId(response.data.data.id); // Save company ID to authStore
          authStore.setCompanyName(response.data.data.name); // Save company name to authStore
          authStore.setCompanyAddress(response.data.data.address); // Save company address to authStore
          authStore.setAdminEmail(response.data.data.admin_email); // Save admin email to authStore
        } else {
          toast.error('Failed to fetch company details.');
        }
      } catch (error) {
        console.error('Error fetching company details:', error);
        let message = 'Failed to load company details.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
        // Optionally, redirect to login if token is invalid
        if (error.response && error.response.status === 401) {
          localStorage.removeItem('admin_token');
          axios.defaults.headers.common['Authorization'] = '';
          router.push('/login');
        }
      }
    };

    onMounted(() => {
      fetchCompanyDetails();
    });

    return {
      companyName,
      isSidebarOpen,
      handleLogout,
    };
  },
};
</script>
<style scoped>
/* Tailwind handles styling */
</style>

