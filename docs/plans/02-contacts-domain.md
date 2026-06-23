# Plano — Domínio Contacts

CRM pessoal leve: investidores, clientes em potencial, clientes ativos, sócios e contatos estratégicos.

## Objetivo

Permitir que o agente (e o frontend) **busquem, desambigüem e atualizem pessoas**:

> "Qual Claudio? O da UNIPE ou o corretor?"

Não é agenda acadêmica genérica nem réplica do LinkedIn.

## Tipos de relacionamento

| `relationship` | Uso |
|----------------|-----|
| `lead` | Interesse em assinar / demo |
| `prospect` | Contato frio com potencial |
| `client` | Já paga ou usa produto |
| `investor` | Capital ou mentoria estratégica |
| `partner` | Sócio, co-founder, parceiro |
| `other` | Professor, corretor, contexto livre |

## Estágio (`stage`)

`cold` | `warm` | `active` | `paused` | `churned`

## Modelo `Contact` (v1)

```text
contacts
├── id              UUID PK
├── name            string NOT NULL
├── displayName     string — "Claudio — Professor, UNIPE" (para o agente falar)
├── email           string
├── phone           string
├── telegram        string
├── whatsapp        string
├── organization    string — "UNIPE", "UFPB", "Imobiliária X"
├── roleTitle       string — "Professor", "CEO"
├── relationship    enum (ver acima)
├── stage           enum (ver acima)
├── source          string — "LinkedIn", "evento"
├── notes           text
├── tags            jsonb — ["saas", "notes-app"]
├── projectId       UUID FK → projects (deal principal)
├── lastContactedAt timestamptz
├── nextFollowUpAt  timestamptz
├── active          bool DEFAULT true
├── createdAt / updatedAt
```

**Índices:** `name`, `organization`, `relationship`, `phone`, `email`.

**Decisão v1:** sem tabela `Organization` — string livre em `organization`.

## Modelo `ContactInteraction` (v1.1)

```text
contact_interactions
├── id          UUID PK
├── contactId   UUID FK
├── type        call | meeting | message | email | note
├── channel     telegram | whatsapp | phone | in_person | other
├── summary     text
├── happenedAt  timestamptz
├── createdAt
```

Ao criar interaction → atualizar `Contact.lastContactedAt`.

## API Admin

Pacote: `server/internal/contacts/` (repository → service → `httpserver/contacts.go`)

```text
GET    /v1/admin/contacts?q=&relationship=&organization=&stage=&projectId=
POST   /v1/admin/contacts
GET    /v1/admin/contacts/{id}
PATCH  /v1/admin/contacts/{id}
DELETE /v1/admin/contacts/{id}     → soft: active=false (recomendado)

GET    /v1/admin/contacts/{id}/interactions
POST   /v1/admin/contacts/{id}/interactions
GET    /v1/admin/contacts/{id}/finance   → após F1
```

### Exemplo de busca (desambiguação)

```http
GET /v1/admin/contacts?q=Claudio
```

```json
[
  {
    "id": "...",
    "name": "Claudio",
    "organization": "UNIPE",
    "roleTitle": "Professor",
    "relationship": "other"
  },
  {
    "id": "...",
    "name": "Claudio",
    "organization": "Imobiliária X",
    "roleTitle": "Corretor",
    "relationship": "other"
  }
]
```

## Relação com Projects

- `Contact.projectId` → projeto principal do relacionamento
- Vários contatos podem apontar ao mesmo `Project`
- Agente pode criar projeto e depois associar contato

## Fases

### Fase C1 — Core

- [ ] `models/contacts.go` + GORM AutoMigrate
- [ ] `internal/contacts/repository` + `service`
- [ ] Handlers admin CRUD
- [ ] Busca `q` em name, organization, roleTitle, email, phone
- [ ] Soft delete (`active=false`)
- [ ] Testes unitários repository/service
- [ ] Registrar rotas em `router.go` + wire em `main.go`

**Entregável:** CRUD completo via admin API.

---

### Fase C2 — Interactions

- [ ] Model `ContactInteraction`
- [ ] CRUD aninhado em `/contacts/{id}/interactions`
- [ ] Atualizar `lastContactedAt` automaticamente
- [ ] Tool agente: `log_interaction`, `list_contacts_due_followup`

**Entregável:** histórico de touchpoints por pessoa.

---

### Fase UI (opcional, paralelo)

- [ ] Página Contacts no management frontend
- [ ] Detalhe: abas Interactions + Finance

## Regras de negócio

1. Nome duplicado é permitido — desambiguar por `organization` + `roleTitle`.
2. Agente não apaga contato sem confirmação explícita.
3. `displayName` pode ser gerado: `"{name} — {roleTitle}, {organization}"`.

## Fora do escopo (v1)

- Tabela `Organization`
- `Meeting` + Google Calendar
- Pipeline `Deal` (valor, probabilidade)
- Dados públicos na landing
