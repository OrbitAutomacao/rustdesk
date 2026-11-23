(() => {
  const NativeWebSocket = window.WebSocket

  function OrbitWebSocket(url, protocols) {
    let destination = url
    try {
      const target = new URL(url, location.href)
      if (target.hostname === '179.197.235.244' && (target.port === '21118' || target.port === '21119')) {
        const service = target.port === '21118' ? 'id' : 'relay'
        target.protocol = 'wss:'
        target.hostname = 'suporte.orbitautomacao.com.br'
        target.port = ''
        target.pathname = `/ws/${service}`
        destination = target.toString()
        console.info('Orbit WebSocket', service, 'via HTTPS')
      }
    } catch {
      // The native constructor reports invalid WebSocket URLs.
    }
    return protocols === undefined ? new NativeWebSocket(destination) : new NativeWebSocket(destination, protocols)
  }

  OrbitWebSocket.prototype = NativeWebSocket.prototype
  Object.setPrototypeOf(OrbitWebSocket, NativeWebSocket)
  window.WebSocket = OrbitWebSocket
})()
