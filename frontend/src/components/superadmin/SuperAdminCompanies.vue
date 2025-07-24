<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Perusahaan</h2>

    <BaseDataTable
      :data="companies"
      :columns="companyColumns"
      :loading="loading"
      v-model:filters="filters"
      :globalFilterFields="['name', 'address', 'subscription_status']"
      searchPlaceholder="Cari Perusahaan..."
    >
      <template #column-subscription_package="{ item }">
        {{ item.subscription_package ? item.subscription_package.package_name : 'N/A' }}
      </template>

      <template #column-created_at="{ item }">
        {{ new Date(item.created_at).toLocaleDateString() }}
      </template>

      <template #column-actions="{ item }">
        <BaseButton @click="openCreateOfferModal(item)" class="btn-primary btn-sm">
          <i class="pi pi-plus"></i> Buat Penawaran Kustom
        </BaseButton>
      </template>
    </BaseDataTable>

    <!-- Create Custom Offer Modal -->
    <BaseModal :isOpen="showCreateOfferModal" @close="showCreateOfferModal = false" title="Buat Penawaran Kustom" maxWidth="lg">
      <form @submit.prevent="createCustomOffer">
        <div class="mb-4">
          <label class="block text-text-muted text-sm font-bold mb-2">Perusahaan:</label>
          <p class="text-text-base font-semibold">{{ newOffer.company_name }}</p>
        </div>

        <div class="mb-4">
          <label for="packageName" class="block text-text-muted text-sm font-bold mb-2">Nama Paket:</label>
          <BaseInput id="packageName" v-model="newOffer.package_name" required />
        </div>

        <div class="mb-4">
          <label for="maxEmployees" class="block text-text-muted text-sm font-bold mb-2">Jumlah Maksimal Karyawan:</label>
          <InputNumber id="maxEmployees" v-model="newOffer.max_employees" :min="1" required class="w-full" />
        </div>

        <div class="mb-4">
          <label for="maxLocations" class="block text-text-muted text-sm font-bold mb-2">Jumlah Maksimal Lokasi:</label>
          <InputNumber id="maxLocations" v-model="newOffer.max_locations" :min="0" required class="w-full" />
        </div>

        <div class="mb-4">
          <label for="maxShifts" class="block text-text-muted text-sm font-bold mb-2">Jumlah Maksimal Shift:</label>
          <InputNumber id="maxShifts" v-model="newOffer.max_shifts" :min="0" required class="w-full" />
        </div>

        <div class="mb-4">
          <label for="features" class="block text-text-muted text-sm font-bold mb-2">Fitur (pisahkan dengan koma):</label>
          <Textarea id="features" v-model="newOffer.features" rows="3" required class="w-full" />
        </div>

        <div class="mb-4">
          <label for="finalPrice" class="block text-text-muted text-sm font-bold mb-2">Harga Akhir (IDR):</label>
          <InputNumber id="finalPrice" v-model="newOffer.final_price" mode="currency" currency="IDR" locale="id-ID" :min="0" required class="w-full" />
        </div>

        <div class="mb-4">
          <label for="billingCycle" class="block text-text-muted text-sm font-bold mb-2">Siklus Penagihan:</label>
          <Dropdown id="billingCycle" v-model="newOffer.billing_cycle" :options="billingCycleOptions" optionLabel="label" optionValue="value" placeholder="Pilih Siklus" required class="w-full" />
        </div>

        <div v-if="generatedOfferLink" class="mt-6 p-4 bg-green-100 text-green-800 rounded-md break-all">
          <p class="font-semibold mb-2">Tautan Penawaran Berhasil Dibuat:</p>
          <a :href="generatedOfferLink" target="_blank" class="text-blue-600 hover:underline">{{ generatedOfferLink }}</a>
          <BaseButton @click="copyLink" type="button" class="btn-outline-success btn-sm ml-4">
            <i class="pi pi-copy"></i> Salin Tautan
          </BaseButton>
        </div>

        <div class="flex justify-end space-x-4 mt-6">
          <BaseButton type="button" @click="showCreateOfferModal = false" class="btn-outline-primary">
            Batal
          </BaseButton>
          <BaseButton type="submit" :loading="isSubmittingOffer" :disabled="generatedOfferLink !== ''">
            <i class="pi pi-send"></i> Buat Tautan Penawaran
          </BaseButton>
        </div>
      </form>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { FilterMatchMode } from '@primevue/core/api';
import BaseDataTable from '../ui/BaseDataTable.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseModal from '../ui/BaseModal.vue';
import BaseInput from '../ui/BaseInput.vue';
import InputNumber from 'primevue/inputnumber';
import Textarea from 'primevue/textarea';
import Dropdown from 'primevue/dropdown';

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

const companyColumns = ref([
    { field: 'id', header: 'ID' },
    { field: 'name', header: 'Nama Perusahaan' },
    { field: 'address', header: 'Alamat' },
    { field: 'subscription_status', header: 'Status Langganan' },
    { field: 'subscription_package', header: 'Paket Langganan' },
    { field: 'billing_cycle', header: 'Siklus Penagihan' },
    { field: 'created_at', header: 'Tanggal Dibuat' },
    { field: 'actions', header: 'Aksi', sortable: false }
]);

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