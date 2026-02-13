const WS_URL = `ws://${window.location.host}/ws`

class WebSocketService {
  constructor() {
    this.ws = null
    this.listeners = {}
  }

  connect() {
    this.ws = new WebSocket(WS_URL)
    
    this.ws.onopen = () => console.log('WebSocket connected')
    
    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        this.emit(message.type, message.payload)
      } catch (e) {
        console.error('WebSocket message error:', e)
      }
    }
    
    this.ws.onclose = () => {
      console.log('WebSocket disconnected, reconnecting...')
      setTimeout(() => this.connect(), 3000)
    }
    
    this.ws.onerror = (err) => console.error('WebSocket error:', err)
  }

  on(event, callback) {
    if (!this.listeners[event]) this.listeners[event] = []
    this.listeners[event].push(callback)
  }

  off(event, callback) {
    if (!this.listeners[event]) return
    this.listeners[event] = this.listeners[event].filter(cb => cb !== callback)
  }

  emit(event, data) {
    if (!this.listeners[event]) return
    this.listeners[event].forEach(cb => cb(data))
  }

  disconnect() {
    if (this.ws) this.ws.close()
  }
}

export const wsService = new WebSocketService()
export default wsService
