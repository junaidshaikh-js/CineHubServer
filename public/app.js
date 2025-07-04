import { Router } from './services/Router.js'
import './components/YoutubeEmbed.js'

window.app = {
  search: event => {
    event.preventDefault()
    const q = document.querySelector('input[name="search"]').value
  },
  Router,
}

window.addEventListener('DOMContentLoaded', event => {
  app.Router.init()
})
