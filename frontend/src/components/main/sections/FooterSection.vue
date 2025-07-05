<template>
  <footer class="bg-bg-muted text-text-muted py-8">
    <div class="container mx-auto px-4 text-center">
      <p>&copy; {{ currentYear }} HadirBos.com. All rights reserved.</p>
      <div class="flex justify-center space-x-4 mt-4">
        <a href="#" class="text-text-muted hover:text-text-base transition duration-300">Kebijakan Privasi</a>
        <a href="#" class="text-text-muted hover:text-text-base transition duration-300">Syarat & Ketentuan</a>
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
