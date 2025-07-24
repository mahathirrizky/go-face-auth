<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6 text-text-base">Daftar Perusahaan Baru</h2>
      <BaseForm :resolver="resolver" :initialValues="initialValues" @submit="registerCompany" v-slot="{ $form }">
        <BaseInput
          id="companyName"
          name="company_name"
          label="Nama Perusahaan:"
          :invalid="$form.company_name?.invalid"
          :errors="$form.company_name?.errors"
          required
        />
        <BaseInput
          id="companyAddress"
          name="company_address"
          label="Alamat Perusahaan:"
          :invalid="$form.company_address?.invalid"
          :errors="$form.company_address?.errors"
        />
        <BaseInput
          id="adminEmail"
          name="admin_email"
          label="Email Admin:"
          type="email"
          :invalid="$form.admin_email?.invalid"
          :errors="$form.admin_email?.errors"
          required
        />
        <BaseInput
          id="adminPassword"
          name="admin_password"
          label="Password Admin:"
          type="password"
          :invalid="$form.admin_password?.invalid"
          :errors="$form.admin_password?.errors"
          :required="true"
          :toggleMask="true"
          :feedback="false"
          :fluid="true"
        />
        <BaseInput
          id="confirmAdminPassword"
          name="confirm_admin_password"
          label="Konfirmasi Password Admin:"
          type="password"
          :invalid="$form.confirm_admin_password?.invalid"
          :errors="$form.confirm_admin_password?.errors"
          :required="true"
          :toggleMask="true"
          :feedback="false"
          :fluid="true"
        />
        <BaseInput
          id="subscriptionPackage"
          name="subscription_package_id"
          v-model="selectedPackageName"
          label="Paket Langganan:"
          readonly
          class="cursor-not-allowed"
          :invalid="$form.subscription_package_id?.invalid"
          :errors="$form.subscription_package_id?.errors"
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

const props = defineProps(['packageId']);
const toast = useToast();
const route = useRoute();

const selectedPackageName = ref('');
const subscriptionPackages = ref([]);

const initialValues = ref({
  company_name: '',
  company_address: '',
  admin_email: '',
  admin_password: '',
  confirm_admin_password: '',
  subscription_package_id: null,
  billing_cycle: route.query.billingCycle || 'monthly',
});

const passwordSchema = z.string()
  .min(8, { message: 'Minimal 8 karakter.' })
  .refine((value) => /[a-z]/.test(value), {
    message: 'Minimal satu huruf kecil.'
  })
  .refine((value) => /[A-Z]/.test(value), {
    message: 'Minimal satu huruf besar.'
  })
  .refine((value) => /\d/.test(value), {
    message: 'Minimal satu angka.'
  });

const schema = z.object({
  company_name: z.string().min(1, { message: 'Nama perusahaan wajib diisi.' }),
  company_address: z.string().optional(),
  admin_email: z.string().email({ message: 'Email admin tidak valid.' }),
  admin_password: passwordSchema,
  confirm_admin_password: z.string(),
  subscription_package_id: z.number().nullable().refine(val => val !== null, { message: 'Paket langganan wajib dipilih.' }),
  billing_cycle: z.string().optional(),
}).refine((data) => data.admin_password === data.confirm_admin_password, {
  message: 'Konfirmasi password tidak cocok dengan password admin.',
  path: ['confirm_admin_password'],
});

const resolver = zodResolver(schema);

const fetchSubscriptionPackages = async () => {
  try {
    const response = await axios.get('/api/subscription-packages');
    subscriptionPackages.value = response.data.data;
  } catch (error) {
    console.error('Error fetching subscription packages:', error);
  }
};

watch(subscriptionPackages, (newPackages) => {
  if (newPackages.length > 0 && props.packageId) {
    const packageId = parseInt(props.packageId);
    const selectedPackage = newPackages.find(pkg => pkg.id === packageId);
    if (selectedPackage) {
      selectedPackageName.value = selectedPackage.package_name;
      initialValues.value.subscription_package_id = packageId;
    } else {
      selectedPackageName.value = 'Paket tidak ditemukan';
      initialValues.value.subscription_package_id = null;
    }
  }
}, { immediate: true });

onMounted(() => {
  fetchSubscriptionPackages();
});

const registerCompany = async (event) => {
  const { valid, data } = event;

  if (!valid) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Mohon lengkapi semua kolom yang wajib diisi.', life: 3000 });
    return;
  }

  console.log('Submitting registration data:', data);

  // Create a new object without confirm_admin_password
  const dataToSend = { ...data };
  delete dataToSend.confirm_admin_password;

  try {
    const response = await axios.post('/api/register-company', dataToSend);
    toast.add({ severity: 'success', summary: 'Success', detail: response.data.message, life: 3000 });
    setTimeout(() => {
      const baseDomain = getBaseDomain();
      const adminLoginURL = `${window.location.protocol}//admin.${baseDomain}${window.location.port ? ':' + window.location.port : ''}/`;
      window.location.href = adminLoginURL;
    }, 2000);
  } catch (error) {
    console.error('Registration failed:', error.response || error.message);
    const errorMessage = error.response?.data?.message || error.message;
    toast.add({ severity: 'error', summary: 'Error', detail: 'Registration failed: ' + errorMessage, life: 3000 });
  }
};
</script>