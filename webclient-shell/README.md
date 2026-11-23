# Cliente web Orbit

`index.html` é a moldura em português para o cliente web beta distribuído pela API RustDesk. A interface Flutter continua servida pela API dentro de um `iframe`.

O `Caddyfile` distingue a navegação principal (`Sec-Fetch-Dest: document`) do carregamento do `iframe` (`Sec-Fetch-Dest: iframe`). Assim, `/webclient/` mostra a moldura na aba e entrega o cliente original ao quadro interno. Os demais arquivos e rotas continuam no proxy da API.

Publicação atual: `index.html` em `/opt/orbit-admin/dist/webclient-shell/index.html`, `Caddyfile` em `/opt/orbit-caddy/Caddyfile`.

O cliente beta original tenta abrir WebSocket TLS diretamente na porta 21118 do IP do VPS, que atende WebSocket sem TLS. `prepare_core.py` cria `core.html` a partir do HTML original da API e carrega `ws-shim.js` antes do cliente. O shim redireciona apenas os sockets da Orbit nas portas 21118 e 21119 para `/ws/id` e `/ws/relay` no domínio HTTPS. O Caddy encaminha esses caminhos aos serviços locais. O HTML e o shim usam `Cache-Control: no-store` para facilitar correções futuras.

Validação: o cliente abriu a sessão WebSocket de ID, passou ao relay e exibiu a solicitação de senha para um dispositivo online. A senha e a área remota não foram acessadas no teste.
