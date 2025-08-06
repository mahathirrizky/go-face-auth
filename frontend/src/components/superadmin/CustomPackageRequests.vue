<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <h2 class="text-2xl font-bold text-text-base mb-6">Permintaan Paket Kustom</h2>

    <DataTable
      :value="customRequests"
      :loading="isLoading"
      :totalRecords="totalRecords"
      :lazy="true"
      v-model:filters="filters"
      @page="onPage"
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      dataKey="id"
      :globalFilterFields="['company_name', 'name', 'email', 'phone', 'message']"
    >
      <template #header>
        <div class="flex justify-end">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Permintaan..." @keydown.enter="onFilter" fluid />
          </IconField>
        </div>
      </template>

      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="company_name" header="Perusahaan" :sortable="true"></Column>
      <Column field="name" header="Nama Kontak" :sortable="true"></Column>
      <Column field="email" header="Email Kontak" :sortable="true"></Column>
      <Column field="phone" header="Telepon Kontak" :sortable="true"></Column>
      <Column field="message" header="Pesan" :sortable="true"></Column>
      <Column field="status" header="Status" :sortable="true">
        <template #body="{ data }">
          <Tag :value="data.status" :severity="getStatusSeverity(data.status)" />
        </template>
      </Column>
      <Column field="created_at" header="Tanggal Permintaan" :sortable="true">
        <template #body="{ data }">
          {{ new Date(data.created_at).toLocaleString('id-ID') }}
        </template>
      </Column>
      <Column header="Aksi">
        <template #body="{ data }">
          <div class="flex space-x-2">
            <Button 
              @click="openCreateOfferModal(data)" 
              class="p-button-primary p-button-sm" 
              icon="pi pi-file-edit" 
              label="Buat Penawaran" 
              :disabled="data.status !== 'pending'"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="showOfferModal" header="Buat Penawaran Kustom" :modal="true" class="w-full max-w-2xl">
      <form @submit.prevent="submitCustomOffer" class="p-fluid grid grid-cols-1 md:grid-cols-2 gap-6 pt-4">
        <div class="col-span-2">
          <FloatLabel variant="on">
            <InputText id="companyName" v-model="currentOffer.company_name" disabled />
            <label for="companyName">Perusahaan</label>
          </FloatLabel>
        </div>
        <div class="col-span-2">
          <FloatLabel variant="on">
            <InputText id="packageName" v-model="currentOffer.package_name" required />
            <label for="packageName">Nama Paket Kustom</label>
          </FloatLabel>
        </div>
        <div class="field">
          <FloatLabel variant="on">
            <InputNumber id="maxEmployees" v-model="currentOffer.max_employees" required />
            <label for="maxEmployees">Maksimal Karyawan</label>
          </FloatLabel>
        </div>
        <div class="field">
          <FloatLabel variant="on">
            <InputNumber id="maxLocations" v-model="currentOffer.max_locations" required />
            <label for="maxLocations">Maksimal Lokasi</label>
          </FloatLabel>
        </div>
        <div class="field">
          <FloatLabel variant="on">
            <InputNumber id="maxShifts" v-model="currentOffer.max_shifts" required />
            <label for="maxShifts">Maksimal Shift</label>
          </FloatLabel>
        </div>
        <div class="field">
           <FloatLabel variant="on">
            <Dropdown id="billingCycle" v-model="currentOffer.billing_cycle" :options="billingCycles" optionLabel="name" optionValue="value" required />
            <label for="billingCycle">Siklus Penagihan</label>
          </FloatLabel>
        </div>
        <div class="col-span-2">
          <FloatLabel variant="on">
            <Textarea id="features" v-model="currentOffer.features" rows="4" required autoResize />
            <label for="features">Fitur (pisahkan dengan koma)</label>
          </FloatLabel>
        </div>
        <div class="col-span-2">
          <FloatLabel variant="on">
            <InputNumber id="finalPrice" v-model="currentOffer.final_price" mode="currency" currency="IDR" locale="id-ID" required />
            <label for="finalPrice">Harga Total (IDR)</label>
          </FloatLabel>
        </div>
        <div class="col-span-2 flex justify-end space-x-2 mt-6">
          <Button type="button" @click="showOfferModal = false" label="Batal" class="p-button-outlined" />
          <Button type="submit" :loading="isSubmittingOffer" label="Kirim Penawaran & Invoice" class="p-button-primary" />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useWebSocketStore } from '../../stores/websocket';
