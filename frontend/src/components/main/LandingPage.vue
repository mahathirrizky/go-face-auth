<template>
  <div class="relative min-h-screen bg-bg-base flex flex-col">
    <Header />
    <HeroSection />
    <FeaturesSection />
    <TestimonialsSection />
    <PricingSection :packages="subscriptionPackages" />
    <ContactSection />
    <FooterSection />
    <ScrollToTopButton />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import Header from '../main/Header.vue';
import HeroSection from './sections/HeroSection.vue';
import FeaturesSection from './sections/Features.vue';
import TestimonialsSection from './sections/TestimonialsSection.vue';
import PricingSection from './sections/PricingSection.vue';
import ContactSection from './sections/ContactSection.vue';
import FooterSection from './sections/FooterSection.vue';
import ScrollToTopButton from '../main/ScrollToTopButton.vue';
import axios from 'axios';

const subscriptionPackages = ref([]);

const fetchSubscriptionPackages = async () => {
  try {
    const response = await axios.get('/api/subscription-packages');
    subscriptionPackages.value = response.data.data;
    console.log('Fetched subscription packages:', subscriptionPackages.value);
  } catch (error) {
    console.error('Error fetching subscription packages:', error);
  }
};

onMounted(() => {
  fetchSubscriptionPackages();
});
</script>

<style scoped>
/* No custom scoped styles needed as Tailwind handles styling */
</style>