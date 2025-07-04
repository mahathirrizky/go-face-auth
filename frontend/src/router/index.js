import { createRouter, createWebHistory } from 'vue-router';
import LandingPage from '../components/main/LandingPage.vue';
import RegisterCompany from '../components/admin/RegisterCompany.vue';
import PaymentPage from '../components/admin/PaymentPage.vue';
import PaymentStatusPage from '../components/admin/PaymentStatusPage.vue';
import ForgotPassword from '../components/auth/ForgotPassword.vue'; // New component
import ResetPassword from '../components/auth/ResetPassword.vue';   // New component

const routes = [
  {
    path: '/',
    name: 'LandingPage',
    component: LandingPage,
  },
  {
    path: '/register/:packageId',
    name: 'RegisterCompany',
    component: RegisterCompany,
    props: true,
  },
  {
    path: '/checkout/:companyId',
    name: 'PaymentPage',
    component: PaymentPage,
    props: true,
  },
  {
    path: '/payment/finish',
    name: 'PaymentFinish',
    component: PaymentStatusPage,
    props: (route) => ({ order_id: route.query.order_id }),
  },
  {
    path: '/payment/error',
    name: 'PaymentError',
    component: PaymentStatusPage,
    props: (route) => ({ order_id: route.query.order_id }),
  },
  {
    path: '/payment/pending',
    name: 'PaymentPending',
    component: PaymentStatusPage,
    props: (route) => ({ order_id: route.query.order_id }),
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
  // Add other routes here as needed
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;