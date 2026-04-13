import { emitter } from './event-emitter.js'

const BASE = '/api'

export const DataService = {
  async createUser(payload) {
    try {
      const res = await fetch(`${BASE}/users`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify(payload),
      })
      if (!res.ok) {
        throw new Error(`${res.status} ${res.statusText}`)
      }
      const user = await res.json()
      console.log('User created:', user)
    } catch (err) {
      console.error(`Error creating user: ${err.message}`)
    }
  },
}
