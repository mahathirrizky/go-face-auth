<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h1 class="text-2xl font-bold text-text-base mb-6">Pengaturan Lokasi Absensi</h1>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex justify-end">
      <BaseButton @click="openAddModal">
        + Tambah Lokasi
      </BaseButton>
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
                <BaseButton @click="openEditModal(location)" class="text-accent hover:text-secondary">
                  <font-awesome-icon :icon="['fas', 'edit']" class="h-5 w-5" />
                </BaseButton>
                <BaseButton @click="deleteLocation(location.ID)" class="text-danger hover:opacity-80">
                  <font-awesome-icon :icon="['fas', 'trash-alt']" class="h-5 w-5" />
                </BaseButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Tambah/Edit Lokasi -->
    <BaseModal :isOpen="isModalOpen" @close="closeModal" :title="modalTitle" maxWidth="2xl">
      <form @submit.prevent="saveLocation">
        <BaseInput
          id="name"
          label="Nama Lokasi"
          v-model="currentLocation.name"
          required
        />
        <div class="mb-4">
          <label for="locationSearch" class="block text-text-muted text-sm font-bold mb-2">Cari Lokasi (Nama Tempat/Alamat)</label>
          <div class="flex space-x-2">
            <BaseInput
              id="locationSearch"
              v-model="searchQuery"
              @keyup.enter="searchLocation"
              placeholder="Contoh: Jakarta, Kantor Pusat, Jl. Sudirman 1"
              class="w-full"
              :label-sr-only="true"
            />
            <BaseButton @click="searchLocation" type="button">Cari</BaseButton>
          </div>
        </div>
        <div id="map-container" class="mb-4 h-80 rounded-md overflow-hidden"></div>
        <div class="grid grid-cols-2 gap-4 mb-4">
          <BaseInput
            id="latitude"
            label="Latitude"
            v-model.number="currentLocation.latitude"
            type="number"
            step="any"
            required
          />
          <BaseInput
            id="longitude"
            label="Longitude"
            v-model.number="currentLocation.longitude"
            type="number"
            step="any"
            required
          />
        </div>
        <BaseInput
          id="radius"
          label="Radius (meter)"
          v-model.number="currentLocation.radius"
          type="number"
          required
        />
        <div class="flex justify-end space-x-4 mt-6">
          <BaseButton @click="closeModal" type="button" class="btn-outline-primary">
            Batal
          </BaseButton>
          <BaseButton type="submit">
            Simpan
          </BaseButton>
        </div>
      </form>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal :isOpen="showDeleteConfirmModal" @close="cancelDelete" title="Konfirmasi Hapus" maxWidth="sm">
      <p class="text-text-muted mb-6 text-center">Apakah Anda yakin ingin menghapus lokasi ini? Tindakan ini tidak dapat dibatalkan.</p>
      <template #footer>
        <BaseButton @click="cancelDelete" class="btn-outline-primary">Batal</BaseButton>
        <BaseButton @click="confirmDelete" class="btn-danger">Ya, Hapus</BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch } from 'vue';
import axios from 'axios';
import L from 'leaflet';
import { useAuthStore } from '../../../stores/auth';
import { useToast } from 'vue-toastification';
import BaseModal from '../../ui/BaseModal.vue';
import BaseInput from '../../ui/BaseInput.vue';
import BaseButton from '../../ui/BaseButton.vue';

const authStore = useAuthStore();
const toast = useToast();

const locations = ref([]);
const isLoading = ref(true);
const isModalOpen = ref(false);
const modalTitle = ref('');
const isEditMode = ref(false);

const showDeleteConfirmModal = ref(false);
const locationToDeleteId = ref(null);

let map = null;
let marker = null;
let circle = null;

const searchQuery = ref('');

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
  if (map) { // Destroy existing map instance if it exists
    map.remove();
    map = null;
  }
  if (circle) {
    circle.remove();
    circle = null;
  }
  
  const { latitude, longitude, radius } = currentLocation.value;
  
  map = L.map('map-container').setView([latitude, longitude], 15);

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map);

  const customIcon = L.icon({
    iconUrl: '/images/marker-icon.png',
    shadowUrl: '/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
  });

  marker = L.marker([latitude, longitude], { icon: customIcon, draggable: true }).addTo(map);

  circle = L.circle([latitude, longitude], {
    color: 'blue',
    fillColor: '#0000ff',
    fillOpacity: 0.2,
    radius: radius
  }).addTo(map);

  marker.on('dragend', (event) => {
    const position = event.target.getLatLng();
    currentLocation.value.latitude = position.lat;
    currentLocation.value.longitude = position.lng;
    circle.setLatLng(position);
  });

  map.on('click', (e) => {
    currentLocation.value.latitude = e.latlng.lat;
    currentLocation.value.longitude = e.latlng.lng;
    marker.setLatLng(e.latlng);
    circle.setLatLng(e.latlng);
  });

  setTimeout(() => {
      map.invalidateSize();
  }, 100);
};

watch(() => currentLocation.value.radius, (newRadius) => {
  if (circle) {
    circle.setRadius(newRadius);
  }
});

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

      map.setView([lat, lon], 15);
      marker.setLatLng([lat, lon]);
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
  searchQuery.value = '';
  isModalOpen.value = true;
  nextTick(() => {
    initMap();
  });
};

const openEditModal = (location) => {
  isEditMode.value = true;
  modalTitle.value = 'Edit Lokasi';
  currentLocation.value = { ...location };
  searchQuery.value = '';
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
  if (circle) {
    circle.remove();
    circle = null;
  }
};

const saveLocation = async () => {
  if (!currentLocation.value.name.trim()) {
    toast.error('Nama lokasi belum diisi.');
    return;
  }

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
    fetchLocations();
  } catch (error) {
     console.error("Error saving location:", error);
     toast.error('Gagal menyimpan lokasi. Pastikan semua field terisi.');
  }
};

const deleteLocation = async (id) => {
  locationToDeleteId.value = id;
  showDeleteConfirmModal.value = true;
};

const confirmDelete = async () => {
  try {
    await axios.delete(`/api/company/locations/${locationToDeleteId.value}`);
    toast.success('Lokasi berhasil dihapus.');
    fetchLocations();
  } catch (error) {
    console.error("Error deleting location:", error);
    toast.error('Gagal menghapus lokasi.');
  } finally {
    showDeleteConfirmModal.value = false;
    locationToDeleteId.value = null;
  }
};

const cancelDelete = () => {
  showDeleteConfirmModal.value = false;
  locationToDeleteId.value = null;
};
</script>
