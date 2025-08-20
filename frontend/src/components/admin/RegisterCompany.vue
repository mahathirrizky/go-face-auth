<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6 text-text-base">Daftar Perusahaan Baru</h2>
      <Form :resolver="resolver" @submit="registerCompany" v-slot="{ errors, handleSubmit }">
        <div class="mb-4">
           <FloatLabel variant="on">

             <label for="companyName" class="block mb-1 text-text-base">Nama Perusahaan:</label>
             <InputText
             id="companyName"
             name="company_name"
             v-model="formData.company_name"
             class="w-full border rounded-md p-2"
             :class="{ 'border-red-500': errors?.company_name || serverErrors.company_name }"
             required
             fluid
             />
            </FloatLabel>
          <small v-if="errors?.company_name || serverErrors.company_name" class="text-red-500">{{ errors?.company_name || serverErrors.company_name }}</small>
        </div>
        <div class="mb-4">
           <FloatLabel variant="on">

             <label for="companyAddress" class="block mb-1 text-text-base">Alamat Perusahaan:</label>
             <InputText
          id="companyAddress"
          name="company_address"
          v-model="formData.company_address"
          class="w-full border rounded-md p-2"
          :class="{ 'border-red-500': errors?.company_address || serverErrors.company_address }"
          fluid
          />
        </FloatLabel>
          <small v-if="errors?.company_address || serverErrors.company_address" class="text-red-500">{{ errors?.company_address || serverErrors.company_address }}</small>
        </div>
        <div class="mb-4">
          <FloatLabel variant="on">

            <label for="adminEmail" class="block mb-1 text-text-base">Email Admin:</label>
            <InputText
            id="adminEmail"
            name="admin_email"
            v-model="formData.admin_email"
            type="email"
            class="w-full border rounded-md p-2"
            :class="{ 'border-red-500': errors?.admin_email || serverErrors.admin_email }"
            required
            fluid
            />
          </FloatLabel>
          <small v-if="errors?.admin_email || serverErrors.admin_email" class="text-red-500">{{ errors?.admin_email || serverErrors.admin_email }}</small>
        </div>
        <div class="mb-4">
           <FloatLabel variant="on">

             <Password
             id="adminPassword"
             name="admin_password"
             v-model="formData.admin_password"
             class="w-full"
            :class="{ 'p-invalid': errors?.admin_password || serverErrors.admin_password }"
            :toggleMask="true"
            :feedback="true"
            required
            inputClass="w-full border rounded-md p-2"
            fluid
            >
              <template #footer>
                <p class="mt-2">Saran Kata Sandi:</p>
                <ul class="pl-4 ml-2 mt-0" style="line-height: 1.5">
                  <li :class="{ 'text-green-500': isLengthValid }">Minimal 8 karakter <i v-if="isLengthValid" class="pi pi-check"></i></li>
                  <li :class="{ 'text-green-500': hasLowercase }">Minimal satu huruf kecil (a-z) <i v-if="hasLowercase" class="pi pi-check"></i></li>
                  <li :class="{ 'text-green-500': hasUppercase }">Minimal satu huruf besar (A-Z) <i v-if="hasUppercase" class="pi pi-check"></i></li>
                  <li :class="{ 'text-green-500': hasNumber }">Minimal satu angka (0-9) <i v-if="hasNumber" class="pi pi-check"></i></li>
                </ul>
              </template>
            </Password>
            <label for="adminPassword" class="block mb-1 text-text-base">Password Admin:</label>
          </FloatLabel>
          <small v-if="errors?.admin_password || serverErrors.admin_password" class="text-red-500">{{ errors?.admin_password || serverErrors.admin_password }}</small>
        </div>
        <div class="mb-4">
           <FloatLabel variant="on">

             <Password
             id="confirmAdminPassword"
             name="confirm_admin_password"
             v-model="formData.confirm_admin_password"
             class="w-full"
             :class="{ 'p-invalid': errors?.confirm_admin_password || serverErrors.confirm_admin_password }"
             :toggleMask="true"
             :feedback="false"
             required
             inputClass="w-full border rounded-md p-2"
             fluid
             />
             <label for="confirmAdminPassword" class="block mb-1 text-text-base">Konfirmasi Password Admin:</label>
          </FloatLabel>
          <small v-if="errors?.confirm_admin_password || serverErrors.confirm_admin_password" class="text-red-500">
            {{ errors?.confirm_admin_password || serverErrors.confirm_admin_password }}
          </small>
        </div>
        <div v-if="isLoadingPackages" class="flex items-center justify-center py-2">
          <i class="pi pi-spin pi-spinner mr-2 text-text-muted"></i>
          <span class="text-text-muted">Memuat paket langganan...</span>
        </div>
        <div v-else class="mb-4">
          
          <label for="subscriptionPackage" class="block mb-1 text-text-base">Paket Langganan:</label>
          <InputText
            id="subscriptionPackage"
            name="subscription_package_id"
            :value="selectedPackageName"
            class="w-full border rounded-md p-2 cursor-not-allowed bg-gray-100"
            :class="{ 'border-red-500': errors?.subscription_package_id || serverErrors.subscription_package_id }"
            disabled
          />
          <small v-if="errors?.subscription_package_id || serverErrors.subscription_package_id" class="text-red-500">
            {{ errors?.subscription_package_id || serverErrors.subscription_package_id }}
          </small>
        </div>
        <Button
          type="submit"
          :label="isRegistering ? 'Mendaftar...' : 'Daftar & Mulai Coba Gratis'"
          class="w-full mt-6 btn-primary"
          :disabled="isRegistering"
          :icon="isRegistering ? 'pi pi-spin pi-spinner' : 'pi pi-check'"
        />
      </Form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useToast } from 'primevue/usetoast';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { getBaseDomain } from '../../utils/subdomain';
