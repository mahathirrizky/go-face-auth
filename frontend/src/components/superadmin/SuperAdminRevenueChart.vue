<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Grafik Pendapatan</h2>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md mb-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pilih Rentang Tanggal</h3>
      <div class="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-4">
        <div class="flex flex-col">
          <label for="startDate" class="text-text-muted text-sm mb-1">Tanggal Mulai:</label>
          <flat-pickr
            v-model="startDate"
            :config="startDatePickerConfig"
            placeholder="Pilih Tanggal Mulai"
            class="w-full md:w-auto px-3 py-2 border border-gray-300 rounded-md text-text-base bg-bg-base focus:outline-none focus:ring-secondary focus:border-secondary"
          ></flat-pickr>
        </div>
        <div class="flex flex-col">
          <label for="endDate" class="text-text-muted text-sm mb-1">Tanggal Akhir:</label>
          <flat-pickr
            v-model="endDate"
            :config="endDatePickerConfig"
            placeholder="Pilih Tanggal Akhir"
            class="w-full md:w-auto px-3 py-2 border border-gray-300 rounded-md text-text-base bg-bg-base focus:outline-none focus:ring-secondary focus:border-secondary"
          ></flat-pickr>
        </div>
        
      </div>
    </div>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pendapatan Bulanan</h3>
      <LineChart :data="chartData" :chart-options="chartOptions" />
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import flatPickr from 'vue-flatpickr-component';
import 'flatpickr/dist/flatpickr.css';
import { Line } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement } from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement);

export default {
  name: 'SuperAdminRevenueChart',
  components: {
    flatPickr,
    LineChart: Line,
  },
  setup() {
    const toast = useToast();
    const startDate = ref('');
    const endDate = ref('');
    const revenueData = ref([]);

    const startDatePickerConfig = ref({
      dateFormat: 'Y-m-d',
      maxDate: new Date(), // Cannot select a date after today
    });

    const endDatePickerConfig = ref({
      dateFormat: 'Y-m-d',
      maxDate: new Date(), // Cannot select a date after today
      minDate: null, // Will be updated by watch effect
    });

    const fetchRevenueData = async () => {
      try {
        const response = await axios.get('/api/superadmin/revenue-summary', {
          params: {
            start_date: startDate.value,
            end_date: endDate.value,
          },
        });
        if (response.data && response.data.status === 'success') {
          revenueData.value = response.data.data;
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

      const labels = [];
      const data = [];

      // Create a map for quick lookup of revenue data
      const revenueMap = new Map();
      if (revenueData.value) { // Add null check here
        revenueData.value.forEach(item => {
          revenueMap.set(`${item.year}-${item.month}`, item.total_revenue);
        });
      }

      // Determine the date range for labels
      let start = startDate.value ? new Date(startDate.value) : new Date();
      let end = endDate.value ? new Date(endDate.value) : new Date();

      // Populate labels and data based on the selected date range
      let currentDate = new Date(start.getFullYear(), start.getMonth(), 1); // Start from the beginning of the start month
      while (currentDate <= end) {
        const monthName = allMonths[currentDate.getMonth()];
        const year = currentDate.getFullYear();
        const monthKey = `${year}-${(currentDate.getMonth() + 1).toString().padStart(2, '0')}`;

        labels.push(`${monthName} ${year}`);
        data.push(revenueMap.get(monthKey) || 0);

        currentDate.setMonth(currentDate.getMonth() + 1); // Move to next month
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
            color: '#a0aec0',
          },
        },
        tooltip: {
          callbacks: {
            label: (context) => {
              let label = context.dataset.label || '';
              if (label) {
                label += ': ';
              }
              if (context.parsed.y !== null) {
                label += new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', currencyDisplay: 'symbol' }).format(context.parsed.y);
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
          min: 0,
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

    // Initial fetch when component is mounted
    onMounted(() => {
      const today = new Date();
      const firstDayOfYear = new Date(today.getFullYear(), 0, 1);
      startDate.value = firstDayOfYear.toISOString().slice(0, 10);
      endDate.value = today.toISOString().slice(0, 10);
      fetchRevenueData();
    });

    // Watch for changes in startDate or endDate and fetch data
    watch([startDate, endDate], ([newStartDate, newEndDate], [oldStartDate, oldEndDate]) => {
      if ((newStartDate !== oldStartDate || newEndDate !== oldEndDate) && newStartDate && newEndDate) {
        fetchRevenueData();
      }
    });

    // Watch startDate to update minDate of endDatePickerConfig
    watch(startDate, (newStartDate) => {
      endDatePickerConfig.value.minDate = newStartDate;
    });

    return {
      startDate,
      endDate,
      startDatePickerConfig,
      endDatePickerConfig,
      fetchRevenueData,
      chartData,
      chartOptions,
    };
}
}
</script>

<style scoped>
/* Tailwind handles styling */
</style>
