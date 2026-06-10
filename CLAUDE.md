# CLAUDE.md - gophish

## Project Overview

**Name:** gophish  
**Type:** Security Tool (Go)  
**Description:** Open-source phishing toolkit for businesses and penetration testers. Enables quick setup of phishing engagements and security awareness training. Forked/customized for internal use.  
**Owner:** @polidisio  

## Tech Stack

- **Language:** Go
- **Min Go Version:** 1.10+
- **Platform:** Cross-platform (Windows, macOS, Linux)
- **Docker:** Official Docker image available

## Quick Start

```bash
# Download release for your platform
# Extract zip and run binary

# Or Docker
docker run -p 3333:3333 gophish/gophish

# Access UI
open https://localhost:3333
# Login with admin + password from log output
```

## File Structure

```
gophish/
├── auth/                 # Authentication modules
├── config/               # Configuration
├── static/               # Web UI assets
├── ansible-playbook/     # Deployment playbooks
├── .github/             # GitHub workflows
├── Dockerfile
├── README.md
└── CLAUDE.md
```

## Architecture

- **Web UI:** Go served frontend
- **API:** REST API for campaign management
- **Storage:** SQLite (default)

## Setup

1. Run gophish binary
2. Open browser to https://localhost:3333
3. Login with credentials from log output
4. Default username: admin
5. Default password: (see log output, varies by version)

## Security Note

- For authorized security testing only
- Requires proper authorization before use
- MIT License

## Resources

- Official: https://getgophish.com
- Docs: http://getgophish.com/documentation
- Token optimization tips: shared/claude-optimization-tips.md (Obsidian Vault)

---

**Owner:** Jose Maudisio (@polidisio)  
**Last updated:** 2026-04-24

---

---

## Workflow

### Para tareas simples
Sé directo: "Añade validación al form" — no necesitas explicar contexto.

### Para tareas complejas (>3 pasos)
1. Agent propone plan primero
2. Usuario confirma
3. Agent ejecuta
4. Agent verifica con tests

### Para cada tarea
1. **Plan** → Si son >3 pasos, escribir en `tasks/todo.md`
2. **Verify** → Confirmar antes de cambios grandes
3. **Execute** → Cambio más pequeño posible
4. **Test** → Ejecutar tests, verificar regression
5. **Document** → Actualizar si es necesario

---

## Code Quality

### SIEMPRE
- Código legible y mantenible
- Seguir convenciones del proyecto
- DRY — no duplicar lógica
- Validar input antes de procesar

### NUNCA
- Hardcodear credenciales o tokens
- "Hacky fixes" sin justificación
- Duplicar código sin razón
- Commits sin mensaje descriptivo

---

## Security

- **NUNCA hardcodear** credenciales — usar environment variables
- **NUNCA exponer** tokens en logs o errores
- **Validar input** antes de procesar
- Si hay secrets, usar `.env` y nunca commitearlo

---

## Self-Improvement

### Si cometes un error
1. Documentar en `lessons.md` — qué salió mal, por qué, cómo evitarlo
2. Actualizar este archivo si la convención no estaba clara
3. No repetir

### Si descubres algo útil
- Documentar en notas del proyecto
- Compartir con Jose si es relevante

---

## Token Optimization

### Hacer
- Agrupar múltiples requests en uno
- Editar en vez de reply (menos historial)
- Nuevo tema = nueva conversación
- Planificar en chat, construir en workspace

### Evitar
- Subir carpetas enteras — solo archivos necesarios
- Múltiples prompts cortos seguidos
- Usar Opus para tareas simples
- Mantener contexto irrelevante

**Budget:** ~88% de tokens en conversaciones largas = solo historial. Mantenerlo limpio.

---

## Resources

**Obsidian Vault:** `~/Library/Mobile Documents/iCloud~md~obsidian/Documents/Saraiba/`

| Recurso | Ubicación en Vault |
|---------|---------------------|
| Best practices | `shared/coding-best-practices.md` |
| Optimization tips | `shared/claude-optimization-tips.md` |
| Skills docs | `shared/openclw-skills.md` |
| Guía coding agents | `shared/guia-coding-agents.md` |

---

## Contact

**Jose Maudisio** — @polidisio
**Issues:** Abrir en GitHub o preguntar en Telegram
