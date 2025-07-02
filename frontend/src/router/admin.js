import { createRouter, createWebHistory } from 'vue-router';
import AdminLandingPage from '../components/admin/AdminLandingPage.vue';
import AuthPage from '../components/admin/AuthPage.vue';

const routes = [
  {
    path: '/',
    name: 'AdminLandingPage',
    component: AdminLandingPage,
  },
  {
    path: '/login',
    name: 'AuthPage',
    component: AuthPage,
  },
  // Add other admin-specific routes here
];

const routeradmin = createRouter({
  history: createWebHistory(),
  routes,
});

export default routeradmin;
