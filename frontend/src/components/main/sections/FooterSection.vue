<template>
  <footer class="bg-gray-800 text-white py-8">
    <div class="container mx-auto px-4 text-center">
      <p>&copy; {{ currentYear }} HadirBos.com. All rights reserved.</p>
      <div class="flex justify-center space-x-4 mt-4">
        <a href="#" class="text-gray-400 hover:text-white transition duration-300">Kebijakan Privasi</a>
        <a href="#" class="text-gray-400 hover:text-white transition duration-300">Syarat & Ketentuan</a>
      </div>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'FooterSection',
  data() {
    return {
      currentYear: new Date().getFullYear(),
    };
  },
};
</script>
