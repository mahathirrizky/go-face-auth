<template>
  <header :class="['w-full', 'py-4', 'px-6', 'flex', 'justify-between', 'items-center', 'fixed', 'top-0', 'left-0', 'z-50', 'transition-all', 'duration-500' ,{ 'bg-bg-muted shadow-md': isScrolled, 'bg-transparent': !isScrolled }]">
    <div class="flex items-center">
      <!-- Logo -->
      <a href="/" class="text-4xl text-secondary font-bold">Hadir Bos</a>
    </div>

    <!-- Navigation Links (Desktop) -->
    <nav class="hidden md:flex space-x-8">
      <a @click="scrollToSection('features')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Fitur</a>
      <a @click="scrollToSection('testimonials')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Testimoni</a>
      <a @click="scrollToSection('pricing')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Harga</a>
      <a @click="scrollToSection('contact')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Kontak</a>
    </nav>

    <!-- Call to Action Button -->
    <div class="hidden md:block">
      <button class="btn btn-secondary">
        Daftar Sekarang
      </button>
    </div>

    <!-- Mobile Menu Button (Hamburger) -->
    <div class="md:hidden">
      <button :class="['focus:outline-none', 'text-text-base']">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
        </svg>
      </button>
    </div>
  </header>
</template>

<script>
export default {
  name: 'Header',
  data() {
    return {
      isScrolled: false,
    };
  },
  methods: {
    handleScroll() {
      this.isScrolled = window.scrollY > 0;
    },
    scrollToSection(id) {
      const element = document.getElementById(id);
      if (element) {
        element.scrollIntoView({
          behavior: 'smooth',
          block: 'start',
        });
      }
    },
  },
  mounted() {
    window.addEventListener('scroll', this.handleScroll);
    this.handleScroll();
  },
  beforeUnmount() {
    window.removeEventListener('scroll', this.handleScroll);
  },
}
</script>

<style scoped>
/* Tailwind handles styling */
</style>

