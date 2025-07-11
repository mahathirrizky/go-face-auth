<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h1 class="text-2xl font-bold text-text-base mb-6">Pengaturan Lokasi Absensi</h1>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex justify-end">
      <button @click="openAddModal" class="btn btn-secondary">
        + Tambah Lokasi
      </button>
    </div>

    <!-- Tabel Lokasi -->
    <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
      <table class="min-w-full divide-y divide-bg-base">
        <thead class="bg-primary">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama Lokasi</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Latitude</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Longitude</th>
            <th scope="col" class="px-6 py-3 text-center text-xs font-medium text-white uppercase tracking-wider">Radius (m)</th>
            <th scope="col" class="px-6 py-3 text-center text-xs font-medium text-white uppercase tracking-wider">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-bg-base">
          <tr v-if="isLoading" class="text-center">
            <td colspan="5" class="px-6 py-4 text-text-muted">Loading...</td>
          </tr>
          <tr v-else-if="locations.length === 0">
             <td colspan="5" class="px-6 py-4 text-center text-text-muted">Belum ada lokasi yang ditambahkan.</td>
          </tr>
          <tr v-for="location in locations" :key="location.ID" class="hover:bg-bg-base">
            <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ location.name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ location.latitude.toFixed(6) }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ location.longitude.toFixed(6) }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-center text-text-muted">{{ location.radius }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-center">
              <div class="flex items-center justify-center space-x-2">
                <button @click="openEditModal(location)" class="text-accent hover:text-secondary">
                  <font-awesome-icon :icon="['fas', 'edit']" class="h-5 w-5" />
                </button>
                <button @click="deleteLocation(location.ID)" class="text-danger hover:opacity-80">
                  <font-awesome-icon :icon="['fas', 'trash-alt']" class="h-5 w-5" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Tambah/Edit Lokasi -->
    <div v-if="isModalOpen" class="fixed inset-0 bg-black bg-opacity-20 flex items-center justify-center z-50">
      <div class="bg-bg-muted p-8 rounded-lg shadow-lg w-full max-w-2xl">
        <h3 class="text-2xl font-bold text-text-base mb-6">{{ modalTitle }}</h3>
        <form @submit.prevent="saveLocation">
          <div class="mb-4">
            <label for="name" class="block text-text-muted text-sm font-bold mb-2">Nama Lokasi</label>
            <input v-model="currentLocation.name" type="text" id="name" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary" required>
          </div>
          <div class="mb-4">
            <label for="locationSearch" class="block text-text-muted text-sm font-bold mb-2">Cari Lokasi (Nama Tempat/Alamat)</label>
            <div class="flex space-x-2">
              <input
                type="text"
                id="locationSearch"
                v-model="searchQuery"
                @keyup.enter="searchLocation"
                placeholder="Contoh: Jakarta, Kantor Pusat, Jl. Sudirman 1"
                class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
              >
              <button @click="searchLocation" type="button" class="btn btn-secondary flex-shrink-0">Cari</button>
            </div>
          </div>
          <div id="map-container" class="mb-4 h-80 rounded-md overflow-hidden"></div>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label for="latitude" class="block text-text-muted text-sm font-bold mb-2">Latitude</label>
              <input v-model.number="currentLocation.latitude" type="number" step="any" id="latitude" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary" required>
            </div>
            <div>
              <label for="longitude" class="block text-text-muted text-sm font-bold mb-2">Longitude</label>
              <input v-model.number="currentLocation.longitude" type="number" step="any" id="longitude" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary" required>
            </div>
          </div>
          <div class="mb-6">
            <label for="radius" class="block text-text-muted text-sm font-bold mb-2">Radius (meter)</label>
            <input v-model.number="currentLocation.radius" type="number" id="radius" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary" required>
          </div>
          <div class="flex justify-end space-x-4">
            <button @click="closeModal" type="button" class="btn btn-outline-primary">
              Batal
            </button>
            <button @click="saveLocation" type="button" class="btn btn-secondary">
              Simpan
            </button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import axios from 'axios';
import L from 'leaflet';
import { useAuthStore } from '../../../stores/auth';
import { useToast } from 'vue-toastification';

const authStore = useAuthStore();
const toast = useToast();

const locations = ref([]);
const isLoading = ref(true);
const isModalOpen = ref(false);
const modalTitle = ref('');
const isEditMode = ref(false);

let map = null;
let marker = null;

const searchQuery = ref(''); // Tambahkan ini

const initialLocationState = {
  ID: null,
  name: '',
  latitude: -6.2088, // Default to Jakarta
  longitude: 106.8456,
  radius: 100,
  company_id: authStore.companyId
};

const currentLocation = ref({ ...initialLocationState });

onMounted(() => {
  fetchLocations();
});

