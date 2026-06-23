# Catálogo de entidades — O que o agente pode ler hoje

Inventário do `management/backend` (Postgres principal) e APIs disponíveis. Base para definir tools do agente.

## Bancos de dados

| Banco | Serviço | Conteúdo |
|-------|---------|----------|
| Postgres principal | Go API | Todos os domínios abaixo |
| Postgres whatsapp | whatsapp-worker | `whatsapp_groups`, `message_log`, `worker_state` — **sem API admin** |
| Creatives (externo) | HTTP | Jobs de imagem — indireto via thumbnails |

## Entidades no Postgres principal

### Dev / Projetos (`devproject`)

| Entidade | Campos principais |
|----------|-------------------|
| **Project** | name, slug, descriptions, status, stack, URLs, notes, isPublic, featured, coverImageId, parentProjectId |
| **ProjectLink** | type, url, environment, label |
| **ProjectDomain** | domain, registrar, expiresAt |
| **ProjectSecret** | name, valor criptografado, environment |
| **ProjectGallery** | mediaAssetId, caption |
| **ProjectEnv** | key, value, environment |

### Perfil (`profile`)

| Entidade | Campos principais |
|----------|-------------------|
| **Profile** | displayName, headline, bio, avatarId, availability, socialLinks, resumeAssetId |

### Mídia (`media`)

| Entidade | Campos principais |
|----------|-------------------|
| **MediaAsset** | filename, mimeType, publicUrl, altText |

### Finanças (`finance`)

| Entidade | Campos principais |
|----------|-------------------|
| **IncomeSource** | name, type, amountCents, frequency, projectId |
| **Expense** | name, category, amountCents, dueDate, projectId |
| **Transaction** | type, amountCents, description, date, links |
| **Invoice** | fatura de **cartão** (não cliente B2B) |
| **InvoiceItem** | linhas da fatura |
| **BudgetPlan** | year, month, category, plannedCents |

Agregados: `FinanceDashboard`, `MonthlySummary`, `CalendarEvent`.

### Conteúdo LeetCode (`content`)

| Entidade | Campos principais |
|----------|-------------------|
| **LeetcodeVideo** | título, problema, difficulty, URLs, WhatsApp sent flags |
| **ContentThumbnail** | prompt, status, creativesJobId |
| **ContentPromptTemplate** | templates de imagem |
| **LeetcodeChannelSettings** | horários de post |
| **WhatsappMessageTemplate** | corpo das mensagens automáticas |

### Planejado (ainda não existe)

| Entidade | Doc |
|----------|-----|
| **Contact** | [02-contacts-domain.md](./02-contacts-domain.md) |
| **ContactInteraction** | [02-contacts-domain.md](./02-contacts-domain.md) |
| **AgentPersonality** | [04-agent-personality.md](./04-agent-personality.md) |

## Superfícies de API

### Admin (`X-Admin-Key`) — leitura total

O agente deve espelhar estas rotas via tools internas:

| Domínio | GET principais |
|---------|----------------|
| Dashboard | `/v1/admin/dashboard` |
| Projetos | list, get, secrets*, envs*, links, domains, gallery |
| Profile | `/v1/admin/profile` |
| Media | list, get |
| Finance | dashboard, summary, calendar, income, expenses, transactions, invoices, budgets |
| Content | videos, thumbnails, templates, settings, whatsapp-* |

\* Secrets e envs: dados sensíveis — tools com confirmação.

### Público (sem auth)

Subset para landing — agente raramente precisa exceto para ver o que está público.

### Internal worker (`X-Worker-Key`)

Apenas LeetCode WhatsApp dispatch — não é superfície do agente.

## O que o agente **não** lê hoje

| Dado | Motivo |
|------|--------|
| Contatos | não implementado |
| Personalidade | não implementado |
| `message_log` WhatsApp | DB do worker, sem rota |
| Reuniões / Google Calendar | não existe |
| Memória de conversa longa | não existe (futuro: Postgres ou Redis) |

## Mapeamento sugerido: entidade → tools

| Domínio | Tools (futuro) |
|---------|----------------|
| Projects | `list_projects`, `get_project`, `create_project`, `update_project` |
| Profile | `get_profile` |
| Finance | `finance_dashboard`, `finance_summary`, `list_transactions`, … |
| Contacts | `search_contacts`, `get_contact`, `create_contact`, … |
| Personality | `get_agent_personality`, `update_agent_personality` |
| Content | `list_leetcode_videos`, `whatsapp_preview` (baixa prioridade) |

## Dados sensíveis

| Recurso | Política agente |
|---------|-----------------|
| ProjectSecret | tool dedicada; confirmar antes de revelar |
| ProjectEnv | idem |
| Finance | confirmar antes de criar transação |
| Admin key / Agent key | nunca em prompt |

## Atualização

Quando novos domínios forem implementados, atualizar este arquivo e [README.md](./README.md).
