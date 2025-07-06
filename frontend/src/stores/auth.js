import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    companyId: null,
    companyName: null,
    companyAddress: null,
    adminEmail: null,
  }),
  actions: {
    setAuth(user, token) {
      this.user = user;
      this.token = token;
    },
    setCompanyId(companyId) {
      this.companyId = companyId;
    },
    setCompanyName(companyName) {
      this.companyName = companyName;
    },
    setCompanyAddress(companyAddress) {
      this.companyAddress = companyAddress;
    },
    setAdminEmail(adminEmail) {
      this.adminEmail = adminEmail;
    },
    clearAuth() {
      this.user = null;
      this.token = null;
      this.companyId = null;
      this.companyName = null;
      this.companyAddress = null;
      this.adminEmail = null;
    },
  },
  persist: true,
});