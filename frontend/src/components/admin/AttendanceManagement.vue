<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Absensi</h2>

    <Tabs v-model:value="selectedTab">
      <TabList>
        <Tab :value="0">Semua Absensi</Tab>
        <Tab :value="1">Karyawan Tidak Absen</Tab>
        <Tab :value="2">Lembur</Tab>
      </TabList>

      <TabPanels>
        <TabPanel :value="0">
          <BaseDataTable
            :data="attendanceRecords"
            :columns="attendanceColumns"
            :loading="isLoading"
            :totalRecords="attendancesTotalRecords"
            :lazy="true"
            v-model:filters="attendancesFilters"
            @page="onPage($event, 'attendances')"
            @filter="onFilter($event, 'attendances')"
            searchPlaceholder="Cari Absensi..."
          >
            <template #header-actions>
              <div class="flex flex-wrap items-center gap-2">
                <div class="flex items-center">
                  <label for="startDate" class="text-text-muted mr-2">Dari:</label>
                  <BaseInput
                    type="date"
                    id="startDate"
                    v-model="startDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <div class="flex items-center">
                  <label for="endDate" class="text-text-muted mr-2">Sampai:</label>
                  <BaseInput
                    type="date"
                    id="endDate"
                    v-model="endDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <BaseButton @click="fetchAttendances" class="btn-primary"><i class="pi pi-filter"></i> Filter</BaseButton>
                <BaseButton @click="exportAllToExcel" class="btn-secondary whitespace-nowrap"><i class="pi pi-file-excel"></i> Export</BaseButton>
                <BaseButton @click="openCorrectionModal()" class="btn-accent whitespace-nowrap"><i class="pi pi-plus"></i> Tambah Koreksi</BaseButton>
              </div>
            </template>
            <template #column-date="{ item }">
              {{ new Date(item.check_in_time).toLocaleDateString() }}
            </template>

            <template #column-check_in_time="{ item }">
              {{ new Date(item.check_in_time).toLocaleTimeString() }}
            </template>

            <template #column-check_out_time="{ item }">
              {{ item.check_out_time ? new Date(item.check_out_time).toLocaleTimeString() : '-' }}
            </template>

            <template #column-status="{ item }">
              <span :class="{
                'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
                'bg-green-100 text-green-600': item.status === 'on_time',
                'bg-yellow-100 text-yellow-600': item.status === 'late',
                'bg-blue-100 text-blue-600': item.status === 'overtime_in' || item.status === 'overtime_out',
                'bg-purple-100 text-purple-600': item.is_correction,
              }">
                {{ item.is_correction ? 'Dikoreksi' : (item.status === 'on_time' ? 'Tepat Waktu' : item.status === 'late' ? 'Terlambat' : item.status === 'overtime_in' ? 'Lembur Masuk' : item.status === 'overtime_out' ? 'Lembur Keluar' : item.status) }}
              </span>
            </template>

            <template #column-employee.name="{ item }">
              {{ item.employee ? item.employee.name : 'N/A' }}
            </template>
            
            <template #column-actions="{ item }">
              <BaseButton 
                v-if="!item.check_out_time" 
                @click="openCorrectionModal(item)" 
                class="btn-info btn-sm"
                v-tooltip.top="'Tambah Waktu Pulang'">
                <i class="pi pi-pencil"></i>
              </BaseButton>
            </template>
          </BaseDataTable>
        </TabPanel>

        <!-- Other TabPanels remain the same -->

      </TabPanels>
    </Tabs>

    <!-- Correction Modal -->
    <BaseModal :isOpen="isCorrectionModalOpen" @close="closeCorrectionModal" title="Koreksi Absensi">
      <form @submit.prevent="submitCorrection">
        <div class="mb-4">
          <label for="employee" class="block text-sm font-medium text-text-muted mb-2">Karyawan</label>
          <Dropdown 
            v-model="correctionForm.employee_id" 
            :options="allEmployees" 
            optionLabel="name" 
            optionValue="id" 
            placeholder="Pilih Karyawan" 
            class="w-full"
            :disabled="!!correctionForm.id"
            filter
          />
        </div>

        <div class="mb-4">
          <label for="correction_time" class="block text-sm font-medium text-text-muted mb-2">Tanggal & Waktu Koreksi</label>
          <Calendar v-model="correctionForm.correction_time" showTime hourFormat="24" class="w-full" />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-muted mb-2">Tipe Koreksi</label>
          <div class="flex items-center gap-4">
            <div class="flex items-center">
              <RadioButton v-model="correctionForm.correction_type" inputId="type_check_in" name="correction_type" value="check_in" :disabled="!!correctionForm.id" />
              <label for="type_check_in" class="ml-2">Check-in</label>
            </div>
            <div class="flex items-center">
              <RadioButton v-model="correctionForm.correction_type" inputId="type_check_out" name="correction_type" value="check_out" :disabled="!correctionForm.id && correctionForm.correction_type !== 'check_out'" />
              <label for="type_check_out" class="ml-2">Check-out</label>
            </div>
          </div>
        </div>

        <div class="mb-4">
          <label for="notes" class="block text-sm font-medium text-text-muted mb-2">Alasan Koreksi</label>
          <Textarea v-model="correctionForm.notes" rows="3" class="w-full" placeholder="Contoh: Karyawan lupa absen pulang." />
        </div>

        <div class="flex justify-end space-x-4 mt-6">
          <BaseButton @click="closeCorrectionModal" type="button" class="btn-outline-primary">Batal</BaseButton>
          <BaseButton type="submit" :loading="isSubmitting">Simpan</BaseButton>
        </div>
      </form>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';
