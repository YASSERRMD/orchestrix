<template>
  <div class="templates-page">
    <div class="page-header">
      <h1>Workflow Templates</h1>
      <button @click="showCreate = true" class="btn-primary">+ New Template</button>
    </div>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="templates.length === 0" class="empty">No templates found</div>
    <div v-else class="template-list">
      <div v-for="template in templates" :key="template.id" class="template-card">
        <div class="template-header">
          <h3>{{ template.name }}</h3>
          <span class="version">v{{ template.version }}</span>
        </div>
        <p class="description">{{ template.description || 'No description' }}</p>
        <div class="stages">
          <span v-for="(stage, idx) in template.stages" :key="stage.id" class="stage-badge">
            {{ idx + 1 }}. {{ stage.name }}
          </span>
        </div>
      </div>
    </div>

    <div v-if="showCreate" class="modal">
      <div class="modal-content">
        <h2>Create Template</h2>
        <div class="form-group">
          <label>Name</label>
          <input v-model="form.name" type="text" />
        </div>
        <div class="form-group">
          <label>Description</label>
          <textarea v-model="form.description"></textarea>
        </div>
        <div class="form-group">
          <label>Stages</label>
          <div v-for="(stage, idx) in form.stages" :key="idx" class="stage-form">
            <input v-model="stage.name" placeholder="Stage name" />
            <select v-model="stage.required_role">
              <option value="viewer">Viewer</option>
              <option value="operator">Operator</option>
              <option value="manager">Manager</option>
              <option value="admin">Admin</option>
            </select>
            <select v-model="stage.approval_type">
              <option value="manual">Manual</option>
              <option value="auto">Auto</option>
            </select>
            <input v-model.number="stage.timeout_minutes" type="number" placeholder="Timeout (min)" />
            <button @click="removeStage(idx)" class="btn-remove">×</button>
          </div>
          <button @click="addStage" class="btn-add">+ Add Stage</button>
        </div>
        <div class="modal-actions">
          <button @click="showCreate = false" class="btn-secondary">Cancel</button>
          <button @click="createTemplate" class="btn-primary">Create</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { templateService } from '../services/api'

export default {
  name: 'Templates',
  setup() {
    const templates = ref([])
    const loading = ref(false)
    const showCreate = ref(false)
    const form = ref({
      name: '',
      description: '',
      stages: [{ name: '', required_role: 'operator', approval_type: 'manual', timeout_minutes: 0 }]
    })

    const loadTemplates = async () => {
      loading.value = true
      try {
        const res = await templateService.list()
        templates.value = res.data
      } catch (e) {
        console.error(e)
      } finally {
        loading.value = false
      }
    }

    const addStage = () => {
      form.value.stages.push({ name: '', required_role: 'operator', approval_type: 'manual', timeout_minutes: 0 })
    }

    const removeStage = (idx) => {
      form.value.stages.splice(idx, 1)
    }

    const createTemplate = async () => {
      try {
        await templateService.create(form.value)
        showCreate.value = false
        form.value = { name: '', description: '', stages: [{ name: '', required_role: 'operator', approval_type: 'manual', timeout_minutes: 0 }] }
        loadTemplates()
      } catch (e) {
        alert(e.response?.data?.error || 'Failed to create template')
      }
    }

    onMounted(loadTemplates)

    return { templates, loading, showCreate, form, addStage, removeStage, createTemplate }
  }
}
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.btn-primary { background: #007bff; color: #fff; border: none; padding: 0.75rem 1.5rem; border-radius: 4px; cursor: pointer; }
.template-list { display: grid; gap: 1.5rem; }
.template-card { background: #fff; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
.template-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.version { background: #e9ecef; padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.8rem; }
.description { color: #666; margin-bottom: 1rem; }
.stages { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.stage-badge { background: #f8f9fa; padding: 0.25rem 0.75rem; border-radius: 12px; font-size: 0.85rem; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: #fff; padding: 2rem; border-radius: 8px; width: 600px; max-height: 80vh; overflow-y: auto; }
.modal-actions { display: flex; gap: 1rem; justify-content: flex-end; margin-top: 1.5rem; }
.btn-secondary { background: #6c757d; color: #fff; border: none; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; }
.form-group input, .form-group textarea, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.stage-form { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr auto; gap: 0.5rem; margin-bottom: 0.5rem; }
.btn-add { background: #28a745; color: #fff; border: none; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; }
.btn-remove { background: #dc3545; color: #fff; border: none; padding: 0.5rem; border-radius: 4px; cursor: pointer; }
</style>
