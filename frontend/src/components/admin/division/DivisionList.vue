<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Divisi</h2>

    <BaseDataTable
      :data="divisions"
      :columns="divisionColumns"
      :loading="isLoading"
      v-model:filters="filters"
      :globalFilterFields="['name']"
      searchPlaceholder="Cari Divisi..."
      editModeType="cell"
      editKey="id"
      @cell-edit-complete="onCellEditComplete"
    >
      <template #header-actions>
        <BaseButton @click="openAddModal">
          <i class="pi pi-plus"></i> Tambah Divisi
        </BaseButton>
      </template>

      <template #editor-name="{ data, field }">
        <BaseInput v-model="data[field]" :id="`edit-${field}`" :name="field" autofocus />
      </template>
      <template #editor-description="{ data, field }">
        <Textarea v-model="data[field]" :id="`edit-${field}`" :name="field" rows="2" class="w-full" />
      </template>

      <template #editor-shifts="{ data, field }">
        <div class="p-2 border rounded-md bg-white dark:bg-gray-800 max-h-48 overflow-y-auto" @click.stop>
          <div v-if="isFetchingShifts" class="text-text-muted">Memuat shift...</div>
          <div v-else-if="!shifts.length" class="text-text-muted">Tidak ada shift.</div>
          <div v-else v-for="shift in shifts" :key="shift.id" class="flex items-center my-1">
            <Checkbox :inputId="`edit-shift-${data.id}-${shift.id}`" v-model="data.shift_ids" :value="shift.id" />
            <label :for="`edit-shift-${data.id}-${shift.id}`" class="ml-2 text-text-base">{{ shift.name }}</label>
          </div>
        </div>
      </template>

      <template #editor-locations="{ data, field }">
         <div class="p-2 border rounded-md bg-white dark:bg-gray-800 max-h-48 overflow-y-auto" @click.stop>
          <div v-if="isFetchingLocations" class="text-text-muted">Memuat lokasi...</div>
          <div v-else-if="!locations.length" class="text-text-muted">Tidak ada lokasi.</div>
          <div v-else v-for="loc in locations" :key="loc.id" class="flex items-center my-1">
            <Checkbox :inputId="`edit-loc-${data.id}-${loc.id}`" v-model="data.location_ids" :value="loc.id" />
            <label :for="`edit-loc-${data.id}-${loc.id}`" class="ml-2 text-text-base">{{ loc.name }}</label>
          </div>
        </div>
      </template>

      <template #column-shifts="{ item }">
        {{ item.Shifts ? item.Shifts.map(s => s.name).join(', ') : '-' }}
      </template>
      <template #column-locations="{ item }">
        {{ item.Locations ? item.Locations.map(l => l.Name).join(', ') : '-' }}
      </template>
      <template #column-actions="{ item }">
        <div class="flex items-center justify-center space-x-2">
          <BaseButton @click="openEditModal(item)" class="text-accent hover:opacity-80">
              <i class="pi pi-pencil"></i>
            </BaseButton>
            <BaseButton @click="deleteDivision(item.id)" class="text-danger hover:opacity-80">
              <i class="pi pi-trash"></i>
          </BaseButton>
        </div>
      </template>
    </BaseDataTable>

    <!-- Modal Tambah/Edit Divisi -->
    <BaseModal :isOpen="isModalOpen" @close="closeModal" :title="modalTitle" maxWidth="md">
      <form @submit.prevent="saveDivision">
        <BaseInput
          id="name"
          label="Nama Divisi:"
          v-model="currentDivision.name"
          required
        />
        
        <div class="mb-4">
          <label for="shifts" class="block text-text-muted text-sm font-bold mb-2">Pilih Shift:</label>
          <MultiSelect
            id="shifts"
            v-model="currentDivision.shift_ids"
            :options="shifts"
            optionLabel="name"
            optionValue="id"
            placeholder="Pilih Shift"
            display="chip"
            class="w-full"
            :loading="isFetchingShifts"
          />
        </div>
        <div class="mb-4">
          <label class="block text-text-muted text-sm font-bold mb-2">Pilih Lokasi Absen:</label>
          <div v-if="isFetchingLocations" class="text-text-muted">Memuat lokasi...</div>
          <div v-else-if="locations.length === 0" class="text-text-muted">
            <p class="mb-2">Tidak ada lokasi tersedia.</p>
            <BaseButton @click="goToAddLocation" class="btn-primary btn-sm">
              <i class="pi pi-plus"></i> Tambah Lokasi Absen
            </BaseButton>
          </div>
          <div v-else class="grid grid-cols-1 gap-2">
            <div v-for="loc in locations" :key="loc.id" class="flex items-center">
              <Checkbox :id="`loc-${loc.id}`" :value="loc.id" v-model="currentDivision.location_ids" />
              <label :for="`loc-${loc.id}`" class="ml-2 text-text-base">{{ loc.name }}</label>
            </div>
          </div>
        </div>
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
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';
import { useAuthStore } from '../../../stores/auth';
import { FilterMatchMode } from '@primevue/core/api';

