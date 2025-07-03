export class AnimatedLoading extends HTMLElement {
  constructor() {
    super()
  }

  connectedCallback() {
    const elements = this.dataset.elements ?? 1
    const width = this.dataset.width ?? 100
    const height = this.dataset.height ?? 10

    for (let i = 0; i < elements; i++) {
      const wrapper = document.createElement('div')
      wrapper.classList.add('loading-wave')
      wrapper.style.width = `${width}px`
      wrapper.style.height = `${height}px`
      wrapper.style.margin = `10px`
      wrapper.style.display = 'inline-block'
      this.appendChild(wrapper)
    }
  }
}

customElements.define('animated-loading', AnimatedLoading)
