import { emitter }     from './modules/event-emitter.js'
import { state }       from './state.js'
import { DataService } from './modules/data-service.js'
import { render }      from './render.js'


emitter.on('events:submit', (payload) => {
  state.error = null
  render()
  DataService.registerEvent(payload)
})

emitter.on('events:error', (message) => {
  state.error = message
  render()
})

// Boot 
render()     // initial paint — empty form