const initMap = async () => {
  await nextTick();
  if (map) {
    map.remove();
    map = null;
  }
  
  const { latitude, longitude } = currentLocation.value;
  
  map = L.map('map-container').setView([latitude, longitude], 15);

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map);

  // Define a custom icon using L.divIcon
  const customIcon = L.divIcon({
    className: 'custom-map-marker', // A custom class for styling
    // Membuat ikon sangat besar dan berwarna hijau cerah untuk debugging
    html: '<i class="fas fa-map-marker-alt" style="font-size: 4rem; color: #00FF00;"></i>',
    iconSize: [60, 84], // Sesuaikan ukuran untuk ikon yang lebih besar
    iconAnchor: [30, 84], // Sesuaikan anchor point
    popupAnchor: [0, -70] // Sesuaikan popup anchor
  });

  marker = L.marker([latitude, longitude], { icon: customIcon, draggable: true }).addTo(map);

  marker.on('dragend', (event) => {
    const position = event.target.getLatLng();
    currentLocation.value.latitude = position.lat;
    currentLocation.value.longitude = position.lng;
  });

  // Tambahkan event listener untuk klik peta
  map.on('click', (e) => {
    currentLocation.value.latitude = e.latlng.lat;
    currentLocation.value.longitude = e.latlng.lng;
    marker.setLatLng(e.latlng); // Pindahkan marker ke lokasi yang diklik
  });

  // A hack to fix map rendering issue in modal
  setTimeout(() => {
      map.invalidateSize();
  }, 100);
};

// Tambahkan fungsi searchLocation ini
const searchLocation = async () => {
  if (!searchQuery.value.trim()) {
    toast.warning('Masukkan nama tempat atau alamat untuk mencari.');
    return;
  }

  try {
    const response = await axios.get(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(searchQuery.value)}`);
    if (response.data && response.data.length > 0) {
      const firstResult = response.data[0];
      const lat = parseFloat(firstResult.lat);
      const lon = parseFloat(firstResult.lon);

      currentLocation.value.latitude = lat;
      currentLocation.value.longitude = lon;

      map.setView([lat, lon], 15); // Pindahkan peta ke lokasi baru
      marker.setLatLng([lat, lon]); // Pindahkan marker ke lokasi baru
      toast.success(`Lokasi ditemukan: ${firstResult.display_name}`);
    } else {
      toast.error('Lokasi tidak ditemukan. Coba kata kunci lain.');
    }
  } catch (error) {
    console.error('Error searching location:', error);
    toast.error('Gagal mencari lokasi. Coba lagi nanti.');
  }
};

const fetchLocations = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/company/locations');
    locations.value = response.data;
  } catch (error) {
    console.error("Error fetching locations:", error);
    toast.error('Gagal memuat data lokasi.');
  } finally {
    isLoading.value = false;
  }
};

const openAddModal = () => {
  isEditMode.value = false;
  modalTitle.value = 'Tambah Lokasi Baru';
  currentLocation.value = { ...initialLocationState, company_id: authStore.companyId };
  searchQuery.value = ''; // Bersihkan search query saat membuka modal
  isModalOpen.value = true;
  nextTick(() => {
    initMap();
  });
};

const openEditModal = (location) => {
  isEditMode.value = true;
  modalTitle.value = 'Edit Lokasi';
  currentLocation.value = { ...location };
  searchQuery.value = ''; // Bersihkan search query saat membuka modal
  isModalOpen.value = true;
  nextTick(() => {
    initMap();
  });
};

const closeModal = () => {
  isModalOpen.value = false;
  if (map) {
    map.remove();
    map = null;
  }
};

const saveLocation = async () => {
  // Ensure company_id is set
  if (!currentLocation.value.company_id) {
      currentLocation.value.company_id = authStore.companyId;
  }

  try {
    if (isEditMode.value) {
      await axios.put(`/api/company/locations/${currentLocation.value.ID}`, currentLocation.value);
      toast.success('Lokasi berhasil diperbarui.');
    } else {
      await axios.post('/api/company/locations', currentLocation.value);
      toast.success('Lokasi baru berhasil ditambahkan.');
    }
    closeModal();
    fetchLocations(); // Refresh list
  } catch (error) {
     console.error("Error saving location:", error);
     toast.error('Gagal menyimpan lokasi. Pastikan semua field terisi.');
  }
};

const deleteLocation = async (id) => {
  if (confirm('Apakah Anda yakin ingin menghapus lokasi ini?')) {
    try {
      await axios.delete(`/api/company/locations/${id}`);
      toast.success('Lokasi berhasil dihapus.');
      fetchLocations(); // Refresh list
    } catch (error) {
      console.error("Error deleting location:", error);
      toast.error('Gagal menghapus lokasi.');
    }
  }
};

</script>