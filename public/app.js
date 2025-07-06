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
    app.Router.go(`/movies?q=${q}`)
  },
  searchOrderChange: order => {
    const urlParams = new URLSearchParams(window.location.search)
    const q = urlParams.get('q')
    const genre = urlParams.get('genre') ?? ''
    const genreQuery = genre ? `&genre=${genre}` : ''
    app.Router.go(`/movies?q=${q}&order=${order}${genreQuery}`)
  },
  searchFilterChange: genre => {
    const urlParams = new URLSearchParams(window.location.search)
    const q = urlParams.get('q')
    const order = urlParams.get('order') ?? ''
    const genreQuery = genre ? `&genre=${genre}` : ''
    const orderQuery = order ? `&order=${order}` : ''
    app.Router.go(`/movies?q=${q}${orderQuery}${genreQuery}`)
  },
}

window.addEventListener('DOMContentLoaded', event => {
  app.Router.init()
})
