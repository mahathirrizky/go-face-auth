<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <ConfirmDialog />
    <h1 class="text-2xl font-bold text-text-base mb-6">Pengaturan Lokasi Absensi</h1>

    <DataTable
      :value="locations"
      :loading="isLoading"
      v-model:filters="filters"
      :globalFilterFields="['name']"
      paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
      dataKey="ID"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <Button @click="openAddModal" icon="pi pi-plus" label="Tambah Lokasi" />
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Lokasi..." />
          </IconField>
        </div>
      </template>

      <Column field="name" header="Nama Lokasi" :sortable="true"></Column>
      <Column field="latitude" header="Latitude" :sortable="true">
        <template #body="{ data }">{{ data.latitude.toFixed(6) }}</template>
      </Column>
      <Column field="longitude" header="Longitude" :sortable="true">
        <template #body="{ data }">{{ data.longitude.toFixed(6) }}</template>
      </Column>
      <Column field="radius" header="Radius (m)" :sortable="true"></Column>
      <Column header="Aksi" style="width: 10rem; text-align: center">
        <template #body="{ data }">
          <Button @click="openEditModal(data)" icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" />
          <Button @click="deleteLocation(data.ID)" icon="pi pi-trash" class="p-button-rounded p-button-danger" />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="isModalOpen" :header="modalTitle" :modal="true" class="w-full max-w-2xl" @after-hide="onModalHide" @show="onModalShow">
      <form @submit.prevent="saveLocation" class="p-fluid mt-6">
        <div class="space-y-6">
          <FloatLabel variant="on">
            <InputText id="name" v-model="currentLocation.name" required fluid/>
            <label for="name">Nama Lokasi</label>
          </FloatLabel>

          <FloatLabel variant="on">

           
            <IconField iconPosition="right">
              <InputText id="locationSearch" v-model="searchQuery" @keydown.enter.prevent="searchLocation" fluid />
              <InputIcon :class="['pi', isSearching ? 'pi-spin pi-spinner' : 'pi-search']" @click="searchLocation" />
            </IconField>
            <label for="locationSearch">Cari Lokasi (Nama Tempat/Alamat)</label>
          
         
        </FloatLabel>

          <div class="h-80 w-full rounded-lg overflow-hidden border-2 border-surface-200">
            <div id="map-container" class="h-full w-full"></div>
          </div>

          <div class="grid grid-cols-3 gap-9 ">
            <FloatLabel variant="on">
              <InputNumber id="latitude" v-model="currentLocation.latitude" mode="decimal" :minFractionDigits="6" :maxFractionDigits="6" required />
              <label for="latitude">Latitude</label>
            </FloatLabel>
            <FloatLabel variant="on">
              <InputNumber id="longitude" v-model="currentLocation.longitude" mode="decimal" :minFractionDigits="6" :maxFractionDigits="6" required />
              <label for="longitude">Longitude</label>
            </FloatLabel>
            <FloatLabel variant="on">
              <InputNumber id="radius" v-model="currentLocation.radius" required suffix=" meter" />
              <label for="radius">Radius</label>
            </FloatLabel>
          </div>
        </div>

        <div class="flex justify-end space-x-2 mt-8">
          <Button @click="closeModal" type="button" label="Batal" class="p-button-text"/>
          <Button type="submit" :loading="isSaving" :label="isSaving ? 'Menyimpan...' : 'Simpan'" />
        </div>
      </form>
    </Dialog>
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
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Toast from 'primevue/toast';
import ConfirmDialog from 'primevue/confirmdialog';
import FloatLabel from 'primevue/floatlabel';

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

let map = null;
let marker = null;
let circle = null;

const searchQuery = ref('');

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const initialLocationState = {
  ID: null,
  name: '',
  latitude: -6.2088,
  longitude: 106.8456,
  radius: 100,
  company_id: authStore.companyId
};

const currentLocation = ref({ ...initialLocationState });

onMounted(() => {
  fetchLocations();
});

const onModalShow = () => {
  initMap();
};

const initMap = async () => {
  await nextTick();
  const mapContainer = document.getElementById('map-container');
  
  if (mapContainer && !map) {
    const { latitude, longitude, radius } = currentLocation.value;
    map = L.map(mapContainer).setView([latitude, longitude], 15);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    const customIcon = L.icon({
      iconUrl: '/images/marker-icon.png',
      shadowUrl: '/images/marker-shadow.png',
      iconSize: [25, 41], iconAnchor: [12, 41], popupAnchor: [1, -34], shadowSize: [41, 41]
    });

    marker = L.marker([latitude, longitude], { icon: customIcon, draggable: true }).addTo(map);
    circle = L.circle([latitude, longitude], { radius, color: 'blue', fillColor: '#00f', fillOpacity: 0.2 }).addTo(map);

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
  }

  if (map) {
    setTimeout(() => {
        map.invalidateSize()
    }, 100);
  }
};

watch(() => currentLocation.value.radius, (newRadius) => {
  if (circle) circle.setRadius(newRadius);
});

const searchLocation = async () => {
  if (!searchQuery.value.trim()) {
    toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Masukkan nama tempat atau alamat untuk mencari.', life: 3000 });
    return;
  }
  isSearching.value = true;
  try {
    const response = await axios.get(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(searchQuery.value)}`, { withCredentials: false });
    if (response.data && response.data.length > 0) {
      const { lat, lon, display_name } = response.data[0];
      currentLocation.value.latitude = parseFloat(lat);
      currentLocation.value.longitude = parseFloat(lon);
      map.setView([lat, lon], 15);
      marker.setLatLng([lat, lon]);
      toast.add({ severity: 'success', summary: 'Lokasi Ditemukan', detail: display_name, life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Gagal', detail: 'Lokasi tidak ditemukan.', life: 3000 });
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mencari lokasi.', life: 3000 });
  } finally {
    isSearching.value = false;
  }
};

const fetchLocations = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/company/locations');
    locations.value = response.data.data || [];
  } catch (error) {
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
};

const openEditModal = (location) => {
  isEditMode.value = true;
  modalTitle.value = 'Edit Lokasi';
  currentLocation.value = { ...location };
  searchQuery.value = '';
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
};

const onModalHide = () => {
    if (map) {
        map.remove();
        map = null;
    }
}

const saveLocation = async () => {
  if (!currentLocation.value.name.trim()) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Nama lokasi wajib diisi.', life: 3000 });
    return;
  }
  isSaving.value = true;
  try {
    const url = isEditMode.value ? `/api/company/locations/${currentLocation.value.ID}` : '/api/company/locations';
    const method = isEditMode.value ? 'put' : 'post';
    await axios[method](url, currentLocation.value);
    toast.add({ severity: 'success', summary: 'Berhasil', detail: `Lokasi berhasil ${isEditMode.value ? 'diperbarui' : 'disimpan'}.`, life: 3000 });
    closeModal();
    fetchLocations();
  } catch (error) {
     toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Gagal menyimpan lokasi.', life: 3000 });
  } finally {
    isSaving.value = false;
  }
};

const deleteLocation = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus lokasi ini?',
    header: 'Konfirmasi Hapus',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await axios.delete(`/api/company/locations/${id}`);
        toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Lokasi berhasil dihapus.', life: 3000 });
        fetchLocations();
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus lokasi.', life: 3000 });
      }
    }
  });
};
</script>

<style>
@import "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css";
</style>

