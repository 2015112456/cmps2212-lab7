import { state }   from './state.js'
import { emitter } from './modules/event-emitter.js'


function handleSubmit(e) {
  e.preventDefault()
  
  const data = new FormData(e.target)
  const payload = {
    name:  data.get('name').trim(),
    email: data.get('email').trim(),
  }

  emitter.emit('users:submit', payload)
}

function renderForm() {
  return `
    <div class="form-card">
      <h2>Add a user</h2>
      <form id="user-form" class="ui-form">

        <input
          class="ui-input"
          name="name"
          placeholder="Full name"
          required
        />

        <input
          class="ui-input"
          name="email"
          type="email"
          placeholder="Email address"
          required
        />

        <button class="ui-btn" type="submit">Submit</button>
      </form>
    </div>`
}

export function render() {
 
  const app = document.querySelector('#app')
  if(app) {
     app.innerHTML = renderForm()
  }

  const form = document.querySelector('#user-form')
  if(form) {
    form.addEventListener('submit', handleSubmit)
  }
}