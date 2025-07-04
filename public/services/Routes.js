import { HomePage } from '../components/HomePage.js'
import { MovieDetailsPage } from '../components/MovieDetailsPage.js'

const ROUTES = [
  {
    path: '/',
    component: HomePage,
  },
  {
    path: /^\/movies\/(\d+)$/,
    component: MovieDetailsPage,
  },
]

export default ROUTES
