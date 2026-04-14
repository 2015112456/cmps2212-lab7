import { state }   from './state.js'
import { emitter } from './modules/event-emitter.js'


function handleSubmit(e) {
  e.preventDefault()
  
  const data = new FormData(e.target)
  const payload = {
    date:  data.get('date').trim(),
    tickets: data.get('tickets').trim(),
    terms: data.get('terms') === 'on'
  }

  emitter.emit('events:submit', payload)
}

function renderForm() {
  let errorHtml = ''
  if (state.error) {
    errorHtml = `<div class="error-message">⚠ ${state.error}</div>`
  }
  return `
    <div class="form-card">
      <h2>Event Registration</h2>
      <form id="event-form" class="ui-form">

        <input
          class="ui-input"
          name="date"
          type="date"
          placeholder="Event Date (e.g., 2026-12-31)"
          required
        />

        <input
          class="ui-input"
          name="tickets"
          type="number"
          placeholder="Number of Tickets (1-5)"
          required
        />

        <label for="terms">
          <input
            class="ui-checkbox"
            name="terms"
            type="checkbox"
          />
          I agree to the Terms and Conditions
        </label>
        ${errorHtml}
        <button class="ui-btn" type="submit">Submit</button>
      </form>
    </div>`
}

export function render() {
 
  const app = document.querySelector('#app')
  if(app) {
     app.innerHTML = renderForm()
  }

  const form = document.querySelector('#event-form')
  if(form) {
    form.addEventListener('submit', handleSubmit)
  }
}