<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6 text-text-base">Daftar Perusahaan Baru</h2>
      <BaseForm :resolver="resolver" :initialValues="initialValues" @submit="registerCompany" v-slot="{ $form }">
        <BaseInput
          id="companyName"
          name="company_name"
          label="Nama Perusahaan:"
          required
          :invalid="$form.company_name?.invalid"
          :errorMessage="$form.company_name?.error?.message"
        />
        <BaseInput
          id="companyAddress"
          name="company_address"
          label="Alamat Perusahaan:"
          :invalid="$form.company_address?.invalid"
          :errorMessage="$form.company_address?.error?.message"
        />
        <BaseInput
          id="adminEmail"
          name="admin_email"
          label="Email Admin:"
          type="email"
          required
          :invalid="$form.admin_email?.invalid"
          :errorMessage="$form.admin_email?.error?.message"
        />
        <div class="flex flex-col gap-1">
          <BaseInput
            id="adminPassword"
            name="admin_password"
            label="Password Admin:"
            type="password"
            placeholder="Minimal 8 karakter"
            :required="true"
            :toggleMask="true"
            :feedback="false"
            :invalid="$form.admin_password?.invalid"
          />
          <template v-if="$form.admin_password?.invalid">
              <Message v-for="(error, index) of $form.admin_password.errors" :key="index" severity="error" size="small" variant="simple">{{ error.message }}</Message>
          </template>
        </div>
        <BaseInput
          id="subscriptionPackage"
          name="subscription_package_id"
          label="Paket Langganan:"
          v-model="selectedPackageName"
          readonly
          class="cursor-not-allowed"
          :invalid="$form.subscription_package_id?.invalid"
          :errorMessage="$form.subscription_package_id?.error?.message"
        />
        <BaseButton :fullWidth="true" class="mt-6 btn-primary" type="submit">
          <i class="pi pi-check"></i> Daftar & Mulai Coba Gratis
        </BaseButton>
      </BaseForm>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useToast } from 'primevue/usetoast';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { getBaseDomain } from '../../utils/subdomain';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseForm from '../ui/BaseForm.vue'; // Import BaseForm
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import Message from 'primevue/message';

const props = defineProps(['packageId']);
const toast = useToast();
const route = useRoute();

const form = ref({
  company_name: '',
  company_address: '',
  admin_email: '',
  admin_password: '',
  subscription_package_id: parseInt(props.packageId),
  billing_cycle: route.query.billingCycle || 'monthly',
});
const selectedPackageName = ref('');
const subscriptionPackages = ref([]);

const fetchSubscriptionPackages = async () => {
  try {
    const response = await axios.get('/api/subscription-packages');
    subscriptionPackages.value = response.data.data;
  } catch (error) {
    console.error('Error fetching subscription packages:', error);
  }
};

watch([subscriptionPackages, () => form.value.subscription_package_id], ([newPackages, newPackageId]) => {
  if (newPackages.length > 0 && newPackageId) {
    const selectedPackage = newPackages.find(pkg => pkg.id === newPackageId);
    if (selectedPackage) {
      selectedPackageName.value = selectedPackage.package_name;
    } else {
      selectedPackageName.value = 'Paket tidak ditemukan';
    }
  }
}, { immediate: true });

onMounted(() => {
  fetchSubscriptionPackages();
});

const registerCompanySchema = z.object({
  company_name: z.string().min(1, { message: 'Nama perusahaan wajib diisi.' }),
  company_address: z.string().optional(),
  admin_email: z.string().email({ message: 'Email admin tidak valid.' }),
  admin_password: z.string()
    .min(8, { message: 'Kata sandi minimal 8 karakter.' })
    .refine((value) => /[a-z]/.test(value), {
      message: 'Kata sandi harus memiliki setidaknya satu huruf kecil.'
    })
    .refine((value) => /[A-Z]/.test(value), {
      message: 'Kata sandi harus memiliki setidaknya satu huruf besar.'
    })
    .refine((value) => /\d/.test(value), {
      message: 'Kata sandi harus memiliki setidaknya satu angka.'
    }),
  subscription_package_id: z.number().int().positive({ message: 'Paket langganan wajib dipilih.' }),
});

const resolver = zodResolver(registerCompanySchema);

const initialValues = ref({
  company_name: '',
  company_address: '',
  admin_email: '',
  admin_password: '',
  subscription_package_id: parseInt(props.packageId),
});

const registerCompany = async (event) => {
  const { valid, data } = event;

  if (!valid) {
    toast.add({ severity: 'error', summary: 'Validasi Gagal', detail: 'Silakan periksa kembali input Anda.', life: 3000 });
    return;
  }

  try {
    const response = await axios.post('/api/register-company', data);
    toast.add({ severity: 'success', summary: 'Success', detail: response.data.message, life: 3000 });
    setTimeout(() => {
      const baseDomain = getBaseDomain();
      const adminLoginURL = `${window.location.protocol}//admin.${baseDomain}${window.location.port ? ':' + window.location.port : ''}/`;
      window.location.href = adminLoginURL;
    }, 2000);
  } catch (error) {
    const errorMessage = error.response?.data?.message || error.message;
    toast.add({ severity: 'error', summary: 'Error', detail: 'Registration failed: ' + errorMessage, life: 3000 });
  }
};

</script>