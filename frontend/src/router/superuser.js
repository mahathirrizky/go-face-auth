import { createRouter, createWebHistory } from 'vue-router';
import SuperUserAuth from '../components/superuser/SuperUserAuth.vue';
import SuperUserDashboard from '../components/superuser/SuperUserDashboard.vue';
import SuperUserDashboardOverview from '../components/superuser/SuperUserDashboardOverview.vue';
import SuperUserCompanies from '../components/superuser/SuperUserCompanies.vue';
import SuperUserSubscriptions from '../components/superuser/SuperUserSubscriptions.vue';
import SuperUserRevenueChart from '../components/superuser/SuperUserRevenueChart.vue';
import SuperUserSubscriptionPackages from '../components/superuser/SuperUserSubscriptionPackages.vue';

const superuserRoutes = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/auth',
    name: 'SuperUserAuth',
    component: SuperUserAuth,
  },
  {
    path: '/dashboard',
    name: 'SuperUserDashboard',
    component: SuperUserDashboard,
    children: [
      {
        path: '', // Default child route
        name: 'SuperUserDashboardOverview',
        component: SuperUserDashboardOverview,
        meta: { requiresAuth: true, role: 'superuser' },
      },
      {
        path: '/companies',
        name: 'SuperUserCompanies',
        component: SuperUserCompanies,
        meta: { requiresAuth: true, role: 'superuser' },
      },
      {
        path: '/subscriptions',
        name: 'SuperUserSubscriptions',
        component: SuperUserSubscriptions,
        meta: { requiresAuth: true, role: 'superuser' },
      },
      {
        path: '/revenue-chart',
        name: 'SuperUserRevenueChart',
        component: SuperUserRevenueChart,
        meta: { requiresAuth: true, role: 'superuser' },
      },
      {
        path: '/subscription-packages',
        name: 'SuperUserSubscriptionPackages',
        component: SuperUserSubscriptionPackages,
        meta: { requiresAuth: true, role: 'superuser' },
      },
    ],
    meta: { requiresAuth: true, role: 'superuser' },
  },
  // Add more superuser-specific routes here
];

const router = createRouter({
  history: createWebHistory(),
  routes: superuserRoutes,
});

export default router;
