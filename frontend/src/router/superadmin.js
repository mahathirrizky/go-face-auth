import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const superadminRoutes = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/auth',
    name: 'SuperAdminAuth',
    component: () => import('../components/superadmin/SuperAdminAuth.vue'),
  },
  {
    path: '/dashboard',
    name: 'SuperAdminDashboard',
    component: () => import('../components/superadmin/SuperAdminDashboard.vue'),
    children: [
      {
        path: '', // Default child route
        name: 'SuperAdminDashboardOverview',
        component: () => import('../components/superadmin/SuperAdminDashboardOverview.vue'),
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/companies',
        name: 'SuperAdminCompanies',
        component: () => import('../components/superadmin/SuperAdminCompanies.vue'),
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/subscriptions',
        name: 'SuperAdminSubscriptions',
        component: () => import('../components/superadmin/SuperAdminSubscriptions.vue'),
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/revenue-chart',
        name: 'SuperAdminRevenueChart',
        component: () => import('../components/superadmin/SuperAdminRevenueChart.vue'),
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/subscription-packages',
        name: 'SuperAdminSubscriptionPackages',
        component: () => import('../components/superadmin/SuperAdminSubscriptionPackages.vue'),
        meta: { requiresAuth: true, role: 'superadmin' },
      },
    ],
    meta: { requiresAuth: true, role: 'superadmin' },
  },
  // Add more superadmin-specific routes here
  // Catch-all 404 route
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../components/main/NotFound.vue'),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes: superadminRoutes,
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  // Check for token and that the user object exists and has the correct role.
  const userIsSuperAdmin = authStore.token && authStore.user && authStore.user.role === 'super_admin';

  if (requiresAuth && !userIsSuperAdmin) {
    // If the route requires auth and the user is not a logged-in superadmin,
    // redirect to the login page.
    next({ name: 'SuperAdminAuth' });
  } else {
    // Otherwise, allow the navigation.
    next();
  }
});

export default router;
