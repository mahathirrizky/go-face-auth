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

<script>
import Header from '../main/Header.vue';
import HeroSection from './sections/HeroSection.vue';
import FeaturesSection from './sections/Features.vue';
import TestimonialsSection from './sections/TestimonialsSection.vue';
import PricingSection from './sections/PricingSection.vue';
import ContactSection from './sections/ContactSection.vue';
import FooterSection from './sections/FooterSection.vue';
import ScrollToTopButton from '../main/ScrollToTopButton.vue';
import axios from 'axios';

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
  data() {
    return {
      subscriptionPackages: [],
    };
  },
  created() {
    this.fetchSubscriptionPackages();
  },
  methods: {
    async fetchSubscriptionPackages() {
      try {
        const response = await axios.get('/api/subscription-packages');
        this.subscriptionPackages = response.data.data;
        console.log('Fetched subscription packages:', this.subscriptionPackages);
      } catch (error) {
        console.error('Error fetching subscription packages:', error);
      }
    },
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