import ROUTES from './Routes.js'

export const Router = {
  init() {
    window.addEventListener('popstate', () => {
      Router.go(location.pathname, false)
    })

    document.querySelectorAll('a.navlink').forEach(link => {
      link.addEventListener('click', event => {
        event.preventDefault()
        const href = link.getAttribute('href')
        Router.go(href)
      })
    })

    Router.go(location.pathname + location.search)
  },
  go(route, addToHistory = true) {
    if (addToHistory) {
      history.pushState(null, '', route)
    }

    const routePath = route.includes('?') ? route.split('?')[0] : route

    let pageElement = null

    for (const r of ROUTES) {
      if (typeof r.path === 'string' && r.path === routePath) {
        pageElement = new r.component()
        break
      } else if (r.path instanceof RegExp) {
        const match = r.path.exec(route)
        if (match) {
          pageElement = new r.component()
          const params = match.slice(1)
          pageElement.params = params
          break
        }
      }
    }

    if (pageElement == null) {
      pageElement = document.createElement('h1')
      pageElement.textContent = 'Page Not Found'
    }

    const oldPage = document.querySelector('main').firstElementChild
    if (oldPage) oldPage.style.viewTransitionName = 'old'
    pageElement.style.viewTransitionName = 'new'

    const updatePage = () => {
      document.querySelector('main').innerHTML = ''
      document.querySelector('main').appendChild(pageElement)
    }

    if (!document.startViewTransition) {
      updatePage()
    } else {
      document.startViewTransition(() => {
        updatePage()
      })
    }
  },
}
