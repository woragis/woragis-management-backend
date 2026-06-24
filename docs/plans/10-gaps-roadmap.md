# Roadmap — gaps consolidados (fases 10.x)

Plano de implementação dos gaps identificados após Messaging T1–T3 e Presence P4a.

## Fases

| Fase | Escopo | Entregáveis |
|------|--------|-------------|
| **10.0** | Plano + docs README | Este doc, `README.md` atualizado |
| **10.1** | Messaging foundation | Índice único `(channel, external_id)`; resolve inbound → destination; sync Telegram; deprecar UI legada LeetCode templates |
| **10.2** | Projetos dashboard | Métricas por maturidade/intent; filtros avançados no dashboard |
| **10.3** | Presence P4b-lite | Marcar post como publicado + URL real na UI |
| **10.4** | Contacts UI | Lista, busca, follow-ups, detalhe com timeline |
| **10.5** | Finance + CRM | Finanças por contato na página do contato; alertas de vencimento no dashboard finance |
| **10.6** | Agent | Personality UI; inbound com contexto de destination; áudio Telegram (voice → Whisper) |

## Princípios

- Um commit por sub-entrega (backend / frontend / doc separados quando fizer sentido).
- P4b completo (OAuth LinkedIn/Reddit/X) fica **fora** deste roadmap; só “marcar publicado + link”.
- A4 Twilio voz permanece stub; não bloqueia as demais fases.

## Dependências

```text
10.1 resolve-destination → 10.6 agent context
10.4 Contacts UI → 10.5 finance por contato (UI no detalhe do contato)
```

## Status

| Fase | Status |
|------|--------|
| 10.0 | Concluído |
| 10.1 | Concluído |
| 10.2 | Concluído |
| 10.3 | Concluído |
| 10.4 | Concluído |
| 10.5 | Concluído |
| 10.6 | Concluído (A4 Twilio permanece stub) |
