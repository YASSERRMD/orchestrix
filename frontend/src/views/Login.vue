<template>
  <div class="login-page">
    <div class="login-card">
      <h1>FlowForge</h1>
      <p class="subtitle">Workflow & Approvals Engine</p>
      
      <div v-if="isRegister" class="form-group">
        <label>Organization Name</label>
        <input v-model="form.org_name" type="text" placeholder="Enter org name" />
      </div>
      
      <div class="form-group">
        <label>Email</label>
        <input v-model="form.email" type="email" placeholder="Enter email" />
      </div>
      
      <div class="form-group">
        <label>Password</label>
        <input v-model="form.password" type="password" placeholder="Enter password" />
      </div>
      
      <div v-if="isRegister" class="form-group">
        <label>Name</label>
        <input v-model="form.name" type="text" placeholder="Your name" />
      </div>
      
      <div v-if="isRegister" class="form-group">
        <label>Role</label>
        <select v-model="form.role">
          <option value="admin">Admin</option>
          <option value="manager">Manager</option>
          <option value="operator">Operator</option>
          <option value="viewer">Viewer</option>
        </select>
      </div>
      
      <button @click="submit" class="btn-primary" :disabled="loading">
        {{ loading ? 'Processing...' : (isRegister ? 'Register' : 'Login') }}
      </button>
      
      <p class="toggle-link">
        {{ isRegister ? 'Already have an account?' : "Don't have an account?" }}
        <a href="#" @click.prevent="isRegister = !isRegister">
          {{ isRegister ? 'Login' : 'Register' }}
        </a>
      </p>
      
      <p v-if="error" class="error">{{ error }}</p>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../services/api'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const isRegister = ref(false)
    const loading = ref(false)
    const error = ref('')
    const form = reactive({
      org_name: '',
      email: '',
      password: '',
      name: '',
      role: 'operator'
    })

    const submit = async () => {
      loading.value = true
      error.value = ''
      try {
        if (isRegister.value) {
          await authService.register(form)
        } else {
          await authService.login(form.email, form.password)
        }
        router.push('/dashboard')
      } catch (e) {
        error.value = e.response?.data?.error || 'An error occurred'
      } finally {
        loading.value = false
      }
    }

    return { isRegister, loading, error, form, submit }
  }
}
</script>

<style scoped>
.login-page { display: flex; justify-content: center; align-items: center; min-height: 80vh; }
.login-card { background: #fff; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 12px rgba(0,0,0,0.1); width: 400px; }
.login-card h1 { text-align: center; color: #333; margin-bottom: 0.5rem; }
.subtitle { text-align: center; color: #666; margin-bottom: 2rem; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; color: #555; font-size: 0.9rem; }
.form-group input, .form-group select { width: 100%; padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; font-size: 1rem; }
.btn-primary { width: 100%; padding: 0.75rem; background: #007bff; color: #fff; border: none; border-radius: 4px; cursor: pointer; font-size: 1rem; }
.btn-primary:disabled { background: #ccc; }
.toggle-link { text-align: center; margin-top: 1rem; color: #666; }
.error { color: #dc3545; margin-top: 1rem; text-align: center; }
</style>
