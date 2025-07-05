import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    companyId: null,
  }),
  actions: {
    setAuth(user, token) {
      this.user = user;
      this.token = token;
    },
    setCompanyId(companyId) {
      this.companyId = companyId;
    },
    clearAuth() {
      this.user = null;
      this.token = null;
      this.companyId = null;
    },
  },
  persist: true,
});