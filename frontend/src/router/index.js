import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  {
    path: '/',
    name: 'LandingPage',
    component: () => import('../components/main/LandingPage.vue'),
  },
  {
    path: '/register/:packageId',
    name: 'RegisterCompany',
    component: () => import('../components/admin/RegisterCompany.vue'),
    props: true,
  },
  {
    path: '/checkout/:companyId',
    name: 'PaymentPage',
    component: () => import('../components/admin/PaymentPage.vue'),
    props: true,
  },
  {
    path: '/payment/finish',
    name: 'PaymentFinish',
    component: () => import('../components/admin/PaymentStatusPage.vue'),
    props: (route) => ({ order_id: route.query.order_id }),
  },
  {
    path: '/payment/error',
    name: 'PaymentError',
    component: () => import('../components/admin/PaymentStatusPage.vue'),
    props: (route) => ({ order_id: route.query.order_id }),
  },
  {
    path: '/payment/pending',
    name: 'PaymentPending',
    component: () => import('../components/admin/PaymentStatusPage.vue'),
    props: (route) => ({ order_id: route.query.order_id }),
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
    path: '/initial-password-setup',
    name: 'InitialPasswordSetup',
    component: () => import('../components/auth/InitialPasswordSetup.vue'),
  },
  {
    path: '/initial-password-success',
    name: 'InitialPasswordSuccess',
    component: () => import('../components/auth/InitialPasswordSuccess.vue'),
  },
  {
    path: '/employee-reset-password',
    name: 'EmployeeResetPassword',
    component: () => import('../components/auth/EmployeeResetPassword.vue'),
  },
  // Catch-all 404 route
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../components/main/NotFound.vue'),
  },
  // Add other routes here as needed
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;