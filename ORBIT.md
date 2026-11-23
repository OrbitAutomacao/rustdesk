# Orbit Suporte Remoto

Este fork mantém o aplicativo Windows da Orbit. A versão distribuída sai da branch `orbit-windows` em `https://github.com/OrbitAutomacao/rustdesk`.

## Publicar uma nova versão Windows

1. Faça as alterações na branch `orbit-windows` e envie os commits ao GitHub.
2. A publicação da branch inicia `.github/workflows/orbit-windows.yml`. Também é possível iniciar o workflow manualmente pelo GitHub Actions.
3. Aguarde os testes e o build terminarem. O job `publish-orbit-release` cria a release `orbit-v<VERSION>-<numero da execução>` com `Orbit-Suporte-Remoto-Windows-x64.exe` e `SHA256SUMS.txt`.
4. Confira o SHA-256 do instalador antes de distribuí-lo. Instale e teste login, conexão remota e duas sessões em janelas separadas.

O código define o servidor de ID, o relay, a chave pública e a API em `src/common.rs` (`load_custom_client`). A troca desses endereços exige novo build. Nunca coloque senhas, tokens, chaves privadas ou dados de clientes no repositório.

## Outros componentes

- API e regras de acesso: `https://github.com/OrbitAutomacao/rustdesk-api`.
- Portal em português e cliente web: `https://github.com/OrbitAutomacao/rustdesk-api-web`.
- Servidor de ID (`hbbs`): a instalação atual usa a correção de troca de chaves do [PR rustdesk-server#706](https://github.com/rustdesk/rustdesk-server/pull/706). O procedimento de reprodução e reversão está no `ORBIT.md` do repositório da API.

O aplicativo, o portal, a API e o servidor de ID têm ciclos de publicação independentes. Enviar esta branch gera um instalador; não atualiza o VPS.
