<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h1 class="text-2xl font-bold text-text-base mb-6">Pengaturan Lokasi Absensi</h1>

    <BaseDataTable
      :data="locations"
      :columns="locationColumns"
      :loading="isLoading"
      v-model:filters="filters"
      :globalFilterFields="['name']"
      searchPlaceholder="Cari Lokasi..."
    >
      <template #header-actions>
        <BaseButton @click="openAddModal">
          <i class="pi pi-plus"></i>Tambah Lokasi
        </BaseButton>
      </template>

      <template #column-latitude="{ item }">
        {{ item.latitude.toFixed(6) }}
      </template>

      <template #column-longitude="{ item }">
        {{ item.longitude.toFixed(6) }}
      </template>

      <template #column-actions="{ item }">
        <div class="flex items-center justify-center space-x-2">
          <BaseButton @click="openEditModal(item)" class="text-accent hover:opacity-80">
              <i class="pi pi-pencil"></i>
            </BaseButton>
            <BaseButton @click="deleteLocation(item.ID)" class="text-danger hover:opacity-80">
              <i class="pi pi-trash"></i>
          </BaseButton>
        </div>
      </template>
    </BaseDataTable>

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
            <BaseButton @click="searchLocation" type="button" :disabled="isSearching">
              <i v-if="!isSearching" class="pi pi-search"></i>
              <i v-else class="pi pi-spin pi-spinner"></i>
              {{ isSearching ? 'Mencari...' : 'Cari' }}
            </BaseButton>
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
            <i class="pi pi-times"></i> Batal
          </BaseButton>
          <BaseButton type="submit" :disabled="isSaving">
            <i v-if="!isSaving" class="pi pi-save"></i>
            <i v-else class="pi pi-spin pi-spinner"></i>
            {{ isSaving ? 'Menyimpan...' : 'Simpan' }}
          </BaseButton>
        </div>
      </form>
    </BaseModal>

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch } from 'vue';
import axios from 'axios';
import L from 'leaflet';
import { FilterMatchMode } from '@primevue/core/api';
import { useAuthStore } from '../../../stores/auth';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';
import ConfirmDialog from 'primevue/confirmdialog';
import BaseModal from '../../ui/BaseModal.vue';
import BaseInput from '../../ui/BaseInput.vue';
import BaseButton from '../../ui/BaseButton.vue';
import BaseDataTable from '../../ui/BaseDataTable.vue';

const authStore = useAuthStore();
const toast = useToast();
const confirm = useConfirm();

const locations = ref([]);
const isLoading = ref(true);
const isModalOpen = ref(false);
const modalTitle = ref('');
const isEditMode = ref(false);
const isSearching = ref(false);
const isSaving = ref(false);
const isDeleting = ref(false);

let map = null;
let marker = null;
let circle = null;

const searchQuery = ref('');

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const locationColumns = ref([
    { field: 'name', header: 'Nama Lokasi' },
    { field: 'latitude', header: 'Latitude' },
    { field: 'longitude', header: 'Longitude' },
    { field: 'radius', header: 'Radius (m)' },
    { field: 'actions', header: 'Aksi' }
]);

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
    toast.add({ severity: 'warning', summary: 'Warning', detail: 'Masukkan nama tempat atau alamat untuk mencari.', life: 3000 });
    return;
  }

  isSearching.value = true;
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
      toast.add({ severity: 'success', summary: 'Success', detail: `Lokasi ditemukan: ${firstResult.display_name}`, life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Lokasi tidak ditemukan. Coba kata kunci lain.', life: 3000 });
    }
  } catch (error) {
    console.error('Error searching location:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mencari lokasi. Coba lagi nanti.', life: 3000 });
  } finally {
    isSearching.value = false;
  }
};

const fetchLocations = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/company/locations');
    locations.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error("Error fetching locations:", error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data lokasi.', life: 3000 });
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
    toast.add({ severity: 'error', summary: 'Error', detail: 'Nama lokasi belum diisi.', life: 3000 });
    return;
  }

  if (!currentLocation.value.company_id) {
      currentLocation.value.company_id = authStore.companyId;
  }

  isSaving.value = true;
  try {
    if (isEditMode.value) {
      await axios.put(`/api/company/locations/${currentLocation.value.ID}`, currentLocation.value);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Lokasi berhasil diperbarui.', life: 3000 });
    } else {
      await axios.post('/api/company/locations', currentLocation.value);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Lokasi baru berhasil ditambahkan.', life: 3000 });
    }
    closeModal();
    fetchLocations();
  } catch (error) {
     console.error("Error saving location:", error);
     toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menyimpan lokasi. Pastikan semua field terisi.', life: 3000 });
  } finally {
    isSaving.value = false;
  }
};

const deleteLocation = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus lokasi ini? Tindakan ini tidak dapat dibatalkan.',
    header: 'Konfirmasi Hapus Lokasi',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      isDeleting.value = true;
      try {
        await axios.delete(`/api/company/locations/${id}`);
        toast.add({ severity: 'success', summary: 'Success', detail: 'Lokasi berhasil dihapus.', life: 3000 });
        fetchLocations();
      } catch (error) {
        console.error("Error deleting location:", error);
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus lokasi.', life: 3000 });
      } finally {
        isDeleting.value = false;
      }
    },
    reject: () => {
      toast.add({ severity: 'info', summary: 'Dibatalkan', detail: 'Penghapusan lokasi dibatalkan', life: 3000 });
    }
  });
};
</script>

<style scoped>
@import "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css";
</style>