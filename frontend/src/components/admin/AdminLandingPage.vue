<template>
  <div class="admin-login-page flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-gray-800 mb-6">Login Admin Perusahaan</h1>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label for="email" class="block text-gray-700 text-sm font-bold mb-2">Email:</label>
          <input
            type="email"
            id="email"
            v-model="email"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Masukkan email Anda"
            required
          />
        </div>

        <div class="mb-6">
          <label for="password" class="block text-gray-700 text-sm font-bold mb-2">Kata Sandi:</label>
          <input
            type="password"
            id="password"
            v-model="password"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Masukkan kata sandi Anda"
            required
          />
        </div>

        <div class="flex items-center justify-between mb-4">
          <button
            type="submit"
            class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
          >
            Login
          </button>
        </div>
      </form>

      <div class="text-center mt-4">
        <a
          href="#"
          @click.prevent="goToForgotPassword"
          class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
        >
          Lupa Kata Sandi?
        </a>
      </div>

      <div class="text-center mt-4">
        <p class="text-gray-600 text-sm">
          Belum punya akun?
          <a
            href="#"
            @click.prevent="goToPricingSection"
            class="font-bold text-blue-500 hover:text-blue-800"
          >
            Daftar Sekarang
          </a>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

export default {
  name: 'AdminLandingPage',
  setup() {
    const email = ref('');
    const password = ref('');
    const router = useRouter();

    const handleLogin = () => {
      // Implement your login logic here
      console.log('Login attempt with:', email.value, password.value);
      // Example: Call an API
      // axios.post('/api/login/admin-company', { email: email.value, password: password.value })
      //   .then(response => { /* handle success */ })
      //   .catch(error => { /* handle error */ });
    };

    const goToForgotPassword = () => {
      // Navigate to a dedicated forgot password page
      router.push({ name: 'ForgotPassword' }); // Assuming you'll create a route named 'ForgotPassword'
    };

    const goToPricingSection = () => {
      // Navigate to the LandingPage on the main domain
      // Assuming the main domain is accessible via window.location.origin or a configured base URL
      window.location.href = `${window.location.origin}/#pricing-section`;
    };

    return {
      email,
      password,
      handleLogin,
      goToForgotPassword,
      goToPricingSection,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>