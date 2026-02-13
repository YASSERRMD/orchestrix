<template>
  <div class="workflow-detail">
    <button @click="$router.back()" class="btn-back">← Back</button>
    
    <div v-if="loading" class="loading">Loading...</div>
    <template v-else-if="workflow">
      <div class="workflow-header">
        <h1>{{ workflow.title }}</h1>
        <span :class="['status', workflow.status]">{{ workflow.status }}</span>
      </div>
      
      <div class="workflow-info">
        <div class="info-card">
          <h3>Current Stage</h3>
          <p>{{ workflow.current_stage?.name || 'N/A' }}</p>
        </div>
        <div class="info-card">
          <h3>Template</h3>
          <p>{{ workflow.template?.name }} (v{{ workflow.template?.version }})</p>
        </div>
        <div class="info-card">
          <h3>Created By</h3>
          <p>{{ workflow.created_by_user?.name || 'N/A' }}</p>
        </div>
      </div>

      <div class="actions" v-if="canAct">
        <button @click="handleAdvance" class="btn-success" v-if="canAdvance">Advance</button>
        <button @click="handleReject" class="btn-danger" v-if="canReject">Reject</button>
        <button @click="handleRollback" class="btn-warning" v-if="canRollback">Rollback</button>
        <button @click="handleCancel" class="btn-secondary" v-if="canCancel">Cancel</button>
      </div>

      <div class="timeline">
        <h3>Timeline</h3>
        <div v-for="log in auditLogs" :key="log.id" class="timeline-item">
          <div class="timeline-marker"></div>
          <div class="timeline-content">
            <div class="timeline-action">{{ log.action }}</div>
            <div class="timeline-meta">
              <span>{{ log.actor?.name || 'System' }}</span>
              <span>{{ log.actor_role }}</span>
              <span>{{ formatDate(log.timestamp) }}</span>
            </div>
            <div v-if="log.reason" class="timeline-reason">{{ log.reason }}</div>
          </div>
        </div>
      </div>
    </template>

    <div v-if="showActionModal" class="modal">
      <div class="modal-content">
        <h2>{{ actionTitle }}</h2>
        <div class="form-group">
          <label>Reason (optional)</label>
          <textarea v-model="actionReason"></textarea>
        </div>
        <div class="modal-actions">
          <button @click="showActionModal = false" class="btn-secondary">Cancel</button>
          <button @click="confirmAction" class="btn-primary">Confirm</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { workflowService } from '../services/api'
import { wsService } from '../services/websocket'

export default {
  name: 'WorkflowDetail',
  setup() {
    const route = useRoute()
    const workflow = ref(null)
    const auditLogs = ref([])
    const loading = ref(false)
    const showActionModal = ref(false)
    const actionType = ref('')
    const actionReason = ref('')

    const user = computed(() => JSON.parse(localStorage.getItem('user') || '{}'))

    const canAct = computed(() => workflow.value && ['pending', 'escalated'].includes(workflow.value.status))
    const canAdvance = computed(() => canAct.value)
    const canReject = computed(() => canAct.value)
    const canRollback = computed(() => canAct.value && ['admin', 'manager'].includes(user.value.role))
    const canCancel = computed(() => canAct.value || workflow.value?.created_by === user.value.id)

    const actionTitle = computed(() => {
      const titles = { advance: 'Advance Workflow', reject: 'Reject Workflow', rollback: 'Rollback Workflow', cancel: 'Cancel Workflow' }
      return titles[actionType.value] || ''
    })

    const loadWorkflow = async () => {
      loading.value = true
      try {
        const res = await workflowService.get(route.params.id)
        workflow.value = res.data
      } catch (e) {
        console.error(e)
      } finally {
        loading.value = false
      }
    }

    const loadAuditLogs = async () => {
      try {
        const res = await workflowService.getAudit(route.params.id)
        auditLogs.value = res.data
      } catch (e) {
        console.error(e)
      }
    }

    const handleAdvance = () => { actionType.value = 'advance'; showActionModal.value = true }
    const handleReject = () => { actionType.value = 'reject'; showActionModal.value = true }
    const handleRollback = () => { actionType.value = 'rollback'; showActionModal.value = true }
    const handleCancel = () => { actionType.value = 'cancel'; showActionModal.value = true }

    const confirmAction = async () => {
      try {
        const actions = { advance: workflowService.advance, reject: workflowService.reject, rollback: workflowService.rollback, cancel: workflowService.cancel }
        await actions[actionType.value](workflow.value.id, actionReason.value)
        showActionModal.value = false
        actionReason.value = ''
        loadWorkflow()
        loadAuditLogs()
      } catch (e) {
        alert(e.response?.data?.error || 'Action failed')
      }
    }

    const formatDate = (date) => new Date(date).toLocaleString()

    onMounted(() => {
      loadWorkflow()
      loadAuditLogs()
      wsService.connect()
      wsService.on('stage_advanced', loadWorkflow)
      wsService.on('rejected', loadWorkflow)
      wsService.on('escalated', loadWorkflow)
    })

    return { workflow, auditLogs, loading, canAct, canAdvance, canReject, canRollback, canCancel, showActionModal, actionType, actionReason, actionTitle, handleAdvance, handleReject, handleRollback, handleCancel, confirmAction, formatDate }
  }
}
</script>

<style scoped>
.btn-back { background: none; border: none; color: #666; cursor: pointer; margin-bottom: 1rem; }
.workflow-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.status { padding: 0.5rem 1rem; border-radius: 4px; font-weight: 500; }
.status.pending { background: #fff3cd; color: #856404; }
.status.completed { background: #d4edda; color: #155724; }
.status.escalated { background: #f8d7da; color: #721c24; }
.status.rejected { background: #f8d7da; color: #721c24; }
.status.canceled { background: #e2e3e5; color: #383d41; }
.workflow-info { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin-bottom: 2rem; }
.info-card { background: #fff; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
.info-card h3 { color: #666; font-size: 0.85rem; margin-bottom: 0.5rem; }
.actions { display: flex; gap: 1rem; margin-bottom: 2rem; }
.btn-success { background: #28a745; color: #fff; border: none; padding: 0.75rem 1.5rem; border-radius: 4px; cursor: pointer; }
.btn-danger { background: #dc3545; color: #fff; border: none; padding: 0.75rem 1.5rem; border-radius: 4px; cursor: pointer; }
.btn-warning { background: #ffc107; color: #000; border: none; padding: 0.75rem 1.5rem; border-radius: 4px; cursor: pointer; }
.btn-secondary { background: #6c757d; color: #fff; border: none; padding: 0.75rem 1.5rem; border-radius: 4px; cursor: pointer; }
.timeline { background: #fff; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
.timeline h3 { margin-bottom: 1rem; }
.timeline-item { display: flex; gap: 1rem; padding: 1rem 0; border-bottom: 1px solid #eee; }
.timeline-marker { width: 12px; height: 12px; border-radius: 50%; background: #007bff; margin-top: 4px; }
.timeline-content { flex: 1; }
.timeline-action { font-weight: 500; text-transform: capitalize; }
.timeline-meta { display: flex; gap: 1rem; color: #666; font-size: 0.85rem; margin-top: 0.25rem; }
.timeline-reason { color: #666; margin-top: 0.5rem; font-style: italic; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: #fff; padding: 2rem; border-radius: 8px; width: 400px; }
.modal-actions { display: flex; gap: 1rem; justify-content: flex-end; margin-top: 1.5rem; }
.btn-primary { background: #007bff; color: #fff; border: none; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; }
.form-group textarea { width: 100%; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
</style>