import { Form } from '@primevue/forms';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Button from 'primevue/button';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import FloatLabel from 'primevue/floatlabel';

const props = defineProps(['packageId']);
const toast = useToast();
const route = useRoute();

const selectedPackageName = ref('');
const subscriptionPackages = ref([]);
const isRegistering = ref(false);
const isLoadingPackages = ref(false);
const formData = ref({
  company_name: '',
  company_address: '',
  admin_email: '',
  admin_password: '',
  confirm_admin_password: '',
  subscription_package_id: props.packageId ? parseInt(props.packageId, 10) : null,
  billing_cycle: route.query.billingCycle || 'monthly',
});
const serverErrors = ref({});

const isLengthValid = computed(() => formData.value.admin_password.length >= 8);
const hasLowercase = computed(() => /[a-z]/.test(formData.value.admin_password));
const hasUppercase = computed(() => /[A-Z]/.test(formData.value.admin_password));
const hasNumber = computed(() => /\d/.test(formData.value.admin_password));

const passwordSchema = z
  .string()
  .min(8, { message: 'Minimal 8 karakter.' })
  .refine((value) => /[a-z]/.test(value), { message: 'Minimal satu huruf kecil.' })
  .refine((value) => /[A-Z]/.test(value), { message: 'Minimal satu huruf besar.' })
  .refine((value) => /\d/.test(value), { message: 'Minimal satu angka.' });

const schema = z
  .object({
    company_name: z.string().min(1, { message: 'Nama perusahaan wajib diisi.' }),
    company_address: z.string().optional(),
    admin_email: z.string().email({ message: 'Email admin tidak valid.' }),
    admin_password: passwordSchema,
    confirm_admin_password: z.string(),
    subscription_package_id: z
      .number()
      .nullable()
      .refine((val) => val !== null, { message: 'Paket langganan wajib dipilih.' }),
    billing_cycle: z.string().optional(),
  })
  .refine((data) => data.admin_password === data.confirm_admin_password, {
    message: 'Konfirmasi password tidak cocok dengan password admin.',
    path: ['confirm_admin_password'],
  });

const resolver = zodResolver(schema);

const fetchSubscriptionPackages = async () => {
  isLoadingPackages.value = true;
  try {
    const response = await axios.get('/api/subscription-packages');
    subscriptionPackages.value = response.data.data;
    console.log('Fetched subscription packages:', subscriptionPackages.value);
  } catch (error) {
    console.error('Error fetching subscription packages:', error);
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Gagal memuat paket langganan.',
      life: 3000,
    });
  } finally {
    isLoadingPackages.value = false;
  }
};

onMounted(async () => {
  console.log('props.packageId:', props.packageId);
  await fetchSubscriptionPackages();
  if (props.packageId && subscriptionPackages.value.length > 0) {
    const packageId = parseInt(props.packageId, 10);
    if (isNaN(packageId)) {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'ID paket tidak valid dari URL.',
        life: 3000,
      });
      return;
    }
    const selectedPackage = subscriptionPackages.value.find((pkg) => pkg.id === packageId);
    if (selectedPackage) {
      selectedPackageName.value = selectedPackage.package_name;
      formData.value.subscription_package_id = packageId;
      console.log('Selected package:', selectedPackage);
    } else {
      selectedPackageName.value = 'Paket tidak ditemukan';
      formData.value.subscription_package_id = null;
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Paket langganan tidak ditemukan.',
        life: 3000,
      });
    }
  } else {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'ID paket tidak disediakan. Silakan pilih paket.',
      life: 3000,
    });
  }
  console.log('Form data after onMounted:', formData.value);
});

const registerCompany = async ({ values: data }) => {
  serverErrors.value = {}; // Clear previous server errors

  const dataToSend = { ...data };
  delete dataToSend.confirm_admin_password;
  console.log('API payload:', dataToSend);

  isRegistering.value = true;
  try {
    const response = await axios.post('/api/register-company', dataToSend);
    toast.add({
      severity: 'success',
      summary: 'Sukses',
      detail: response.data.message,
      life: 3000,
    });
    setTimeout(() => {
      const baseDomain = getBaseDomain();
      const adminLoginURL = `${window.location.protocol}//admin.${baseDomain}${
        window.location.port ? ':' + window.location.port : ''
      }/`;
      window.location.href = adminLoginURL;
    }, 2000);
  } catch (error) {
    console.error('Registration failed:', error.response || error.message);

    if (error.response && error.response.data && error.response.data.errors) {
      serverErrors.value = error.response.data.errors;
    }

    const errorMessage = error.response?.data?.message || error.message;
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Pendaftaran gagal: ' + errorMessage,
      life: 3000,
    });
  } finally {
    isRegistering.value = false;
  }
};
</script>

<style scoped>

</style>