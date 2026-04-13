import { emitter }     from './modules/event-emitter.js'
import { state }       from './state.js'
import { DataService } from './modules/data-service.js'
import { render }      from './render.js'


emitter.on('users:submit', (payload) => {
  render()
  DataService.createUser(payload)
})

// Boot 
render()     // initial paint — empty form
