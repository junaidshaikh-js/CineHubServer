import { HomePage } from '../components/HomePage.js'
import { MovieDetailsPage } from '../components/MovieDetailsPage.js'
import { MoviesPage } from '../components/MoviesPage.js'
import { RegisterPage } from '../components/RegisterPage.js'
import { LoginPage } from '../components/LoginPage.js'

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
  {
    path: '/account/register',
    component: RegisterPage,
  },
  {
    path: '/account/login',
    component: LoginPage,
  },
]

export default ROUTES
