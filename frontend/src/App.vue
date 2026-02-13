<template>
  <div id="app">
    <nav v-if="isAuthenticated" class="navbar">
      <div class="nav-brand">FlowForge</div>
      <div class="nav-links">
        <router-link to="/dashboard">Dashboard</router-link>
        <router-link to="/templates">Templates</router-link>
        <button @click="logout" class="btn-logout">Logout</button>
      </div>
    </nav>
    <div class="container">
      <router-view />
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from './services/api'

export default {
  name: 'App',
  setup() {
    const router = useRouter()
    const isAuthenticated = computed(() => !!authService.getToken())

    const logout = () => {
      authService.logout()
      router.push('/login')
    }

    return { isAuthenticated, logout }
  }
}
</script>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; }
.navbar { background: #fff; padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
.nav-brand { font-size: 1.5rem; font-weight: bold; color: #333; }
.nav-links { display: flex; gap: 1.5rem; align-items: center; }
.nav-links a { color: #666; text-decoration: none; }
.nav-links a.router-link-active { color: #007bff; }
.btn-logout { background: none; border: 1px solid #ddd; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; }
.container { max-width: 1200px; margin: 2rem auto; padding: 0 1rem; }
</style>