import BaseModal from '../ui/BaseModal.vue';
import Tabs from 'primevue/tabs';
import Tab from 'primevue/tab';
import TabList from 'primevue/tablist';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import Dropdown from 'primevue/dropdown';
import Calendar from 'primevue/calendar';
import RadioButton from 'primevue/radiobutton';
import Textarea from 'primevue/textarea';
import Tooltip from 'primevue/tooltip';
import { FilterMatchMode } from '@primevue/core/api';

const attendanceRecords = ref([]);
const unaccountedEmployees = ref([]);
const overtimeRecords = ref([]);
const selectedTab = ref(0);
const toast = useToast();
const authStore = useAuthStore();
const isLoading = ref(false);
const isSubmitting = ref(false);

// State for Correction Modal
const isCorrectionModalOpen = ref(false);
const allEmployees = ref([]);
const correctionForm = ref({
  id: null,
  employee_id: null,
  correction_time: null,
  correction_type: 'check_in',
  notes: ''
});

// State for All Attendances table
const attendancesTotalRecords = ref(0);
const attendancesLazyParams = ref({});
const attendancesFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

// ... (other table states remain the same)

const attendanceColumns = ref([
    { field: 'date', header: 'Tanggal' },
    { field: 'employee.name', header: 'Nama Karyawan' },
    { field: 'check_in_time', header: 'Waktu Masuk' },
    { field: 'check_out_time', header: 'Waktu Keluar' },
    { field: 'status', header: 'Status' },
    { field: 'actions', header: 'Aksi' }
]);

// ... (other column definitions remain the same)

const fetchAllEmployees = async () => {
  if (!authStore.companyId) return;
  try {
    const response = await axios.get(`/api/companies/${authStore.companyId}/employees`);
    allEmployees.value = response.data.employees;
  } catch (error) {
    console.error('Error fetching employees:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat daftar karyawan.', life: 3000 });
  }
};

const openCorrectionModal = (item = null) => {
  if (item) {
    // Editing existing record (adding check-out)
    correctionForm.value = {
      id: item.id,
      employee_id: item.employee_id,
      correction_time: new Date(),
      correction_type: 'check_out',
      notes: ''
    };
  } else {
    // Creating new record
    correctionForm.value = {
      id: null,
      employee_id: null,
      correction_time: new Date(),
      correction_type: 'check_in',
      notes: ''
    };
  }
  isCorrectionModalOpen.value = true;
};

const closeCorrectionModal = () => {
  isCorrectionModalOpen.value = false;
};

const submitCorrection = async () => {
  isSubmitting.value = true;
  try {
    const payload = {
      employee_id: correctionForm.value.employee_id,
      correction_time: correctionForm.value.correction_time.toISOString(),
      correction_type: correctionForm.value.correction_type,
      notes: correctionForm.value.notes
    };
    
    await axios.post('/api/attendances/correction', payload);
    
    toast.add({ severity: 'success', summary: 'Sukses', detail: 'Absensi berhasil dikoreksi.', life: 3000 });
    closeCorrectionModal();
    fetchAttendances(); // Refresh the table
  } catch (error) {
    console.error('Error submitting correction:', error);
    const message = error.response?.data?.message || 'Gagal menyimpan koreksi.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isSubmitting.value = false;
  }
};

// ... (all other existing functions like fetchAttendances, onPage, etc. remain the same)

onMounted(() => {
  // ... (existing onMounted logic)
  fetchAllEmployees(); // Fetch employees for the dropdown
});

// ... (existing watch logic)
</script>

<style scoped>
/* Tailwind handles styling */
</style>
