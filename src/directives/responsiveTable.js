// Give Element Plus rows readable field labels in the mobile card layout.
const states = new WeakMap()

function annotateTable (root) {
  root.querySelector('.el-table__body-wrapper table.el-table__body')?.classList.add('mobile-card-body')
  const headers = [...root.querySelectorAll('.el-table__header-wrapper thead th')]
    .map(th => th.querySelector('.cell')?.textContent?.trim() || (th.classList.contains('el-table-column--selection') ? 'Selecionar' : ''))
  const rows = root.querySelectorAll('.el-table__body-wrapper > .el-scrollbar > .el-scrollbar__wrap > .el-scrollbar__view > table > tbody > tr, .el-table__body-wrapper > table > tbody > tr')
  rows.forEach(row => {
    ;[...row.children].filter(cell => cell.tagName === 'TD').forEach((cell, index) => {
      const label = headers[index] || ''
      cell.setAttribute('data-mobile-label', label)
      cell.classList.toggle('mobile-actions-cell', /^(ações|actions)$/i.test(label))
    })
  })
}

export const responsiveTable = {
  mounted (el) {
    el.classList.add('responsive-data-table')
    const state = { frame: 0, observer: null }
    const schedule = () => {
      if (state.frame) return
      state.frame = requestAnimationFrame(() => {
        state.frame = 0
        annotateTable(el)
      })
    }
    state.observer = new MutationObserver(schedule)
    state.observer.observe(el, { childList: true, subtree: true })
    states.set(el, state)
    schedule()
  },
  updated (el) {
    annotateTable(el)
  },
  unmounted (el) {
    const state = states.get(el)
    if (!state) return
    state.observer.disconnect()
    cancelAnimationFrame(state.frame)
    states.delete(el)
  },
}
