import { createRouter, createWebHistory } from 'vue-router';
import LandingPage from '../components/main/LandingPage.vue';
import RegisterCompany from '../components/admin/RegisterCompany.vue';
import PaymentPage from '../components/admin/PaymentPage.vue';
import PaymentStatusPage from '../components/admin/PaymentStatusPage.vue';
import ForgotPassword from '../components/auth/ForgotPassword.vue'; // New component
import ResetPassword from '../components/auth/ResetPassword.vue';   // New component
import InitialPasswordSetup from '../components/auth/InitialPasswordSetup.vue'; // New component
import EmployeeResetPassword from '../components/auth/EmployeeResetPassword.vue'; // New component
import InitialPasswordSuccess from '../components/auth/InitialPasswordSuccess.vue'; // New component
import NotFound from '../components/main/NotFound.vue'; // New 404 component
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
  {
    path: '/initial-password-setup',
    name: 'InitialPasswordSetup',
    component: InitialPasswordSetup,
  },
  {
    path: '/initial-password-success',
    name: 'InitialPasswordSuccess',
    component: InitialPasswordSuccess,
  },
  {
    path: '/employee-reset-password',
    name: 'EmployeeResetPassword',
    component: EmployeeResetPassword,
  },
  // Catch-all 404 route
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
  },
  // Add other routes here as needed
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;