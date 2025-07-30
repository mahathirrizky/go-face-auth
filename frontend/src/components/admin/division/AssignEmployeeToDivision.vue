<template>
  <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
    <Toast />
    <h3 class="text-xl font-semibold text-text-base mb-4">Tetapkan Karyawan ke Divisi</h3>

    <div class="mb-6">
       <FloatLabel variant="on">

         <Select
         id="divisionSelect"
         v-model="selectedDivisionId"
         :options="divisions"
         optionLabel="name"
         optionValue="id"
         fluid
         class="w-full"
         :loading="isFetchingDivisions"
         />
         <label for="divisionSelect" class="block text-text-muted text-sm font-bold mb-2">Pilih Divisi:</label>
        </FloatLabel>
    </div>

    <div v-if="isLoading" class="text-center p-8">
      <ProgressSpinner />
      <p class="text-text-muted mt-2">Memuat data karyawan...</p>
    </div>

    <div v-else-if="selectedDivisionId">
      <PickList
        v-model="employeeLists"
        listStyle="height: calc(100vh - 300px); overflow-y: auto;"
        dataKey="id"
        @move-to-target="onMoveToOther"
        @move-to-source="onMoveToDivision"
        @move-all-to-target="onMoveAllToOther"
        @move-all-to-source="onMoveAllToDivision"
      >
        <template #sourceheader>
          Karyawan di Divisi ({{ getDivisionName(selectedDivisionId) }})
        </template>
        <template #targetheader>
          <div class="flex flex-col">
            <span>Karyawan Lain</span>
            <IconField iconPosition="left" class="mt-2">
              <InputIcon class="pi pi-search"></InputIcon>
              <InputText v-model="searchTerm" placeholder="Cari Karyawan..." @input="onSearch" class="w-full" fluid/>
            </IconField>
          </div>
        </template>
        <template #item="slotProps">
          <div class="flex items-center p-2">
            <div>{{ slotProps.item.name }} ({{ slotProps.item.position }})</div>
          </div>
        </template>
      </PickList>
    </div>

    <div v-else class="text-center p-8 border-2 border-dashed rounded-lg">
      <i class="pi pi-info-circle text-4xl text-text-muted mb-2"></i>
      <p class="text-text-muted">Silakan pilih divisi untuk mengelola karyawan.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../../stores/auth';
import Select from 'primevue/select';
import Toast from 'primevue/toast';
import PickList from 'primevue/picklist';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import ProgressSpinner from 'primevue/progressspinner';
import debounce from 'lodash.debounce';
import FloatLabel from 'primevue/floatlabel';

const toast = useToast();
const authStore = useAuthStore();

const divisions = ref([]);
const selectedDivisionId = ref(null);
const isFetchingDivisions = ref(false);

const employeeLists = ref([[], []]);
const allOtherEmployees = ref([]);
const isLoading = ref(false);
const searchTerm = ref('');

const fetchDivisions = async () => {
  isFetchingDivisions.value = true;
  try {
    const response = await axios.get('/api/admin/divisions');
    divisions.value = response.data.data.map(div => ({ id: div.ID, name: div.Name }));
    if (divisions.value.length > 0) {
      selectedDivisionId.value = divisions.value[0].id;
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat divisi.', life: 3000 });
  } finally {
    isFetchingDivisions.value = false;
  }
};

const fetchEmployees = async () => {
  if (!authStore.companyId || !selectedDivisionId.value) {
    employeeLists.value = [[], []];
    return;
  }
  isLoading.value = true;
  try {
    const [inDivisionRes, noDivisionRes] = await Promise.all([
      axios.get(`/api/companies/${authStore.companyId}/employees`, { params: { division_id: selectedDivisionId.value, limit: 1000 } }),
      axios.get(`/api/companies/${authStore.companyId}/employees`, { params: { no_division: true, limit: 1000 } })
    ]);

    const employeesInDivision = inDivisionRes.data.data.items.map(emp => ({ id: emp.ID, name: emp.Name, email: emp.Email, position: emp.Position }));
    const otherEmployees = noDivisionRes.data.data.items.map(emp => ({ id: emp.ID, name: emp.Name, email: emp.Email, position: emp.Position }));

    allOtherEmployees.value = otherEmployees;
    employeeLists.value = [employeesInDivision, filterOtherEmployees(allOtherEmployees.value, searchTerm.value)];

  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengambil data karyawan.', life: 3000 });
    employeeLists.value = [[], []];
    allOtherEmployees.value = [];
  } finally {
    isLoading.value = false;
  }
};

const handleEmployeeMove = async (items, targetDivisionId) => {
  const updates = items.map(item => ({ employee_id: item.id, division_id: targetDivisionId }));
  try {
    await axios.put(`/api/admin/employees/division/batch`, { updates });
    toast.add({ severity: 'success', summary: 'Sukses', detail: 'Divisi karyawan berhasil diperbarui.', life: 3000 });
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memperbarui divisi karyawan.', life: 3000 });
    fetchEmployees(); // Revert on failure
  }
};

const onMoveToOther = (event) => handleEmployeeMove(event.items, null);
const onMoveToDivision = (event) => handleEmployeeMove(event.items, selectedDivisionId.value);
const onMoveAllToOther = (event) => handleEmployeeMove(event.items, null);
const onMoveAllToDivision = (event) => handleEmployeeMove(event.items, selectedDivisionId.value);

const getDivisionName = (divisionId) => {
  const division = divisions.value.find(d => d.id === divisionId);
  return division ? division.name : 'Tidak Diketahui';
};

const filterOtherEmployees = (employees, term) => {
  if (!term) return employees;
  const lowerCaseTerm = term.toLowerCase();
  return employees.filter(emp => 
    emp.name.toLowerCase().includes(lowerCaseTerm) || 
    (emp.email && emp.email.toLowerCase().includes(lowerCaseTerm)) || 
    (emp.position && emp.position.toLowerCase().includes(lowerCaseTerm))
  );
};

const onSearch = debounce(() => {
  employeeLists.value[1] = filterOtherEmployees(allOtherEmployees.value, searchTerm.value);
}, 300);

watch(selectedDivisionId, (newVal) => {
  if (newVal) {
    fetchEmployees();
  }
});

onMounted(() => {
  fetchDivisions();
});
</script>
