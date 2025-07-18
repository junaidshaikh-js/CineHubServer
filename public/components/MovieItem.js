export class MovieItem extends HTMLElement {
  constructor(movie) {
    super()
    this.movie = movie
  }

  connectedCallback() {
    const url = `/movies/${this.movie.id}`

    this.innerHTML = `
      <a onClick="app.Router.go('${url}')">
        <article class="movie-item">
          <img 
            src="${this.movie.poster_url}" 
            alt="${this.movie.title} Poster" 
            class="movie-item__poster"
            />
          <p class="movie-item__title">
            ${this.movie.title} (${this.movie.release_year})
          </p>
        </article>
      </a>
    `
  }
}

customElements.define('movie-item', MovieItem)
