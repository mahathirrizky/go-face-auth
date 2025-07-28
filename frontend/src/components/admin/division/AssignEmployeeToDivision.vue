<template>
  <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
    <h3 class="text-xl font-semibold text-text-base mb-4">Tetapkan Karyawan ke Divisi</h3>

    <!-- Division Selection -->
    <div class="mb-6">
      <label for="divisionSelect" class="block text-text-muted text-sm font-bold mb-2">Pilih Divisi:</label>
      <Select
        id="divisionSelect"
        v-model="selectedDivisionId"
        :options="divisions"
        optionLabel="name"
        optionValue="id"
        placeholder="Pilih Divisi"
        class="w-full"
        :loading="isFetchingDivisions"
      />
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center p-8">
      <i class="pi pi-spin pi-spinner text-4xl text-primary"></i>
      <p class="text-text-muted mt-2">Memuat data karyawan...</p>
    </div>

    <!-- PickList for Employee Assignment -->
    <div v-else-if="selectedDivisionId">
      <PickList
        v-model="employeeLists"
        listStyle="height: calc(100vh - 250px); overflow-y: auto;"
        dataKey="id"
        @update:modelValue="onPickListMove"
      >
        <template #sourceheader>
          Karyawan di Divisi ({{ getDivisionName(selectedDivisionId) }})
        </template>
        <template #targetheader>
          <div class="flex flex-col">
            <span>Karyawan Lain</span>
            <span class="relative mt-2">
              <i class="pi pi-search absolute top-1/2 -translate-y-1/2 left-3 text-text-muted"></i>
              <InputText 
                v-model="searchTerm" 
                placeholder="Cari Karyawan..." 
                @input="onSearch"
                class="w-full pl-10"
              />
            </span>
          </div>
        </template>
        <template #item="slotProps">
          <div class="flex items-center p-2">
            <div>{{ slotProps.item.name }}</div>
          </div>
        </template>
      </PickList>
    </div>

    <!-- No Division Selected State -->
    <div v-else class="text-center p-8 border-2 border-dashed rounded-lg">
      <i class="pi pi-info-circle text-4xl text-text-muted mb-2"></i>
      <p class="text-text-muted">Silakan pilih divisi untuk mengelola karyawan.</p>
    </div>

    <Toast />
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
import debounce from 'lodash.debounce';

const toast = useToast();
const authStore = useAuthStore();

const divisions = ref([]);
const selectedDivisionId = ref(null);
const isFetchingDivisions = ref(false);

const employeeLists = ref([[], []]); // [employeesInDivision, otherEmployees]
const allOtherEmployees = ref([]); // Store all other employees for filtering
const isLoading = ref(false);
const searchTerm = ref('');

// Fetch all divisions
const fetchDivisions = async () => {
  isFetchingDivisions.value = true;
  try {
    const response = await axios.get('/api/admin/divisions');
    divisions.value = response.data.data.map(div => ({ id: div.ID, name: div.Name }));
    if (divisions.value.length > 0) {
      selectedDivisionId.value = divisions.value[0].id; // Select the first division by default
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat divisi.', life: 3000 });
  } finally {
    isFetchingDivisions.value = false;
  }
};

// Fetch employees based on selected division
const fetchEmployees = async () => {
  if (!authStore.companyId || !selectedDivisionId.value) {
    employeeLists.value = [[], []];
    return;
  }
  isLoading.value = true;
  try {
    // Fetch employees for the selected division (left list)
    const employeesInDivisionResponse = await axios.get(`/api/companies/${authStore.companyId}/employees`, {
      params: { division_id: selectedDivisionId.value }
    });
    const employeesInDivision = employeesInDivisionResponse.data.data.items.map(emp => ({
      id: emp.ID,
      name: emp.Name,
      email: emp.Email,
      position: emp.Position,
      division_id: emp.DivisionID
    }));

    // Fetch employees with no division (right list)
    const otherEmployeesResponse = await axios.get(`/api/companies/${authStore.companyId}/employees`, {
      params: { no_division: true }
    });
    const otherEmployees = otherEmployeesResponse.data.data.items.map(emp => ({
      id: emp.ID,
      name: emp.Name,
      email: emp.Email,
      position: emp.Position,
      division_id: emp.DivisionID
    }));

    // Filter out employees who are already in the selected division from the 'otherEmployees' list
    const filteredOtherEmployees = otherEmployees.filter(
      emp => !employeesInDivision.some(divEmp => divEmp.id === emp.id)
    );

    allOtherEmployees.value = filteredOtherEmployees; // Store for search filtering
    
    // Apply initial search filter if any
    employeeLists.value = [employeesInDivision, filterOtherEmployees(allOtherEmployees.value, searchTerm.value)];

  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengambil data karyawan.', life: 3000 });
    employeeLists.value = [[], []];
    allOtherEmployees.value = [];
  } finally {
    isLoading.value = false;
  }
};

// Handle employee movement between lists
const onPickListMove = async (event) => {
  // Update the local state immediately
  employeeLists.value = [event.source, event.target];

  const updates = [];
  if (event.direction === 'toTarget') { // Moved from Division to Other
    event.items.forEach(item => {
      updates.push({ employee_id: item.id, division_id: null }); // Assign to null (no division)
    });
  } else if (event.direction === 'toSource') { // Moved from Other to Division
    event.items.forEach(item => {
      updates.push({ employee_id: item.id, division_id: selectedDivisionId.value }); // Assign to selected division
    });
  }

  if (updates.length > 0) {
    try {
      for (const update of updates) {
        await axios.put(`/api/admin/employees/${update.employee_id}/division`, { division_id: update.division_id });
      }
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Divisi karyawan berhasil diperbarui.', life: 3000 });
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memperbarui divisi karyawan.', life: 3000 });
      // Revert UI changes if backend update fails
      fetchEmployees(); // Re-fetch to ensure data consistency
    }
  }
};

const getDivisionName = (divisionId) => {
  const division = divisions.value.find(d => d.id === divisionId);
  return division ? division.name : 'Tidak Diketahui';
};

const filterOtherEmployees = (employeesToFilter, term) => {
  if (!term) return employeesToFilter;
  const lowerCaseTerm = term.toLowerCase();
  return employeesToFilter.filter(
    emp => emp.name.toLowerCase().includes(lowerCaseTerm) ||
           emp.email.toLowerCase().includes(lowerCaseTerm) ||
           emp.position.toLowerCase().includes(lowerCaseTerm)
  );
};

const onSearch = debounce(() => {
  employeeLists.value[1] = filterOtherEmployees(allOtherEmployees.value, searchTerm.value);
}, 300);

// Watch for changes in selectedDivisionId to re-fetch employees
watch(selectedDivisionId, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    fetchEmployees();
  }
});

onMounted(() => {
  fetchDivisions();
});
</script>