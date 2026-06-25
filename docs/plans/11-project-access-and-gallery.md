# Projetos — acesso secreto e galeria multimídia

## Access level

| Valor | Comportamento |
|-------|----------------|
| `private` | Só no admin; não aparece na landing |
| `public` | Visível na landing (`isPublic=true`) |
| `secret` | Nunca na landing; downgrade exige senha |

### Senha de desbloqueio

Configure no servidor:

```bash
# Gerar hash (exemplo com Go):
go run golang.org/x/crypto/bcrypt@latest -cost 12 <<< 'sua-senha-forte'

# Ou em um one-liner:
python -c "import bcrypt; print(bcrypt.hashpw(b'SUA_SENHA', bcrypt.gensalt(12)).decode())"

PROJECT_SECRET_UNLOCK_PASSWORD_HASH='$2a$12$...'
```

Sem essa variável, **não é possível** mudar um projeto de `secret` para `private`/`public`.

## Galeria

- Aceita imagens, GIFs (`image/gif`) e vídeos (`video/mp4`, `video/webm`, `video/quicktime`)
- Limite: 10 MiB imagens/GIFs; 100 MiB vídeos
- Upload em **Media** → adicionar na aba **gallery** do projeto
