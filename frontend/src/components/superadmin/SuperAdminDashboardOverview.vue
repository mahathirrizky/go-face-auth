<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Ringkasan Dashboard SuperAdmin</h2>

    <div class="flex flex-col lg:flex-row lg:space-x-6 mt-8">
      <!-- Chart on the left -->
      <div class="bg-bg-muted p-6 rounded-lg shadow-md mb-6 lg:mb-0 lg:w-1/2 max-h-96">
        <h3 class="text-xl font-semibold text-text-base mb-4">Pendapatan Bulanan</h3>
        <LineChart :data="chartData" :chart-options="chartOptions" />
      </div>

      <!-- Cards on the right -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6 lg:w-1/2">
        <!-- Total Companies Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Total Perusahaan</h3>
          <p class="text-3xl font-bold text-secondary">{{ summary.total_companies }}</p>
        </div>
        <!-- Active Subscriptions Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Langganan Aktif</h3>
          <p class="text-3xl font-bold text-accent">{{ summary.active_subscriptions }}</p>
        </div>
        <!-- Expired Subscriptions Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Langganan Kedaluwarsa</h3>
          <p class="text-3xl font-bold text-danger">{{ summary.expired_subscriptions }}</p>
        </div>
        <!-- Trial Subscriptions Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Langganan Percobaan</h3>
          <p class="text-3xl font-bold text-blue-400">{{ summary.trial_subscriptions }}</p>
        </div>
      </div>
    </div>

    <div class="mt-8 bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Aktivitas Terbaru</h3>
      <ul>
        <li v-for="activity in recentActivities" :key="activity.id" class="border-b border-bg-base py-2 text-gray-300">
          {{ activity.description }} - {{ new Date(activity.timestamp).toLocaleString() }}
        </li>
        <li v-if="recentActivities.length === 0" class="py-2 text-gray-300">Tidak ada aktivitas terbaru.</li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';
import { Line } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement } from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement);

export default {
  name: 'SuperAdminDashboardOverview',
  components: {
    LineChart: Line,
  },
  setup() {
    const summary = ref({
      total_companies: 0,
      active_subscriptions: 0,
      expired_subscriptions: 0,
      trial_subscriptions: 0,
    });
    const recentActivities = ref([]);
    const revenueData = ref([]);
    const toast = useToast();
    const authStore = useAuthStore();
    let ws = null;

    const fetchDashboardSummary = async () => {
      try {
        const response = await axios.get('/api/superadmin/dashboard-summary');
        if (response.data && response.data.data) {
          summary.value = response.data.data;
          recentActivities.value = response.data.data.recent_activities || [];
        } else {
          toast.error('Failed to fetch superadmin dashboard summary.');
        }
      } catch (error) {
        console.error('Error fetching superadmin dashboard summary:', error);
        let message = 'Failed to load superadmin dashboard summary.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    const fetchRevenueData = async () => {
      try {
        const response = await axios.get('/api/superadmin/revenue-summary');
        if (response.data && response.data.status === 'success') {
          revenueData.value = response.data.data || [];
        } else {
          toast.error(response.data.message || 'Failed to fetch revenue data.');
        }
      } catch (error) {
        console.error('Error fetching revenue data:', error);
        toast.error('An error occurred while fetching revenue data.');
      }
    };

    const chartData = computed(() => {
      const allMonths = [
        'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
        'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
      ];

      const currentYear = new Date().getFullYear();
      const labels = [];
      const data = [];

      // Map revenueData to a more accessible format
      const revenueMap = new Map();
      revenueData.value.forEach(item => {
        // Ensure we only consider data for the current year
        if (item.year === currentYear) {
          revenueMap.set(`${item.month}-${item.year}`, item.total_revenue);
        }
      });

      // Populate labels and data for all 12 months of the current year
      for (let i = 0; i < 12; i++) {
        const monthName = allMonths[i];
        // Format month to YYYY-MM for matching with backend data
        const monthNumber = (i + 1).toString().padStart(2, '0');
        const monthYearKey = `${currentYear}-${monthNumber}`;

        labels.push(`${monthName} ${currentYear}`);
        data.push(revenueMap.get(monthYearKey) || 0);
      }

      return {
        labels: labels,
        datasets: [
          {
            label: 'Pendapatan',
            backgroundColor: '#42A5F5',
            borderColor: '#42A5F5',
            data: data,
            tension: 0.4,
          },
        ],
      };
    });

    const chartOptions = ref({
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          labels: {
            color: '#a0aec0', // text-text-muted
          },
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || '';
              if (label) {
                label += ': ';
              }
              if (context.parsed.y !== null) {
                label += new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(context.parsed.y);
              }
              return label;
            }
          }
        }
      },
      scales: {
        x: {
          ticks: {
            color: '#a0aec0',
          },
          grid: {
            color: 'rgba(160, 174, 192, 0.1)',
          },
        },
        y: {
          min: 0, // Ensure Y-axis starts from 0
          ticks: {
            color: '#a0aec0',
            callback: function(value) {
              return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', currencyDisplay: 'symbol' }).format(value);
            }
          },
          grid: {
            color: 'rgba(160, 174, 192, 0.1)',
          },
        },
      },
    });

    const connectWebSocket = () => {
      if (!authStore.token) {
        console.warn('No auth token found, cannot establish WebSocket connection.');
        return;
      }

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//api.4commander.my.id/ws/superadmin-dashboard?token=${authStore.token}`;

      ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log('SuperAdmin WebSocket connected.');
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data) {
          summary.value = data;
          recentActivities.value = data.recent_activities || [];
          revenueData.value = data.monthly_revenue || []; // Update revenue data from WebSocket
        }
      };

      ws.onclose = (event) => {
        console.log('SuperAdmin WebSocket disconnected:', event.code, event.reason);
        if (event.code !== 1000) {
          setTimeout(() => {
            console.log('Attempting to reconnect SuperAdmin WebSocket...');
            connectWebSocket();
          }, 3000);
        }
      };

      ws.onerror = (error) => {
        console.error('SuperAdmin WebSocket error:', error);
        toast.error('SuperAdmin WebSocket connection error. Dashboard updates may be delayed.');
      };
    };

    onMounted(() => {
      fetchDashboardSummary();
      fetchRevenueData();
      connectWebSocket();
    });

    onUnmounted(() => {
      if (ws) {
        ws.close();
      }
    });

    return {
      summary,
      recentActivities,
      chartData,
      chartOptions,
    };
  },
};

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(value);
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
