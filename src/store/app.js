import { defineStore, acceptHMRUpdate } from 'pinia'
import logo from '@/assets/orbit-logo.webp'
import ptBR from 'element-plus/es/locale/lang/pt-br'
import { admin, app, server } from '@/api/config'

const ptLocale = { value: ptBR, sideBarWidth: '252px' }
export const useAppStore = defineStore({
  id: 'App',
  state: () => ({
    setting: {
      title: 'Orbit Suporte Remoto',
      sideIsCollapse: false,
      mobileMenuOpen: false,
      logo,
      lang: 'pt-BR',
      locale: ptLocale,
      appConfig: {
        web_client: 1,
      },
      rustdeskConfig: {
        'id_server': '',
        'key': '',
        'relay_server': '',
        'api_server': '',
      },
    },
  }),

  actions: {
    sideCollapse () {
      this.setting.sideIsCollapse = !this.setting.sideIsCollapse
    },
    toggleMobileMenu () {
      this.setting.mobileMenuOpen = !this.setting.mobileMenuOpen
    },
    closeMobileMenu () {
      this.setting.mobileMenuOpen = false
    },
    loadConfig () {
      this.getAppConfig()
      this.getAdminConfig()
      this.loadRustdeskConfig()
    },
    getAppConfig () {
      console.log('getAppConfig')
      return app().then(res => {
        this.setting.appConfig = res.data
      })
    },
    getAdminConfig () {
      console.log('getAdminConfig')
      return admin().then(res => {
        this.replaceAdminTitle(res.data.title)
      })
    },
    replaceAdminTitle (newTitle) {
      this.setting.title = 'Orbit Suporte Remoto'
      document.title = 'Orbit Suporte Remoto'
    },
    async loadRustdeskConfig () {
      console.log('loadRustdeskConfig')
      const res = await server().catch(_ => false)
      if (res) {
        this.setting.rustdeskConfig = res.data
        const prefix = 'wc-'
        localStorage.setItem(`${prefix}custom-rendezvous-server`, res.data.id_server)
        localStorage.setItem(`${prefix}key`, res.data.key)
        localStorage.setItem(`${prefix}api-server`, res.data.api_server)
      }

    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAppStore, import.meta.hot))
}
