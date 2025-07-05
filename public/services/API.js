export const API = {
  BASE_URL: '/api',
  fetch: async (url, args) => {
    try {
      const query = args ? `?${new URLSearchParams(args).toString()}` : ''
      const res = await fetch(`${API.BASE_URL}${url}${query}`)
      const data = await res.json()
      return data
    } catch (error) {
      console.error(error)
      app.showError()
    }
  },
  getTopMovies: async () => {
    return await API.fetch('/movies/top')
  },
  getRandomMovies: async () => {
    return await API.fetch('/movies/random')
  },
  getMovieById: async id => {
    return await API.fetch(`/movies/${id}`)
  },
  searchMovies: async (q, order, genre) => {
    return await API.fetch(`/movies/search`, { q, order, genre })
  },
  getGenres: async () => {
    return await API.fetch('/genres')
  },
}
