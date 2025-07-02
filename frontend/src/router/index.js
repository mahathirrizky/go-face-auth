import { createRouter, createWebHistory } from 'vue-router';
import LandingPage from '../components/main/LandingPage.vue';

const routes = [
  {
    path: '/',
    name: 'LandingPage',
    component: LandingPage,
  },
  // Add other routes here as needed
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
