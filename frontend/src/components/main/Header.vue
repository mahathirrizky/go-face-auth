<template>
  <header :class="[
    'w-full', 'py-4', 'px-6', 'flex', 'justify-between', 'items-center', 'fixed', 'top-0', 'left-0', 'z-50', 'transition-all', 'duration-500',
    { 'bg-bg-muted shadow-md': isScrolled, 'bg-transparent': !isScrolled }
  ]">
    <div class="flex items-center">
      <a href="/" class="text-4xl text-secondary font-bold">Hadir Bos</a>
    </div>

    <nav class="hidden md:flex space-x-8">
      <a @click="scrollToSection('features')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Fitur</a>
      <a @click="scrollToSection('testimonials')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Testimoni</a>
      <a @click="scrollToSection('pricing')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Harga</a>
      <a @click="scrollToSection('contact')" class="cursor-pointer font-medium transition-colors duration-300 text-text-muted hover:text-text-base">Kontak</a>
    </nav>

    <div class="hidden md:block">
      <Button @click="scrollToSection('pricing')" icon="pi pi-play" label="Mulai Coba Gratis" />
    </div>

    <div class="md:hidden">
      <Button icon="pi pi-bars" @click="toggleMobileMenu" class="p-button-text text-text-base" />
    </div>

    <Sidebar v-model:visible="mobileMenuOpen" position="right">
      <template #header>
        <h3 class="text-2xl font-bold">Menu</h3>
      </template>
      <ul class="flex flex-col gap-4">
        <li><a @click="scrollToSectionAndCloseMenu('features')" class="cursor-pointer text-xl font-medium text-text-base hover:text-secondary">Fitur</a></li>
        <li><a @click="scrollToSectionAndCloseMenu('testimonials')" class="cursor-pointer text-xl font-medium text-text-base hover:text-secondary">Testimoni</a></li>
        <li><a @click="scrollToSectionAndCloseMenu('pricing')" class="cursor-pointer text-xl font-medium text-text-base hover:text-secondary">Harga</a></li>
        <li><a @click="scrollToSectionAndCloseMenu('contact')" class="cursor-pointer text-xl font-medium text-text-base hover:text-secondary">Kontak</a></li>
        <li><Button @click="scrollToSectionAndCloseMenu('pricing')" icon="pi pi-play" label="Mulai Coba Gratis" class="w-full mt-4" /></li>
      </ul>
    </Sidebar>
  </header>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import Button from 'primevue/button';
import Sidebar from 'primevue/sidebar';

const isScrolled = ref(false);
const mobileMenuOpen = ref(false);

const handleScroll = () => {
  isScrolled.value = window.scrollY > 0;
};

const scrollToSection = (id) => {
  const element = document.getElementById(id);
  if (element) {
    element.scrollIntoView({
      behavior: 'smooth',
      block: 'start',
    });
  }
};

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value;
};

const scrollToSectionAndCloseMenu = (id) => {
  scrollToSection(id);
  mobileMenuOpen.value = false;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
  handleScroll();
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
/* Tailwind handles styling */
</style>