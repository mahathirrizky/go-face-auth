<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Divisi</h2>
    <Tabs value="0">
      <TabList>
        <Tab value="0">Daftar Divisi</Tab>
        <Tab value="1">Tetapkan Karyawan</Tab>
      </TabList>
      <TabPanels>
        <TabPanel value="0">
          <DivisionList />
        </TabPanel>
        <TabPanel value="1">
          <AssignEmployeeToDivision />
        </TabPanel>
      </TabPanels>
    </Tabs>
  </div>
</template>

<script setup>
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import DivisionList from './division/DivisionList.vue';
import AssignEmployeeToDivision from './division/AssignEmployeeToDivision.vue';
</script>