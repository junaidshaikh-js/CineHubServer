import { Router } from './services/Router.js'
import { API } from './services/API.js'
import Store from './services/Store.js'
import './components/YoutubeEmbed.js'

window.app = {
  Router,
  Store,
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
  register: async event => {
    event.preventDefault()
    const name = document.getElementById('register-name').value
    const email = document.getElementById('register-email').value
    const password = document.getElementById('register-password').value
    const confirmPassword = document.getElementById('register-confirm-password').value

    const errors = []

    if (name.length < 3) errors.push('Name must be at least 3 characters long')
    if (email.length < 5) errors.push('Email must be at least 5 characters long')
    if (!email.includes('@')) errors.push('Email must contain @')
    if (password.length < 8) errors.push('Password must be at least 8 characters long')
    if (confirmPassword.length < 8) errors.push('Confirm Password must be at least 8 characters long')
    if (password !== confirmPassword) errors.push('Passwords do not match')

    if (errors.length > 0) {
      app.showError(errors.join('. '), false)
      return
    }

    const result = await API.register(name, email, password)

    if (result.success) {
      app.Store.token = result.token
      app.Router.go('/account')
    } else {
      app.showError(result.message, false)
    }
  },
  login: async event => {
    event.preventDefault()
    const email = document.getElementById('login-email').value
    const password = document.getElementById('login-password').value

    const errors = []

    if (email.length < 5) errors.push('Email must be at least 5 characters long')
    if (!email.includes('@')) errors.push('Email must contain @')
    if (password.length < 8) errors.push('Password must be at least 8 characters long')

    if (errors.length > 0) {
      app.showError(errors.join('. '), false)
      return
    }

    const result = await API.login(email, password)

    if (result.success) {
      app.Store.token = result.token
      app.Router.go('/account')
    } else {
      app.showError(result.message, false)
    }
  },
  logout: () => {
    app.Store.token = null
    app.Router.go('/')
  },
}

window.addEventListener('DOMContentLoaded', event => {
  app.Router.init()
})
