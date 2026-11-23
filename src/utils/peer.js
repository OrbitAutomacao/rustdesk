import { ElMessage } from 'element-plus'
import { useAppStore } from '@/store/app'

export const connectByClient = (id) => {
  const { id_server: server, key } = useAppStore().setting.rustdeskConfig
  if (!id || !server || !key) {
    ElMessage.error('Configuração do servidor Orbit indisponível. Atualize o portal e tente novamente.')
    return
  }
  const isWindows = /Windows/i.test(navigator.userAgent)
  const protocol = isWindows ? 'orbitsuporteremoto' : 'rustdesk'
  const a = document.createElement('a')
  a.href = `${protocol}://${encodeURIComponent(id)}/r@${server}?key=${encodeURIComponent(key)}`
  a.target = '_self'
  a.click()
}
