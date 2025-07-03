import { HomePage } from './components/HomePage.js'
import { MovieDetailsPage } from './components/MovieDetailsPage.js'

window.addEventListener('DOMContentLoaded', event => {
  document.querySelector('main').appendChild(new HomePage())

  document.querySelector('main').appendChild(new MovieDetailsPage())
})

window.app = {
  search: event => {
    event.preventDefault()
    const q = document.querySelector('input[name="search"]').value
    console.log(q)
  },
}
