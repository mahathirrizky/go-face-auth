import { defineStore } from 'pinia';
import axios from 'axios';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    companyId: null,
    companyName: null,
    companyAddress: null,
    adminEmail: null,
    subscriptionStatus: null,
    trialEndDate: null,
  }),
  actions: {
    setAuth(user, token) {
      this.user = user;
      this.token = token;
    },
    setCompanyId(companyId) {
      this.companyId = companyId;
    },
    async fetchCompanyDetails() {
      if (!this.token) return;
      try {
        const response = await axios.get('/api/company-details', {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        const { name, address, admin_email, subscription_status, trial_end_date } = response.data.data;
        this.companyName = name;
        this.companyAddress = address;
        this.adminEmail = admin_email;
        this.subscriptionStatus = subscription_status;
        this.trialEndDate = trial_end_date;
      } catch (error) {
        console.error('Failed to fetch company details:', error);
        if (error.response && error.response.status === 403) {
          // If trial expired or subscription inactive, clear auth and reload
          this.clearAuth();
          window.location.reload();
        }
      }
    },
    clearAuth() {
      this.user = null;
      this.token = null;
      this.companyId = null;
      this.companyName = null;
      this.companyAddress = null;
      this.adminEmail = null;
      this.subscriptionStatus = null;
      this.trialEndDate = null;
    },
  },
  persist: true,
});