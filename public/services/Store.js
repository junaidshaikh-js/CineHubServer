const Store = {
  token: null,
  get loggedIn() {
    return this.token !== null
  },
}

if (localStorage.getItem('token')) {
  Store.token = localStorage.getItem('token')
}

const proxyStore = new Proxy(Store, {
  set: (target, prop, value) => {
    if (prop === 'token') {
      target[prop] = value
      localStorage.setItem('token', value)
    }
    return true
  },
})

export default proxyStore
