import { createRouter, createWebHistory } from 'vue-router';
import AdminLandingPage from '../components/admin/AdminLandingPage.vue';
import { useAuthStore } from '../stores/auth';

import ForgotPassword from '../components/auth/ForgotPassword.vue';
import ResetPassword from '../components/auth/ResetPassword.vue';
import ConfirmEmail from '../components/auth/ConfirmEmail.vue';
import AdminDashboard from '../components/admin/AdminDashboard.vue';
import DashboardOverview from '../components/admin/DashboardOverview.vue'; // New import
import EmployeeManagement from '../components/admin/EmployeeManagement.vue';
import EmployeeAttendanceHistory from '../components/admin/EmployeeAttendanceHistory.vue';
import AttendanceManagement from '../components/admin/AttendanceManagement.vue';
import SettingsPage from '../components/admin/SettingsPage.vue';
import ShiftManagement from '../components/admin/ShiftManagement.vue';
import GeneralSettings from '../components/admin/GeneralSettings.vue'; // New import
import AdminAccountSettings from '../components/admin/AdminAccountSettings.vue'; // New import
import SubscriptionPage from '../components/admin/SubscriptionPage.vue'; // New import
import PaymentPage from '../components/admin/PaymentPage.vue';
import NotFound from '../components/main/NotFound.vue'; // Added for 404

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
    path: '/confirm-email',
    name: 'ConfirmEmail',
    component: ConfirmEmail,
  },
  {
    path: '/dashboard',
    name: 'AdminDashboard',
    component: AdminDashboard,
    children: [
      {
        path: '',
        name: 'DashboardOverview',
        component: DashboardOverview,
      },
      {
        path: 'employees',
        name: 'EmployeeManagement',
        component: EmployeeManagement,
      },
      {
        path: 'employees/:employeeId/attendance-history',
        name: 'EmployeeAttendanceHistory',
        component: EmployeeAttendanceHistory,
      },
      {
        path: 'attendance',
        name: 'AttendanceManagement',
        component: AttendanceManagement,
      },
      {
        path: 'subscribe',
        name: 'SubscriptionPage',
        component: SubscriptionPage,
      },
      {
        path: 'payment/:companyId',
        name: 'PaymentPage',
        component: PaymentPage,
        props: true,
      },
      {
        path: 'settings',
        name: 'SettingsPage',
        component: SettingsPage,
        children: [
          {
            path: '',
            name: 'SettingsRedirect', // Add a name to the redirect route
            redirect: { name: 'GeneralSettings' }, // Redirect to General Settings by default
          },
          {
            path: 'general',
            name: 'GeneralSettings',
            component: GeneralSettings,
          },
          {
            path: 'admin-account',
            name: 'AdminAccountSettings',
            component: AdminAccountSettings,
          },
          {
            path: 'shifts',
            name: 'ShiftManagement',
            component: ShiftManagement,
          },
        ],
      },
    ],
  },
  // Catch-all 404 route
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
  },
];

const routeradmin = createRouter({
  history: createWebHistory(),
  routes,
});

routeradmin.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const publicPages = ['/', '/forgot-password', '/reset-password', '/confirm-email'];
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
