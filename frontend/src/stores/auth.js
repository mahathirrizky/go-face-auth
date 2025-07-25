import { defineStore } from 'pinia';
import axios from 'axios';

export const useAuthStore = defineStore('auth', {
  state: () => {
    console.log("AuthStore: Initializing state.");
    return {
      user: null,
      token: null,
      companyId: null,
      companyName: null,
      companyAddress: null,
      adminEmail: null,
      subscriptionStatus: null,
      trialEndDate: null,
      companyTimezone: null,
      hasConfiguredTimezone: false, // Added hasConfiguredTimezone
      subscriptionPackageId: null,
      billingCycle: null,
    };
  },
  actions: {
    setAuth(user, token) {
      this.user = user;
      this.token = token;
      console.log("AuthStore: setAuth called. Token set.");
    },
    async fetchCompanyDetails() {
      console.log("AuthStore: fetchCompanyDetails called. Current token:", this.token);
      if (!this.token) {
        console.log("AuthStore: Skipping fetchCompanyDetails. Token missing.");
        return;
      }
      try {
        const response = await axios.get('/api/company-details', {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        const { id, name, address, admin_email, subscription_status, trial_end_date, timezone, subscription_package_id, billing_cycle } = response.data.data;
        this.companyId = id; // Set companyId from the response
        this.companyName = name;
        this.companyAddress = address;
        this.adminEmail = admin_email;
        this.subscriptionStatus = subscription_status;
        this.trialEndDate = trial_end_date;
        this.companyTimezone = timezone; // Set companyTimezone
        this.subscriptionPackageId = subscription_package_id;
        this.billingCycle = billing_cycle;
        console.log("AuthStore: Company details fetched successfully.", { companyId: this.companyId, companyName: this.companyName, companyTimezone: this.companyTimezone });
      } catch (error) {
        console.error('AuthStore: Failed to fetch company details:', error);
        if (error.response && error.response.status === 403) {
          // If trial expired or subscription inactive, clear auth and reload
          this.clearAuth();
          window.location.reload();
        }
      }
    },
    clearAuth() {
      console.log("AuthStore: clearAuth called. Clearing all auth state.");
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
  persist: {
    key: 'auth-store',
    paths: ['user', 'token', 'companyId', 'companyName', 'companyAddress', 'adminEmail', 'subscriptionStatus', 'trialEndDate', 'companyTimezone', 'hasConfiguredTimezone'],
    afterRestore: (ctx) => {
      console.log("AuthStore: State rehydrated. companyId after restore:", ctx.store.companyId);
    },
  },
});