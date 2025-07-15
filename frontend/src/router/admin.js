import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const routes = [
  {
    path: '/',
    name: 'AdminLandingPage',
    component: () => import('../components/admin/AdminLandingPage.vue'),
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('../components/auth/ForgotPassword.vue'),
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('../components/auth/ResetPassword.vue'),
  },
  {
    path: '/confirm-email',
    name: 'ConfirmEmail',
    component: () => import('../components/auth/ConfirmEmail.vue'),
  },
  {
    path: '/dashboard',
    name: 'AdminDashboard',
    component: () => import('../components/admin/AdminDashboard.vue'),
    children: [
      {
        path: '',
        name: 'DashboardOverview',
        component: () => import('../components/admin/DashboardOverview.vue'),
      },
      {
        path: 'employees',
        name: 'EmployeeManagement',
        component: () => import('../components/admin/EmployeeManagement.vue'),
      },
      {
        path: 'employees/:employeeId/attendance-history',
        name: 'EmployeeAttendanceHistory',
        component: () => import('../components/admin/EmployeeAttendanceHistory.vue'),
      },
      {
        path: 'attendance',
        name: 'AttendanceManagement',
        component: () => import('../components/admin/AttendanceManagement.vue'),
      },
      {
        path: 'leave-requests',
        name: 'LeaveRequestManagement',
        component: () => import('../components/admin/LeaveRequestManagement.vue'),
      },
      {
        path: 'broadcast',
        name: 'BroadcastMessage',
        component: () => import('../components/admin/BroadcastMessagePage.vue'),
      },
      {
        path: 'subscribe',
        name: 'SubscriptionPage',
        component: () => import('../components/admin/SubscriptionPage.vue'),
      },
      {
        path: 'payment/:companyId',
        name: 'PaymentPage',
        component: () => import('../components/admin/PaymentPage.vue'),
        props: true,
      },
      {
        path: 'settings',
        name: 'SettingsPage',
        component: () => import('../components/admin/SettingsPage.vue'),
        children: [
          {
            path: '',
            name: 'SettingsRedirect',
            redirect: { name: 'GeneralSettings' },
          },
          {
            path: 'general',
            name: 'GeneralSettings',
            component: () => import('../components/admin/GeneralSettings.vue'),
          },
          {
            path: 'admin-account',
            name: 'AdminAccountSettings',
            component: () => import('../components/admin/AdminAccountSettings.vue'),
          },
          {
            path: 'shifts',
            name: 'ShiftManagement',
            component: () => import('../components/admin/ShiftManagement.vue'),
          },
        ],
      },
      {
        path: 'locations',
        name: 'LocationManagement',
        component: () => import('../components/admin/locations/LocationManagement.vue'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../components/main/NotFound.vue'),
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
