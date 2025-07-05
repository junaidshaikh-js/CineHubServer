import { Router } from './services/Router.js'
import './components/YoutubeEmbed.js'

window.app = {
  Router,
  showError: (message = 'Something went wrong!', goToHome = true) => {
    document.getElementById('alert-modal').showModal()
    document.querySelector('#alert-modal p').textContent = message
    if (goToHome) app.Router.go('/')
  },
  closeError: () => {
    document.getElementById('alert-modal').close()
  },
  search: event => {
    event.preventDefault()
    const q = document.querySelector('input[name="search"]').value
  },
}

window.addEventListener('DOMContentLoaded', event => {
  app.Router.init()
})
