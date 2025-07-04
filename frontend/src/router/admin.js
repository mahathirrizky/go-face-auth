import { createRouter, createWebHistory } from 'vue-router';
import AdminLandingPage from '../components/admin/AdminLandingPage.vue';
import AuthPage from '../components/admin/AuthPage.vue';
import ForgotPassword from '../components/auth/ForgotPassword.vue'; // New component
import ResetPassword from '../components/auth/ResetPassword.vue';   // New component

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
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: ForgotPassword,
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: ResetPassword,
  },
  // Add other admin-specific routes here
];

const routeradmin = createRouter({
  history: createWebHistory(),
  routes,
});

export default routeradmin;
