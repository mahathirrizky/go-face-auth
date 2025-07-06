import { createRouter, createWebHistory } from 'vue-router';
import AdminLandingPage from '../components/admin/AdminLandingPage.vue';
import AuthPage from '../components/admin/AuthPage.vue';
import { useAuthStore } from '../stores/auth';

import ForgotPassword from '../components/auth/ForgotPassword.vue';
import ResetPassword from '../components/auth/ResetPassword.vue';
import AdminDashboard from '../components/admin/AdminDashboard.vue';
import DashboardOverview from '../components/admin/DashboardOverview.vue'; // New import
import EmployeeManagement from '../components/admin/EmployeeManagement.vue';
import AttendanceManagement from '../components/admin/AttendanceManagement.vue';
import SettingsPage from '../components/admin/SettingsPage.vue';
import ShiftManagement from '../components/admin/ShiftManagement.vue'; // New import

const routes = [
  {
    path: '/',
    name: 'AdminLandingPage',
    component: AdminLandingPage,
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
  {
    path: '/dashboard',
    name: 'AdminDashboard',
    component: AdminDashboard,
    children: [
      {
        path: '', // Default child route for /dashboard
        name: 'DashboardOverview',
        component: DashboardOverview,
      },
      {
        path: 'employees',
        name: 'EmployeeManagement',
        component: EmployeeManagement,
      },
      {
        path: 'attendance',
        name: 'AttendanceManagement',
        component: AttendanceManagement,
      },
      {
        path: 'settings',
        name: 'SettingsPage',
        component: SettingsPage,
      },
      {
        path: 'shifts', // New route for Shift Management
        name: 'ShiftManagement',
        component: ShiftManagement,
      },
    ],
  },
];

const routeradmin = createRouter({
  history: createWebHistory(),
  routes,
});

routeradmin.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const publicPages = ['/login', '/', '/forgot-password', '/reset-password'];
  const authRequired = !publicPages.includes(to.path);

  // If user is logged in and tries to access the root path, redirect to dashboard
  if (authStore.token && to.path === '/') {
    return next('/dashboard');
  }

  if (authRequired && !authStore.token) {
    return next('/login');
  }

  next();
});

export default routeradmin;
