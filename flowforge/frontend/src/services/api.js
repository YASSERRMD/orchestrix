import axios from 'axios'

const API_URL = '/api'

const api = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' }
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const authService = {
  register: async (data) => {
    const res = await api.post('/auth/register', data)
    if (res.data.token) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
    }
    return res.data
  },
  login: async (email, password) => {
    const res = await api.post('/auth/login', { email, password })
    if (res.data.token) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
    }
    return res.data
  },
  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  },
  getToken: () => localStorage.getItem('token'),
  getUser: () => JSON.parse(localStorage.getItem('user') || '{}'),
  me: () => api.get('/auth/me')
}

export const userService = {
  list: () => api.get('/users'),
  get: (id) => api.get(`/users/${id}`),
  create: (data) => api.post('/users', data),
  update: (id, data) => api.put(`/users/${id}`, data),
  delete: (id) => api.delete(`/users/${id}`)
}

export const templateService = {
  list: () => api.get('/templates'),
  get: (id) => api.get(`/templates/${id}`),
  getVersions: (id) => api.get(`/templates/${id}/versions`),
  create: (data) => api.post('/templates', data),
  update: (id, data) => api.put(`/templates/${id}`, data)
}

export const workflowService = {
  list: (params) => api.get('/workflows', { params }),
  get: (id) => api.get(`/workflows/${id}`),
  create: (data) => api.post('/workflows', data),
  advance: (id, reason) => api.post(`/workflows/${id}/advance`, { reason }),
  reject: (id, reason) => api.post(`/workflows/${id}/reject`, { reason }),
  rollback: (id, reason) => api.post(`/workflows/${id}/rollback`, { reason }),
  cancel: (id, reason) => api.post(`/workflows/${id}/cancel`, { reason }),
  getAudit: (id) => api.get(`/workflows/${id}/audit`)
}

export default api
