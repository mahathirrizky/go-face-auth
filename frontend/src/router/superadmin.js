import { createRouter, createWebHistory } from 'vue-router';
import SuperAdminAuth from '../components/superadmin/SuperAdminAuth.vue';
import SuperAdminDashboard from '../components/superadmin/SuperAdminDashboard.vue';
import SuperAdminDashboardOverview from '../components/superadmin/SuperAdminDashboardOverview.vue';
import SuperAdminCompanies from '../components/superadmin/SuperAdminCompanies.vue';
import SuperAdminSubscriptions from '../components/superadmin/SuperAdminSubscriptions.vue';
import SuperAdminRevenueChart from '../components/superadmin/SuperAdminRevenueChart.vue';
import SuperAdminSubscriptionPackages from '../components/superadmin/SuperAdminSubscriptionPackages.vue';
import NotFound from '../components/main/NotFound.vue'; // Added for 404
import { useAuthStore } from '../stores/auth';

const superadminRoutes = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/auth',
    name: 'SuperAdminAuth',
    component: SuperAdminAuth,
  },
  {
    path: '/dashboard',
    name: 'SuperAdminDashboard',
    component: SuperAdminDashboard,
    children: [
      {
        path: '', // Default child route
        name: 'SuperAdminDashboardOverview',
        component: SuperAdminDashboardOverview,
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/companies',
        name: 'SuperAdminCompanies',
        component: SuperAdminCompanies,
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/subscriptions',
        name: 'SuperAdminSubscriptions',
        component: SuperAdminSubscriptions,
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/revenue-chart',
        name: 'SuperAdminRevenueChart',
        component: SuperAdminRevenueChart,
        meta: { requiresAuth: true, role: 'superadmin' },
      },
      {
        path: '/subscription-packages',
        name: 'SuperAdminSubscriptionPackages',
        component: SuperAdminSubscriptionPackages,
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
    component: NotFound,
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
