<template>
  <div class="relative min-h-screen bg-blue-700 flex flex-col">
    <Header />
    <HeroSection />
    <FeaturesSection />
    <TestimonialsSection />
    <PricingSection />
    <ContactSection />
    <FooterSection />
    <ScrollToTopButton />
  </div>
</template>

<script>
import Header from '../main/Header.vue';
import HeroSection from './sections/HeroSection.vue';
import FeaturesSection from './sections/Features.vue';
import TestimonialsSection from './sections/TestimonialsSection.vue';
import PricingSection from './sections/PricingSection.vue';
import ContactSection from './sections/ContactSection.vue';
import FooterSection from './sections/FooterSection.vue';
import ScrollToTopButton from '../main/ScrollToTopButton.vue';

export default {
  name: 'LandingPage',
  components: {
    Header,
    HeroSection,
    FeaturesSection,
    TestimonialsSection,
    PricingSection,
    ContactSection,
    FooterSection,
    ScrollToTopButton,
  },
  methods: {
    goToDashboard() {
      console.log('Navigating to dashboard...');
      // this.$router.push('/dashboard');
    }
  }
}
</script>

<style scoped>
/* No custom scoped styles needed as Tailwind handles styling */
</style>
