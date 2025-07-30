<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Perusahaan</h2>

    <DataTable
      :value="companies"
      :loading="loading"
      v-model:filters="filters"
      :globalFilterFields="['name', 'address', 'subscription_status']"
      paginator
      :rows="10"
      class="p-datatable-customers"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      dataKey="id"
    >
      <template #header>
        <div class="flex flex-wrap justify-between items-center gap-4">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Perusahaan..." fluid />
          </IconField>
        </div>
      </template>
      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="id" header="ID" :sortable="true"></Column>
      <Column field="name" header="Nama Perusahaan" :sortable="true"></Column>
      <Column field="address" header="Alamat" :sortable="true"></Column>
      <Column field="subscription_status" header="Status Langganan" :sortable="true"></Column>
      <Column field="subscription_package.package_name" header="Paket Langganan" :sortable="true">
         <template #body="{ data }">
            {{ data.subscription_package ? data.subscription_package.package_name : 'N/A' }}
        </template>
      </Column>
      <Column field="billing_cycle" header="Siklus Penagihan" :sortable="true"></Column>
      <Column field="created_at" header="Tanggal Dibuat" :sortable="true">
        <template #body="{ data }">
          {{ new Date(data.created_at).toLocaleDateString() }}
        </template>
      </Column>
      <Column header="Aksi" style="min-width: 12rem">
        <template #body="{ data }">
          <Button @click="openCreateOfferModal(data)" class="p-button-primary p-button-sm" icon="pi pi-plus" label="Buat Penawaran" />
        </template>
      </Column>
    </DataTable>

    <!-- Create Custom Offer Modal -->
    <Dialog v-model:visible="showCreateOfferModal" header="Buat Penawaran Kustom" :modal="true" class="w-full max-w-lg">
      <form @submit.prevent="createCustomOffer" class="p-fluid">
        <div class="mb-4">
          <label class="block text-text-muted text-sm font-bold mb-2">Perusahaan:</label>
          <p class="text-text-base font-semibold">{{ newOffer.company_name }}</p>
        </div>

        <div class="field mb-4">
          <label for="packageName" class="block text-text-muted text-sm font-bold mb-2">Nama Paket:</label>
          <InputText id="packageName" v-model="newOffer.package_name" required fluid />
        </div>

        <div class="field mb-4">
          <label for="maxEmployees" class="block text-text-muted text-sm font-bold mb-2">Jumlah Maksimal Karyawan:</label>
          <InputNumber id="maxEmployees" v-model="newOffer.max_employees" :min="1" required fluid />
        </div>

        <div class="field mb-4">
          <label for="maxLocations" class="block text-text-muted text-sm font-bold mb-2">Jumlah Maksimal Lokasi:</label>
          <InputNumber id="maxLocations" v-model="newOffer.max_locations" :min="0" required fluid />
        </div>

        <div class="field mb-4">
          <label for="maxShifts" class="block text-text-muted text-sm font-bold mb-2">Jumlah Maksimal Shift:</label>
          <InputNumber id="maxShifts" v-model="newOffer.max_shifts" :min="0" required fluid />
        </div>

        <div class="field mb-4">
          <label for="features" class="block text-text-muted text-sm font-bold mb-2">Fitur (pisahkan dengan koma):</label>
          <Textarea id="features" v-model="newOffer.features" rows="3" required fluid />
        </div>

        <div class="field mb-4">
          <label for="finalPrice" class="block text-text-muted text-sm font-bold mb-2">Harga Akhir (IDR):</label>
          <InputNumber id="finalPrice" v-model="newOffer.final_price" mode="currency" currency="IDR" locale="id-ID" :min="0" required fluid />
        </div>

        <div class="field mb-4">
          <label for="billingCycle" class="block text-text-muted text-sm font-bold mb-2">Siklus Penagihan:</label>
          <Dropdown id="billingCycle" v-model="newOffer.billing_cycle" :options="billingCycleOptions" optionLabel="label" optionValue="value" placeholder="Pilih Siklus" required fluid />
        </div>

        <div v-if="generatedOfferLink" class="mt-6 p-4 bg-green-100 text-green-800 rounded-md break-all">
          <p class="font-semibold mb-2">Tautan Penawaran Berhasil Dibuat:</p>
          <a :href="generatedOfferLink" target="_blank" class="text-blue-600 hover:underline">{{ generatedOfferLink }}</a>
          <Button @click="copyLink" type="button" class="p-button-text p-button-sm" icon="pi pi-copy" label="Salin" />
        </div>

        <div class="flex justify-end space-x-4 mt-6">
          <Button type="button" @click="showCreateOfferModal = false" class="p-button-outlined" label="Batal" />
          <Button type="submit" :loading="isSubmittingOffer" :disabled="generatedOfferLink !== ''" icon="pi pi-send" label="Buat Tautan" />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { FilterMatchMode } from '@primevue/core/api';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Textarea from 'primevue/textarea';
import Select from 'primevue/select';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';

const companies = ref([]);
const loading = ref(true);
const error = ref(null);
const toast = useToast();

const showCreateOfferModal = ref(false);
const selectedCompany = ref(null);
const newOffer = ref({
  company_id: null,
  company_name: '',
  package_name: '',
  max_employees: 1,
  max_locations: 0,
  max_shifts: 0,
  features: '',
  final_price: 0,
  billing_cycle: 'monthly',
});
const generatedOfferLink = ref('');
const isSubmittingOffer = ref(false);

const billingCycleOptions = [
  { label: 'Bulanan', value: 'monthly' },
  { label: 'Tahunan', value: 'yearly' },
];

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const fetchCompanies = async () => {
  try {
    const response = await axios.get('/api/superadmin/companies');
    if (response.data && response.data.status === 'success') {
      companies.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Failed to fetch companies.', life: 3000 });
      error.value = response.data.message || 'Failed to fetch companies.';
    }
  } catch (err) {
    console.error('Error fetching companies:', err);
    toast.add({ severity: 'error', summary: 'Error', detail: 'An error occurred while fetching companies.', life: 3000 });
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchCompanies();
});

const openCreateOfferModal = (company) => {
  selectedCompany.value = company;
  newOffer.value = {
    company_id: company.id,
    company_name: company.name,
    package_name: '',
    max_employees: 1,
    max_locations: 0,
    max_shifts: 0,
    features: '',
    final_price: 0,
    billing_cycle: 'monthly',
  };
  generatedOfferLink.value = ''; // Reset generated link
  showCreateOfferModal.value = true;
};

const createCustomOffer = async () => {
  isSubmittingOffer.value = true;
  try {
    const response = await axios.post('/api/superadmin/custom-offers', newOffer.value);
    if (response.data && response.data.status === 'success') {
      generatedOfferLink.value = response.data.data.link;
      toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Tautan penawaran kustom berhasil dibuat!', life: 5000 });
    } else {
      toast.add({ severity: 'error', summary: 'Gagal', detail: response.data?.message || 'Gagal membuat penawaran kustom.', life: 3000 });
    }
  } catch (error) {
    console.error('Error creating custom offer:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Terjadi kesalahan saat membuat penawaran kustom.', life: 3000 });
  } finally {
    isSubmittingOffer.value = false;
  }
};

const copyLink = () => {
  navigator.clipboard.writeText(generatedOfferLink.value);
  toast.add({ severity: 'info', summary: 'Disalin', detail: 'Tautan penawaran telah disalin ke clipboard.', life: 3000 });
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
