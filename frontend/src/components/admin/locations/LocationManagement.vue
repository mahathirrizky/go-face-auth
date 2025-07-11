<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold mb-6">Pengaturan Lokasi Absensi</h1>

    <div class="mb-4">
      <button @click="openAddModal" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
        + Tambah Lokasi
      </button>
    </div>

    <!-- Tabel Lokasi -->
    <div class="bg-white shadow-md rounded my-6">
      <table class="min-w-max w-full table-auto">
        <thead>
          <tr class="bg-gray-200 text-gray-600 uppercase text-sm leading-normal">
            <th class="py-3 px-6 text-left">Nama Lokasi</th>
            <th class="py-3 px-6 text-left">Latitude</th>
            <th class="py-3 px-6 text-left">Longitude</th>
            <th class="py-3 px-6 text-center">Radius (m)</th>
            <th class="py-3 px-6 text-center">Aksi</th>
          </tr>
        </thead>
        <tbody class="text-gray-600 text-sm font-light">
          <tr v-if="isLoading" class="text-center">
            <td colspan="5" class="py-4">Loading...</td>
          </tr>
          <tr v-else-if="locations.length === 0">
             <td colspan="5" class="py-4 text-center">Belum ada lokasi yang ditambahkan.</td>
          </tr>
          <tr v-for="location in locations" :key="location.ID" class="border-b border-gray-200 hover:bg-gray-100">
            <td class="py-3 px-6 text-left whitespace-nowrap">{{ location.name }}</td>
            <td class="py-3 px-6 text-left">{{ location.latitude.toFixed(6) }}</td>
            <td class="py-3 px-6 text-left">{{ location.longitude.toFixed(6) }}</td>
            <td class="py-3 px-6 text-center">{{ location.radius }}</td>
            <td class="py-3 px-6 text-center">
              <div class="flex item-center justify-center">
                <button @click="openEditModal(location)" class="w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center mr-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.5L15.232 5.232z" /></svg>
                </button>
                <button @click="deleteLocation(location.ID)" class="w-8 h-8 rounded-full bg-red-500 text-white flex items-center justify-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Tambah/Edit Lokasi -->
    <div v-if="isModalOpen" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
        </div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">{{ modalTitle }}</h3>
            <form @submit.prevent="saveLocation">
              <div class="mb-4">
                <label for="name" class="block text-gray-700 text-sm font-bold mb-2">Nama Lokasi</label>
                <input v-model="currentLocation.name" type="text" id="name" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" required>
              </div>
              <div id="map-container" class="mb-4 h-80 z-0"></div>
              <div class="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label for="latitude" class="block text-gray-700 text-sm font-bold mb-2">Latitude</label>
                  <input v-model.number="currentLocation.latitude" type="number" step="any" id="latitude" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" required>
                </div>
                <div>
                  <label for="longitude" class="block text-gray-700 text-sm font-bold mb-2">Longitude</label>
                  <input v-model.number="currentLocation.longitude" type="number" step="any" id="longitude" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" required>
                </div>
              </div>
              <div class="mb-4">
                <label for="radius" class="block text-gray-700 text-sm font-bold mb-2">Radius (meter)</label>
                <input v-model.number="currentLocation.radius" type="number" id="radius" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" required>
              </div>
            </form>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button @click="saveLocation" type="button" class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none sm:ml-3 sm:w-auto sm:text-sm">
              Simpan
            </button>
            <button @click="closeModal" type="button" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm">
              Batal
            </button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import axios from 'axios';
import L from 'leaflet';
import { useAuthStore } from '../../../../stores/auth';
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

  marker = L.marker([latitude, longitude], { draggable: true }).addTo(map);

  marker.on('dragend', (event) => {
    const position = event.target.getLatLng();
    currentLocation.value.latitude = position.lat;
    currentLocation.value.longitude = position.lng;
  });

  // A hack to fix map rendering issue in modal
  setTimeout(() => {
      map.invalidateSize();
  }, 100);
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
  isModalOpen.value = true;
  nextTick(() => {
    initMap();
  });
};

const openEditModal = (location) => {
  isEditMode.value = true;
  modalTitle.value = 'Edit Lokasi';
  currentLocation.value = { ...location };
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