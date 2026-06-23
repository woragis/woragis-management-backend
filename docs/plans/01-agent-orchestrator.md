# Plano — Agent Orchestrator

Assistente pessoal operacional: Telegram, WhatsApp, voz. Orquestra OpenAI + tools no backend Woragis.

## Por que um serviço separado

| Serviço | Papel |
|---------|-------|
| **Go API** (`server/`) | Domínios, Postgres, tools, auth |
| **agent-worker** (novo) | Canais, Realtime/Whisper, loop de tools, sessão de conversa |
| **whatsapp-worker** (existente) | Só entrega Baileys + cron LeetCode |

Não colocar OpenAI nem lógica de agente dentro do `whatsapp-worker`.

## Arquitetura

```text
Telegram / WhatsApp áudio / Twilio voz
        ↓
   agent-worker
        ↓
   OpenAI (GPT / Realtime / Whisper)
        ↓
   POST/GET /v1/internal/agent/tools/*
        ↓
   Go API → devproject | finance | contacts | content | profile
```

## Estrutura no repositório

```text
management/backend/
├── server/                 # tools + domínios
├── whatsapp-worker/        # entrega WhatsApp
└── agent-worker/           # NOVO
    ├── src/
    │   ├── telegram/
    │   ├── openai/
    │   ├── tools/          # HTTP client → API
    │   ├── session/
    │   └── channels/
    └── Dockerfile
```

## Autenticação

| Chave | Uso |
|-------|-----|
| `AGENT_API_KEY` | agent-worker → `/v1/internal/agent/*` |
| `WORKER_API_KEY` | whatsapp-worker → rotas LeetCode internas (inalterado) |
| `ADMIN_API_KEY` | frontend management (inalterado) |

## Fases

### Fase A0 — Fundação na API

**Objetivo:** superfície interna estável para o agente.

- [ ] Variável `AGENT_API_KEY` em `.env.example`
- [ ] Middleware `AgentAuth` (espelho de `WorkerAuth`)
- [ ] Namespace `/v1/internal/agent/tools/*`
- [ ] Documentar contrato OpenAPI ou markdown por tool
- [ ] Logs de `agent_tool_runs` (auditoria opcional na A0)

**Entregável:** API responde a tools com chave de agente; sem canal ainda.

---

### Fase A1 — Telegram texto

**Objetivo:** validar loop completo agente + tools.

- [ ] Criar `agent-worker/` (Node + TypeScript)
- [ ] Bot Telegram (`TELEGRAM_BOT_TOKEN`)
- [ ] GPT com function calling
- [ ] Tools iniciais: `search_contacts`, `get_contact`, `list_projects`, `create_project`, `get_agent_personality`, `update_agent_personality`
- [ ] Tools finanças leitura: `finance_dashboard`, `finance_summary`, `list_transactions`
- [ ] Carregar personalidade do Redis no boot de cada sessão (ver [04-agent-personality.md](./04-agent-personality.md))
- [ ] `docker-compose.yml`: serviço `agent-worker` + Redis

**Entregável:** conversa por texto no Telegram executando ações reais.

---

### Fase A2 — Áudio no Telegram

**Objetivo:** entrada/saída por voz sem ligação telefônica.

- [ ] Áudio recebido → Whisper → mesmo loop de tools
- [ ] Resposta texto ou TTS (OpenAI `tts-1` / Realtime)
- [ ] Latência aceitável com mensagem “processando…” em áudios longos

**Alternativa:** OpenAI Realtime API em WebSocket para áudio bidirecional.

| Pipeline | Latência típica |
|----------|-----------------|
| Whisper + GPT + TTS | 2–4 s por turno |
| Realtime API | 0,3–1 s por turno |

---

### Fase A3 — WhatsApp conversacional

**Objetivo:** agente responde no WhatsApp (não só cron LeetCode).

- [ ] Canal WhatsApp no `agent-worker` (webhook ou bridge)
- [ ] Envio via `whatsappworkerclient` → `/v1/send`
- [ ] Cron LeetCode permanece no `whatsapp-worker` atual

---

### Fase A4 — Ligação telefônica

**Objetivo:** ligar para o agente como no exemplo “marque reunião com Claudio”.

- [ ] Twilio Voice → WebSocket → OpenAI Realtime
- [ ] Mesmas tools da API
- [ ] Personalidade + voz aplicadas na sessão Realtime (ver doc 04)

---

## Fluxo de tool call (referência)

```text
Usuário: "Marque follow-up com Claudio"
  → search_contacts("Claudio")
  → API retorna 2 contatos
  → Agente: "Qual Claudio? Professor UNIPE ou corretor?"
Usuário: "Professor"
  → update_contact(id, nextFollowUpAt=...)
  → log_interaction(...)
  → Agente confirma
```

## O que o agente já lê hoje

Sem tools dedicadas, tudo passa por espelhar `/v1/admin/*`. Ver [05-entity-catalog.md](./05-entity-catalog.md).

## Riscos

| Risco | Mitigação |
|-------|-----------|
| `VITE_ADMIN_API_KEY` / chave no cliente | Agente usa `AGENT_API_KEY` só server-side |
| Secrets em tools | Tool separada; confirmação obrigatória |
| 502 / API down | Health check antes de iniciar sessão |