import BaseDataTable from '../../ui/BaseDataTable.vue';
import BaseModal from '../../ui/BaseModal.vue';
import BaseInput from '../../ui/BaseInput.vue';
import BaseButton from '../../ui/BaseButton.vue';
import ConfirmDialog from 'primevue/confirmdialog';
import Textarea from 'primevue/textarea';
import MultiSelect from 'primevue/multiselect';
import Checkbox from 'primevue/checkbox';
import { useRouter } from 'vue-router';

const toast = useToast();
const confirm = useConfirm();
const authStore = useAuthStore();
const router = useRouter();

const divisions = ref([]);
const isLoading = ref(false);
const isModalOpen = ref(false);
const modalTitle = ref('');
const isEditMode = ref(false);
const currentDivision = ref({});
const isSaving = ref(false);
const isDeleting = ref(false);

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const shifts = ref([]);
const locations = ref([]);
const isFetchingShifts = ref(false);
const isFetchingLocations = ref(false);

const divisionColumns = ref([
  { field: 'name', header: 'Nama Divisi', editable: true },
  { field: 'description', header: 'Deskripsi', editable: true },
  { 
    field: 'shifts', 
    header: 'Shift', 
    editable: true, 
    body: (item) => item.Shifts ? item.Shifts.map(s => s.name).join(', ') : '-',
  },
  { 
    field: 'locations', 
    header: 'Lokasi Absen', 
    editable: true, 
    body: (item) => item.Locations ? item.Locations.map(l => l.Name).join(', ') : '-'
  },
  { field: 'actions', header: 'Aksi', editable: false }
]);

const fetchDivisions = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/admin/divisions');
    console.log('Raw divisions data from backend:', response.data.data);
    divisions.value = response.data.data.map(div => ({
          ...div,
          id: div.ID,
          name: div.Name, 
          description: div.Description,
          shifts: div.Shifts || [], // Add this line
          shift_ids: div.Shifts ? div.Shifts.map(s => s.id) : [],
          location_ids: div.Locations ? div.Locations.map(l => l.ID) : []
        }));
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.message || 'Gagal memuat divisi.', life: 3000 });
    console.error('Error fetching divisions:', err);
  } finally {
    isLoading.value = false;
  }
};

const fetchShifts = async () => {
  isFetchingShifts.value = true;
  try {
    const response = await axios.get('/api/shifts', {
      headers: { Authorization: `Bearer ${authStore.token}` },
    });
    console.log('Full shifts API response:', response);
    console.log('Raw shifts data from backend:', response.data.data);
    shifts.value = Array.isArray(response.data.data) ? response.data.data.map(s => ({
      id: s.id,
      name: s.name,
    })) : [];
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.message || 'Gagal memuat shift.', life: 3000 });
    console.error('Error fetching shifts:', err);
  } finally {
    isFetchingShifts.value = false;
  }
};