import { FilterMatchMode } from '@primevue/core/api';

import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Tag from 'primevue/tag';
import Dialog from 'primevue/dialog';
import Textarea from 'primevue/textarea';
import InputNumber from 'primevue/inputnumber';
import Dropdown from 'primevue/dropdown';
import Toast from 'primevue/toast';
import FloatLabel from 'primevue/floatlabel';

const customRequests = ref([]);
const isLoading = ref(false);
const totalRecords = ref(0);
const lazyParams = ref({});
const toast = useToast();
const webSocketStore = useWebSocketStore();

const showOfferModal = ref(false);
const isSubmittingOffer = ref(false);
const currentOffer = ref({});

const filters = ref({
  'global': { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const billingCycles = ref([
  { name: 'Bulanan', value: 'monthly' },
  { name: 'Tahunan', value: 'yearly' },
]);

const fetchCustomRequests = async () => {
  isLoading.value = true;
  try {
    const params = {
      page: lazyParams.value.page + 1,
      limit: lazyParams.value.rows,
      search: filters.value.global.value || '',
    };
    const response = await axios.get('/api/superadmin/custom-package-requests', { params });
    if (response.data && response.data.status === 'success') {
      customRequests.value = response.data.data.items;
      totalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengambil permintaan kustom.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching custom requests:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengambil permintaan kustom.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const onPage = (event) => {
  lazyParams.value = event;
  fetchCustomRequests();
};

const onFilter = () => {
  lazyParams.value.page = 0; // Reset to first page on filter
  fetchCustomRequests();
};

const openCreateOfferModal = (request) => {
  currentOffer.value = {
    company_id: request.company_id,
    company_name: request.company_name,
    package_name: '',
    max_employees: 10,
    max_locations: 1,
    max_shifts: 1,
    features: '',
    final_price: 0,
    billing_cycle: 'monthly',
    request_id: request.id, // Keep track of the request ID
  };
  showOfferModal.value = true;
};

const submitCustomOffer = async () => {
  isSubmittingOffer.value = true;
  try {
    // Step 1: Create the custom offer
    const offerResponse = await axios.post('/api/superadmin/custom-offers', currentOffer.value);
    
    if (offerResponse.data && offerResponse.data.status === 'success') {
      toast.add({ severity: 'info', summary: 'Berhasil', detail: 'Penawaran kustom berhasil dibuat. Mengirim invoice...', life: 3000 });

      // Step 2: Update the request status to 'contacted' (or a new status like 'offered')
      await axios.put(`/api/superadmin/custom-package-requests/${currentOffer.value.request_id}/contacted`);
      
      toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Invoice telah dikirim ke email perusahaan.', life: 5000 });
      
      showOfferModal.value = false;
      fetchCustomRequests(); // Refresh the list
    } else {
      toast.add({ severity: 'error', summary: 'Gagal', detail: offerResponse.data?.message || 'Gagal membuat penawaran kustom.', life: 3000 });
    }
  } catch (error) {
    console.error('Error submitting custom offer:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan fatal saat mengirim penawaran.', life: 3000 });
  } finally {
    isSubmittingOffer.value = false;
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'pending':
      return 'warning';
    case 'contacted':
      return 'info';
    case 'resolved':
      return 'success';
    default:
      return null;
  }
};

const handleWebSocketMessage = (data) => {
  if (data.type === 'new_custom_package_request') {
    toast.add({ severity: 'info', summary: 'Permintaan Baru!', detail: `Permintaan paket kustom baru dari ${data.company_name}.`, life: 5000 });
    fetchCustomRequests();
  }
};

onMounted(() => {
  lazyParams.value = {
    first: 0,
    rows: 10,
    page: 0,
  };
  fetchCustomRequests();
  webSocketStore.onMessage('superadmin_notification', handleWebSocketMessage);
});

onUnmounted(() => {
  webSocketStore.offMessage('superadmin_notification');
});
</script>
