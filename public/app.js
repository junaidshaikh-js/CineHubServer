window.app = {
  search: event => {
    event.preventDefault()
    const q = document.querySelector('input[name="search"]').value
    console.log(q)
  },
}
