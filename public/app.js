import { HomePage } from './components/HomePage.js'

window.addEventListener('DOMContentLoaded', event => {
  document.querySelector('main').appendChild(new HomePage())
})

window.app = {
  search: event => {
    event.preventDefault()
    const q = document.querySelector('input[name="search"]').value
    console.log(q)
  },
}
