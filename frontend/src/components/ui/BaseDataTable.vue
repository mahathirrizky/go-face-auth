<template>
  <DataTable
    :value="data"
    :loading="loading"
    :paginator="true"
    :rows="10"
    :totalRecords="totalRecords"
    :lazy="lazy"
    @page="$emit('page', $event)"
    @filter="$emit('filter', $event)"
    v-model:filters="localFilters"
    class="p-datatable-customers"
    paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
    :rowsPerPageOptions="[10, 25, 50]"
    currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"

    :editMode="editModeType"
    :dataKey="editKey"
    @cell-edit-complete="onCellEditComplete"
  >
    <template #header>
      <div class="flex flex-wrap justify-between items-center gap-4">
        <IconField iconPosition="left">
          <InputIcon class="pi pi-search"></InputIcon>
          <InputText v-model="localFilters.global.value" :placeholder="searchPlaceholder" />
        </IconField>
        <div class="flex flex-wrap gap-2">
          <slot name="header-actions"></slot>
        </div>
      </div>
    </template>
    <template #empty>
      Tidak ada data ditemukan.
    </template>
    <template #loading>
      Memuat data...
    </template>

    <Column
      v-for="col in columns"
      :key="col.field"
      :field="col.field"
      :header="col.header"
      :sortable="col.sortable !== false"
      :showFilterMenu="col.showFilterMenu !== false"
      :editor="col.editable !== false ? (slotProps) => {
        // Check if a custom editor slot is provided for this column
        if ($slots[`editor-${col.field}`]) {
          return h('div', {}, $slots[`editor-${col.field}`](slotProps));
        }
        // Otherwise, use the default InputText editor
        return h(InputText, {
          modelValue: slotProps.data[slotProps.field],
          'onUpdate:modelValue': (value) => (slotProps.data[slotProps.field] = value),
          class: 'w-full'
        });
      } : null"
    >
      <template #body="slotProps">
        <slot :name="`column-${col.field}`" :item="slotProps.data">
          {{ slotProps.data[col.field] }}
        </slot>
      </template>
    </Column>

    <!-- Dedicated Custom Actions Column (for delete, view, etc.) -->
    <Column v-if="$slots.actions" header="Aksi" style="width: 10%; min-width: 8rem" bodyStyle="text-align:center">
        <template #body="slotProps">
            <slot name="actions" :item="slotProps.data"></slot>
        </template>
    </Column>

  </DataTable>
</template>

<script setup>
import { ref, computed, h } from 'vue';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';


import { FilterMatchMode } from '@primevue/core/api';

const props = defineProps({
  data: { type: Array, required: true },
  columns: { type: Array, required: true },
  loading: { type: Boolean, default: false },
  searchPlaceholder: { type: String, default: 'Cari...' },
  lazy: { type: Boolean, default: false },
  totalRecords: { type: Number, default: 0 },
  filters: {
    type: Object,
    default: () => ({ global: { value: null, matchMode: FilterMatchMode.CONTAINS } })
  },
  // Props for Editing
  editModeType: { type: String, default: null }, // Changed from editable to editModeType
  editKey: { type: String, default: 'id' }
});

const emit = defineEmits(['page', 'filter', 'update:filters', 'cell-edit-complete']); // Updated emit

const localFilters = computed({
  get: () => props.filters,
  set: (value) => emit('update:filters', value)
});

const onCellEditComplete = (event) => {
  emit('cell-edit-complete', event); // Emit the cell-edit-complete event
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>