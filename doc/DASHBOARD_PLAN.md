# Modern Dashboard & Analytics - Plan de Desarrollo

## 📊 Estado Actual

El dashboard actual de Gophish tiene:
- CampaignStats básico: `total`, `sent`, `opened`, `clicked`, `submitted_data`, `email_reported`, `error`
- Vista de Campaigns por estado
- Timeline de eventos básico

---

## 🚀 Mejoras Propuestas

### Fase 1: Métricas Avanzadas (Semana 1)

#### 1.1 Score de Susceptibilidad por Usuario
```go
type UserPhishingScore struct {
    UserId       int64
    Email        string
    Score        int // 0-100
    TimesClicked int
    TimesOpened  int
    LastCampaign time.Time
    RiskLevel    string // low, medium, high, critical
}
```

#### 1.2 Métricas por Departamento
- Agrupar usuarios por dominio de email
- Calcular riesgo por grupo
- Identificar patrones de comportamiento

#### 1.3 Tiempo Promedio
- Tiempo hasta hacer click
- Tiempo hasta reportar el email
- Tiempo hasta abrir el email

### Fase 2: Visualizaciones (Semana 2)

#### 2.1 Gráficos de Tendencia
- Emails enviados por día/semana
- Tasa de click-through histórica
- Evolución del score de seguridad

#### 2.2 Dashboard Interactivo
- Filtros por fecha, campaña, grupo
- Comparativas entre campañas
- Export a PDF/CSV

#### 2.3 Mapas de Calor
- Días/horas con más clicks
- Patrones geográficos (si hay datos)

### Fase 3: Predicciones (Semana 3)

#### 3.1 Riesgo Predictivo
- AI/ML para predecir usuarios en riesgo
- Recomendaciones personalizadas

#### 3.2 Alertas Automáticas
- Notificaciones cuando alguien hace click
- Reports semanales automáticos

### Fase 4: Integraciones (Semana 4)

- API REST expandida
- Webhooks para SIEM
- Export a Grafana/Kibana

---

## 📁 Archivos a Modificar

| Archivo | Descripción |
|---------|-------------|
| `models/campaign.go` | Añadir nuevos tipos |
| `models/result.go` | Extender Result con más campos |
| `controllers/campaign.go` | Nuevos endpoints API |
| `static/js/dashboard.js` | Nuevo frontend |

---

## ✅ Checklist

- [x] Crear nuevos modelos de datos
- [x] Implementar endpoints API
- [ ] Diseñar UI moderna
- [ ] Añadir gráficos (Chart.js/D3.js)
- [ ] Tests unitarios
- [ ] Documentación

---

## 📦 Fase 1: Completado ✅

### Archivos creados:
- `models/analytics.go` - Modelos de datos
- `models/analytics_service.go` - Servicio de cálculo
- `db/db_sqlite3/migrations/20260302000000_modern_analytics.sql`
- `db/db_mysql/migrations/20260302000000_modern_analytics.sql`

### Modelos implementados:
- `UserPhishingScore` - Score de susceptibilidad por usuario
- `DepartmentMetrics` - Métricas por departamento
- `TimeAnalytics` - Análisis de tiempos
- `CampaignAnalytics` - Analytics por campaña
- `DashboardSummary` - Resumen del dashboard
- `AnalyticsService` - Métodos de cálculo

---

## 📦 Fase 2: Completado ✅

### Endpoints API implementados:

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/analytics/dashboard` | Resumen del dashboard |
| GET | `/api/analytics/campaign/{id}` | Analytics de campaña |
| GET | `/api/analytics/campaign/{id}/time` | Análisis de tiempos |
| GET | `/api/analytics/users/scores` | Scores de todos los usuarios |
| GET | `/api/analytics/users/{email}/score` | Score de un usuario |
| GET | `/api/analytics/users/at-risk` | Usuarios en riesgo |
| GET | `/api/analytics/departments` | Métricas por departamento |
| GET | `/api/analytics/risk/distribution` | Distribución de riesgo |

### Archivos modificados/creados:
- `controllers/api/analytics.go` - Nuevos endpoints
- `controllers/api/server.go` - Rutas añadidas
- `models/campaign.go` - Añadido GetEvents()
- `TEST_PLAN.md` - Plan de pruebas completo

---

## 🚀 Cómo probar

Ver `TEST_PLAN.md` para instrucciones completas.

```bash
# 1. Compilar
go build

# 2. Ejecutar migraciones
sqlite3 gophish.db < db/db_sqlite3/migrations/20260302000000_modern_analytics.sql

# 3. Ejecutar
./gophish

# 4. Probar endpoints
curl -H "Authorization: Bearer API_KEY" \
     http://localhost:3333/api/analytics/dashboard
```
