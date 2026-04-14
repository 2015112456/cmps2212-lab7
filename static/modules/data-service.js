import { emitter } from './event-emitter.js'

const BASE = '/api'

export const DataService = {
  async registerEvent(payload) {
    try {
      const res = await fetch(`${BASE}/events`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify(payload),
      })
      if (!res.ok) {
        const err = await res.text()
        throw new Error(err || 'Failed to register event')
      }
      const event = await res.json()
      console.log('Event registered:', event)
    } catch (err) {
      console.error(`Error registering event: ${err.message}`)
      emitter.emit('events:error', err.message)
    }
  },
}