const fetchLocations = async () => {
  isFetchingLocations.value = true;
  try {
    const response = await axios.get('/api/company/locations', {
      headers: { Authorization: `Bearer ${authStore.token}` },
    });
    locations.value = Array.isArray(response.data.data) ? response.data.data : [];
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.message || 'Gagal memuat lokasi.', life: 3000 });
    console.error('Error fetching locations:', err);
    locations.value = []; // Ensure it's an empty array on error
  } finally {
    isFetchingLocations.value = false;
  }
};

const openAddModal = () => {
  isEditMode.value = false;
  modalTitle.value = 'Tambah Divisi Baru';
  currentDivision.value = { name: '', description: '', shift_ids: [], location_ids: [] };
  isModalOpen.value = true;
  fetchShifts();
  fetchLocations();
};

const openEditModal = (division) => {
  isEditMode.value = true;
  modalTitle.value = 'Edit Divisi';
  currentDivision.value = { 
    ...division, 
    shift_ids: division.Shifts ? division.Shifts.map(s => s.id) : [],
    location_ids: division.Locations ? division.Locations.map(l => l.ID) : [],
  };
  isModalOpen.value = true;
  fetchShifts();
  fetchLocations();
};

const closeModal = () => {
  isModalOpen.value = false;
  currentDivision.value = {};
};

const saveDivision = async () => {
  if (!currentDivision.value.name.trim()) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Nama divisi tidak boleh kosong.', life: 3000 });
    return;
  }

  isSaving.value = true;
  try {
    const payload = {
      name: currentDivision.value.name,
      description: currentDivision.value.description,
      shift_ids: currentDivision.value.shift_ids,
      location_ids: currentDivision.value.location_ids,
    };

    if (isEditMode.value) {
      await axios.put(`/api/admin/divisions/${currentDivision.value.id}`, payload);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Divisi berhasil diperbarui.', life: 3000 });
    } else {
      await axios.post('/api/admin/divisions', payload);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Divisi berhasil ditambahkan.', life: 3000 });
    }
    closeModal();
    fetchDivisions();
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.message || 'Gagal menyimpan divisi.', life: 3000 });
    console.error('Error saving division:', err);
  } finally {
    isSaving.value = false;
  }
};

const onCellEditComplete = async (event) => {
  console.log('[DEBUG] onCellEditComplete triggered', event);
  let { data, newData, field } = event;

  isSaving.value = true;
  try {
    let payload = {};

    if (field === 'shifts') {
      payload.shift_ids = Array.isArray(newData.shift_ids) ? newData.shift_ids : [];
    } else if (field === 'locations') {
      payload.location_ids = Array.isArray(newData.location_ids) ? newData.location_ids : [];
    } else {
      payload[field] = newData[field];
    }

    console.log(`[DEBUG] Saving division ${data.id} with payload:`, JSON.stringify(payload));

    await axios.put(`/api/admin/divisions/${data.id}`, payload);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Divisi berhasil diperbarui.', life: 3000 });
    fetchDivisions(); // Re-fetch all divisions to ensure data consistency

  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.message || 'Gagal memperbarui divisi.', life: 3000 });
    console.error('Error updating division:', err);
    
    // Jika gagal, muat ulang data untuk mengembalikan ke kondisi server yang sebenarnya.
    fetchDivisions(); 
  } finally {
    isSaving.value = false;
  }
};



const deleteDivision = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus divisi ini? Tindakan ini tidak dapat dibatalkan.',
    header: 'Konfirmasi Hapus Divisi',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      isDeleting.value = true;
      try {
        await axios.delete(`/api/admin/divisions/${id}`);
        toast.add({ severity: 'success', summary: 'Success', detail: 'Divisi berhasil dihapus.', life: 3000 });
        fetchDivisions();
      } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.message || 'Gagal menghapus divisi.', life: 3000 });
        console.error('Error deleting division:', err);
      } finally {
        isDeleting.value = false;
      }
    },
    reject: () => {
      toast.add({ severity: 'info', summary: 'Dibatalkan', detail: 'Penghapusan divisi dibatalkan', life: 3000 });
    }
  });
};

const goToAddLocation = () => {
  router.push({ name: 'LocationManagement' });
};

onMounted(() => {
  fetchDivisions();
  fetchShifts();
  fetchLocations();
});
</script>