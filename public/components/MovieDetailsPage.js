import { API } from '../services/API.js'

export class MovieDetailsPage extends HTMLElement {
  id = null
  movie = null

  async render() {
    try {
      this.movie = await API.getMovieById(this.id)
    } catch (error) {
      // TODO: show error message
      alert('Movie not found')
      return
    }

    const template = document.getElementById('template-movie-details')
    const content = template.content.cloneNode(true)
    this.appendChild(content)

    this.querySelector('h2').textContent = this.movie.title
    this.querySelector('h3').textContent = this.movie.tagline
    this.querySelector('img').src = this.movie.poster_url
    this.querySelector('youtube-embed').dataset.url = this.movie.trailer_url
    this.querySelector('#overview').textContent = this.movie.overview
    this.querySelector('#metadata').innerHTML = `
      <dt>Release Date</dt>
      <dd>${this.movie.release_year}</dd>
      <dt>Score</dt>
      <dd>${this.movie.score}</dd>
      <dt>Popularity</dt>
      <dd>${this.movie.popularity}</dd>
    `

    const ulGenres = this.querySelector('#genres')
    ulGenres.innerHTML = ''
    this.movie.genres.forEach(genre => {
      const li = document.createElement('li')
      li.textContent = genre.name
      ulGenres.appendChild(li)
    })

    const ulCast = this.querySelector('#cast')
    ulCast.innerHTML = ''
    this.movie.casting.forEach(actor => {
      const li = document.createElement('li')
      li.innerHTML = `
        <img src="${actor.image_url ?? '/images/generic_actor.jpg'}" alt="Picture of ${actor.last_name}">
        <p>${actor.first_name} ${actor.last_name}</p>
      `
      ulCast.appendChild(li)
    })
  }

  connectedCallback() {
    this.id = this.params[0]
    this.render()
  }
}

customElements.define('movie-details-page', MovieDetailsPage)
