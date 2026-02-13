<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1>Workflow Dashboard</h1>
      <button @click="showCreate = true" class="btn-primary" v-if="canCreate">+ New Workflow</button>
    </div>

    <div class="tabs">
      <button :class="{ active: activeTab === 'active' }" @click="activeTab = 'active'">Active</button>
      <button :class="{ active: activeTab === 'completed' }" @click="activeTab = 'completed'">Completed</button>
      <button :class="{ active: activeTab === 'escalated' }" @click="activeTab = 'escalated'">Escalated</button>
      <button :class="{ active: activeTab === 'my' }" @click="activeTab = 'my'">My Approvals</button>
    </div>

    <div class="workflow-list">
      <div v-if="loading" class="loading">Loading...</div>
      <div v-else-if="workflows.length === 0" class="empty">No workflows found</div>
      <div v-else v-for="workflow in workflows" :key="workflow.id" class="workflow-card" @click="viewWorkflow(workflow.id)">
        <div class="workflow-title">{{ workflow.title }}</div>
        <div class="workflow-meta">
          <span :class="['status', workflow.status]">{{ workflow.status }}</span>
          <span>Template: {{ workflow.template?.name || 'N/A' }}</span>
        </div>
      </div>
    </div>

    <div v-if="showCreate" class="modal">
      <div class="modal-content">
        <h2>Create Workflow</h2>
        <div class="form-group">
          <label>Template</label>
          <select v-model="createForm.template_id">
            <option v-for="t in templates" :key="t.id" :value="t.id">{{ t.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>Title</label>
          <input v-model="createForm.title" type="text" placeholder="Workflow title" />
        </div>
        <div class="modal-actions">
          <button @click="showCreate = false" class="btn-secondary">Cancel</button>
          <button @click="createWorkflow" class="btn-primary">Create</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { workflowService, templateService } from '../services/api'
import { wsService } from '../services/websocket'

export default {
  name: 'Dashboard',
  setup() {
    const router = useRouter()
    const activeTab = ref('active')
    const workflows = ref([])
    const templates = ref([])
    const loading = ref(false)
    const showCreate = ref(false)
    const createForm = ref({ template_id: null, title: '' })

    const user = computed(() => JSON.parse(localStorage.getItem('user') || '{}'))
    const canCreate = computed(() => ['admin', 'manager', 'operator'].includes(user.value.role))

    const loadWorkflows = async () => {
      loading.value = true
      try {
        const params = {}
        if (activeTab.value === 'active') params.status = 'pending'
        else if (activeTab.value === 'completed') params.status = 'completed'
        else if (activeTab.value === 'escalated') params.status = 'escalated'
        
        const res = await workflowService.list(params)
        workflows.value = res.data
      } catch (e) {
        console.error(e)
      } finally {
        loading.value = false
      }
    }

    const loadTemplates = async () => {
      try {
        const res = await templateService.list()
        templates.value = res.data
        if (res.data.length) createForm.value.template_id = res.data[0].id
      } catch (e) {
        console.error(e)
      }
    }

    const createWorkflow = async () => {
      try {
        await workflowService.create(createForm.value)
        showCreate.value = false
        loadWorkflows()
      } catch (e) {
        alert(e.response?.data?.error || 'Failed to create workflow')
      }
    }

    const viewWorkflow = (id) => router.push(`/workflow/${id}`)

    onMounted(() => {
      loadWorkflows()
      loadTemplates()
      wsService.connect()
      wsService.on('stage_advanced', loadWorkflows)
      wsService.on('workflow_created', loadWorkflows)
      wsService.on('escalated', loadWorkflows)
    })

    onUnmounted(() => {
      wsService.off('stage_advanced', loadWorkflows)
      wsService.off('workflow_created', loadWorkflows)
      wsService.off('escalated', loadWorkflows)
    })

    return { activeTab, workflows, templates, loading, showCreate, createForm, canCreate, createWorkflow, viewWorkflow }
  }
}
</script>

<style scoped>
.dashboard-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.btn-primary { background: #007bff; color: #fff; border: none; padding: 0.75rem 1.5rem; border-radius: 4px; cursor: pointer; }
.tabs { display: flex; gap: 0.5rem; margin-bottom: 1.5rem; }
.tabs button { padding: 0.75rem 1.5rem; border: none; background: #e9ecef; cursor: pointer; border-radius: 4px; }
.tabs button.active { background: #007bff; color: #fff; }
.workflow-list { display: flex; flex-direction: column; gap: 1rem; }
.workflow-card { background: #fff; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); cursor: pointer; transition: transform 0.2s; }
.workflow-card:hover { transform: translateY(-2px); }
.workflow-title { font-size: 1.1rem; font-weight: 600; margin-bottom: 0.5rem; }
.workflow-meta { display: flex; gap: 1rem; color: #666; font-size: 0.9rem; }
.status { padding: 0.25rem 0.75rem; border-radius: 12px; font-size: 0.8rem; font-weight: 500; }
.status.pending { background: #fff3cd; color: #856404; }
.status.completed { background: #d4edda; color: #155724; }
.status.escalated { background: #f8d7da; color: #721c24; }
.status.rejected { background: #f8d7da; color: #721c24; }
.status.canceled { background: #e2e3e5; color: #383d41; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: #fff; padding: 2rem; border-radius: 8px; width: 400px; }
.modal-actions { display: flex; gap: 1rem; justify-content: flex-end; margin-top: 1.5rem; }
.btn-secondary { background: #6c757d; color: #fff; border: none; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; }
.form-group input, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
</style>
