import { createRouter, createWebHistory } from 'vue-router';
import SuperAdminAuth from '../components/superadmin/SuperAdminAuth.vue';
import SuperAdminDashboard from '../components/superadmin/SuperAdminDashboard.vue';
import SuperAdminDashboardOverview from '../components/superadmin/SuperAdminDashboardOverview.vue';
import SuperAdminCompanies from '../components/superadmin/SuperAdminCompanies.vue';
import SuperAdminSubscriptions from '../components/superadmin/SuperAdminSubscriptions.vue';
import SuperAdminRevenueChart from '../components/superadmin/SuperAdminRevenueChart.vue';
import SuperAdminSubscriptionPackages from '../components/superadmin/SuperAdminSubscriptionPackages.vue';

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
];

const router = createRouter({
  history: createWebHistory(),
  routes: superadminRoutes,
});

export default router;
