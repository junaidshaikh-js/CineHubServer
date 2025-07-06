import { HomePage } from '../components/HomePage.js'
import { MovieDetailsPage } from '../components/MovieDetailsPage.js'
import { MoviesPage } from '../components/MoviesPage.js'

const ROUTES = [
  {
    path: '/',
    component: HomePage,
  },
  {
    path: /^\/movies\/(\d+)$/,
    component: MovieDetailsPage,
  },
  {
    path: '/movies',
    component: MoviesPage,
  },
]

export default ROUTES
