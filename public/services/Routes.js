import { HomePage } from '../components/HomePage.js'
import { MovieDetailsPage } from '../components/MovieDetailsPage.js'
import { MoviesPage } from '../components/MoviesPage.js'
import { RegisterPage } from '../components/RegisterPage.js'
import { LoginPage } from '../components/LoginPage.js'
import { AccountPage } from '../components/AccountPage.js'
import FavoritePage from '../components/FavoritesPage.js'
import WatchlistPage from '../components/WatchlistPage.js'

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
  {
    path: '/account',
    component: AccountPage,
    protected: true,
  },
  {
    path: '/account/favorites',
    component: FavoritePage,
    protected: true,
  },
  {
    path: '/account/watchlist',
    component: WatchlistPage,
    protected: true,
  },
]

export default ROUTES
