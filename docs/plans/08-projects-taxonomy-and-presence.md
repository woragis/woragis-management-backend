# Projetos — taxonomia e Presence

Planejamento do domínio de estratégia de projetos e presença em redes (LinkedIn, Reddit, Twitter/X).

## Motivação

Projetos no Woragis não são só repositórios: têm intenção, maturidade, monetização e canais de distribuição. O domínio **Presence** liga posts sociais a projetos e campanhas para job hunting, lançamentos e visibilidade.

## Dimensões do projeto (implementado)

| Campo | Valores |
|-------|---------|
| `intent` | commercial, academic, personal_tool, portfolio, hobby, nonprofit |
| `distribution` | web, play_store, app_store, desktop, internal_only |
| `monetization` | subscription, one_time, ads, services, indirect, none |
| `maturity` | idea, building, mvp, launched, maintenance, sunset |
| `visibilityGoal` | revenue, job_hunting, academic_credit, community, private |

## Entidades Presence (implementado)

| Entidade | Descrição |
|----------|-----------|
| **SocialCampaign** | Campanha com goal, datas, projectId opcional |
| **PostTemplate** | Slug, platform (ou `any`), goal, body com `{{variáveis}}` |
| **SocialPost** | Post com platform, status, hook/body/cta, links a projeto e campanha |

### Variáveis de template

`{{projectName}}`, `{{projectSlug}}`, `{{shortDescription}}`, `{{demoUrl}}`, `{{githubUrl}}`, `{{repoUrl}}`, `{{stack}}`

## Fases

| Fase | Escopo | Status |
|------|--------|--------|
| **P0** | Dimensões no Project + filtros admin | ✅ |
| **P1** | API + UI Presence (campaigns, templates, posts) | ✅ |
| **P2** | Aba Presence no projeto, templates, limites por plataforma, copiar | ✅ |
| **P3** | Agente: list/create posts, apply template via tools | ✅ |
| **P4a** | Lembretes WhatsApp via scheduler + settings | ✅ |
| **P4b** | Publicação automática (APIs oficiais) | Futuro |

## Limites por plataforma (UI)

| Plataforma | Corpo | Título |
|------------|-------|--------|
| LinkedIn | 3 000 | — |
| Twitter/X | 280 | — |
| Reddit | 40 000 | 300 |

## API admin

```
GET/POST   /v1/admin/presence/campaigns
GET/PATCH/DELETE /v1/admin/presence/campaigns/{id}
GET/POST   /v1/admin/presence/templates
GET/PATCH/DELETE /v1/admin/presence/templates/{id}
GET/POST   /v1/admin/presence/posts
GET/PATCH/DELETE /v1/admin/presence/posts/{id}
```

GET/POST   /v1/internal/agent/tools/presence/posts
GET        /v1/internal/agent/tools/presence/templates
POST       /v1/internal/agent/tools/presence/apply-template
```

Filtros posts: `projectId`, `campaignId`, `platform`, `goal`, `status`.

## Tools do agente (P3)

| Tool | Descrição |
|------|-----------|
| `list_social_posts` | Lista posts com filtros |
| `list_post_templates` | Lista templates ativos |
| `apply_post_template` | Renderiza template + vars do projeto (preview) |
| `create_social_post` | Cria draft/scheduled; bloqueia `published` na API |

## Próximo passo sugerido (P4b)

OAuth + adapters LinkedIn / Reddit / X para publicar direto da API (sem copiar manualmente).
