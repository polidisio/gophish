# Plan de Pruebas - Modern Dashboard Analytics

## 🎯 Objetivo
Probar los nuevos endpoints de analytics para el dashboard moderno de Gophish.

---

## 📋 Requisitos Previos

### 1. Entorno de Desarrollo
```bash
# Instalar dependencias
cd gophish
go mod download

# Compilar el proyecto
go build

# Verificar que compila sin errores
echo "✓ Compilación exitosa"
```

### 2. Base de Datos
```bash
# Ejecutar las migraciones
# Para SQLite:
sqlite3 gophish.db < db/db_sqlite3/migrations/20260302000000_modern_analytics.sql

# Para MySQL:
mysql -u gophish -p gophish < db/db_mysql/migrations/20260302000000_modern_analytics.sql
```

---

## 🧪 Plan de Testing

### Fase 1: Tests Unitarios

#### 1.1 Test de Modelos
```bash
go test -v ./models -run TestAnalytics
```

**Casos de prueba:**
- [ ] `TestCalculateUserScores` - Verificar cálculo de scores
- [ ] `TestCalculateRiskLevel` - Verificar niveles de riesgo
- [ ] `TestCalculateCampaignAnalytics` - Verificar analytics de campaña

#### 1.2 Test de Cálculos
```bash
go test -v ./models -run TestCalculate
```

**Verificar:**
- Score 0 cuando usuario no interactuó
- Score aumenta al hacer click
- Score aumenta al abrir email
- Score baja al reportar

---

### Fase 2: Tests de API

#### 2.1 Iniciar el servidor
```bash
# Configurar archivo de configuración
cat > config.json << 'EOF'
{
    "admin_server": {
        "listen_url": "127.0.0.1:3333"
    },
    "phish_server": {
        "listen_url": "127.0.0.1:8080"
    },
    "db_name": "sqlite3",
    "db_path": "./gophish.db"
}
EOF

# Ejecutar gophish
./gophish
```

#### 2.2 Obtener API Key
```bash
# Login en https://localhost:3333
# Usuario: admin
# Password: (ver en logs)

# Obtener API key del usuario
# Ir a Account Settings > API Key
```

#### 2.3 Test de Endpoints

```bash
API_KEY="tu_api_key_aqui"
BASE_URL="http://localhost:3333"

# Test 1: Dashboard Summary
curl -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     "$BASE_URL/api/analytics/dashboard"

# Test 2: Campaign Analytics
curl -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     "$BASE_URL/api/analytics/campaign/1"

# Test 3: User Scores
curl -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     "$BASE_URL/api/analytics/users/scores"

# Test 4: Department Metrics
curl -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     "$BASE_URL/api/analytics/departments"

# Test 5: Risk Distribution
curl -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     "$BASE_URL/api/analytics/risk/distribution"

# Test 6: Top At-Risk Users
curl -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     "$BASE_URL/api/analytics/users/at-risk?limit=5"

# Test 7: Time Analytics
curl -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     "$BASE_URL/api/analytics/campaign/1/time"
```

---

### Fase 3: Tests de Integración

#### 3.1 Crear Campaña de Prueba
```bash
# 1. Crear grupo de destinatarios
curl -X POST -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Test Group",
       "targets": [
         {"email": "user1@company.com", "first_name": "User", "last_name": "One"},
         {"email": "user2@company.com", "first_name": "User", "last_name": "Two"},
         {"email": "hr@company.com", "first_name": "HR", "last_name": "Manager"},
         {"email": "it@company.com", "first_name": "IT", "last_name": "Admin"}
       ]
     }' \
     "$BASE_URL/api/groups/"

# 2. Crear template de phishing
curl -X POST -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Test Email",
       "subject": "Verify Your Account",
       "html": "<p>Click here to verify: <a href={{.URL}}>Verify</a></p>"
     }' \
     "$BASE_URL/api/templates/"

# 3. Crear página de landing
curl -X POST -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Login Page",
       "html": "<form><input name=\"username\"/><input name=\"password\"/></form>"
     }' \
     "$BASE_URL/api/pages/"

# 4. Crear perfil de envío (SMTP)
curl -X POST -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Test SMTP",
       "from": "IT Support <it@company.com>",
       "host": "smtp.company.com",
       "username": "test@company.com",
       "password": "password"
     }' \
     "$BASE_URL/api/smtp/"

# 5. Lanzar campaña
curl -X POST -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Test Campaign",
       "template_id": 1,
       "page_id": 1,
       "smtp_id": 1,
       "url": "http://127.0.0.1:8080",
       "group_ids": [1]
     }' \
     "$BASE_URL/api/campaigns/"
```

#### 3.2 Simular Interacciones
```bash
# Esperar a que se envíen los emails
sleep 30

# Simular: Usuario abre email
# (Visitar el link de tracking en el email enviado)

# Simular: Usuario hace click
# (Visitar la URL de la página de phishing)

# Simular: Usuario reporta
# (Hacer click en botón de report en el email)
```

#### 3.3 Verificar Analytics
```bash
# Verificar que aparecen los datos
curl -H "Authorization: Bearer $API_KEY" \
     "$BASE_URL/api/analytics/dashboard" | jq .

# Verificar scores de usuarios
curl -H "Authorization: Bearer $API_KEY" \
     "$BASE_URL/api/analytics/users/scores" | jq .
```

---

## 📊 Resultados Esperados

### Dashboard Summary
```json
{
  "total_campaigns": 1,
  "active_campaigns": 0,
  "completed_campaigns": 1,
  "total_emails_sent": 4,
  "total_clicks": 1,
  "total_opens": 3,
  "total_reports": 1,
  "global_click_rate": 25.0,
  "global_open_rate": 75.0,
  "global_report_rate": 25.0,
  "low_risk_users": 2,
  "medium_risk_users": 1,
  "high_risk_users": 1,
  "critical_risk_users": 0
}
```

### User Scores
```json
[
  {
    "email": "user1@company.com",
    "score": 25,
    "times_clicked": 1,
    "times_opened": 0,
    "times_reported": 0,
    "risk_level": "medium"
  },
  {
    "email": "user2@company.com",
    "score": 10,
    "times_clicked": 0,
    "times_opened": 1,
    "times_reported": 0,
    "risk_level": "low"
  }
]
```

---

## ✅ Checklist de Verificación

- [ ] Compilación sin errores
- [ ] Migraciones ejecutadas correctamente
- [ ] Servidor inicia sin errores
- [ ] Login funciona
- [ ] API key funciona
- [ ] Endpoint /dashboard responde
- [ ] Endpoint /campaign/{id} responde
- [ ] Endpoint /users/scores responde
- [ ] Endpoint /departments responde
- [ ] Datos se calculan correctamente

---

## 🐛 Troubleshooting

### Error: "Campaign not found"
- Verificar que la campaña existe
- Verificar permisos del usuario

### Error: "Invalid JSON"
- Verificar formato de headers
- Asegurar Content-Type: application/json

### Error: "401 Unauthorized"
- Verificar API key
- Verificar que no expiró

### Datos no aparecen
- Esperar a que termine el envío
- Verificar que los emails se enviaron
- Revisar logs del servidor
