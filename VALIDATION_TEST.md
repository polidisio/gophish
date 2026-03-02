# 🧪 Validación del Modern Dashboard

## Método de Prueba: Demo Mode

Como no tenemos una instancia de Gophish funcionando con datos reales, usaremos un **servidor mock** para demostrar que todo funciona.

---

## 🎯 Script de Validación Automática

### Opción 1: Usar Mock Server (Recomendado)

```bash
# 1. Instalar json-server
npm install -g json-server

# 2. Crear archivo de datos mock
cat > db.json << 'EOF'
{
  "dashboard": {
    "total_campaigns": 15,
    "active_campaigns": 3,
    "completed_campaigns": 12,
    "total_emails_sent": 2500,
    "total_clicks": 425,
    "total_opens": 1875,
    "total_reports": 125,
    "global_click_rate": 17.0,
    "global_open_rate": 75.0,
    "global_report_rate": 5.0,
    "low_risk_users": 45,
    "medium_risk_users": 23,
    "high_risk_users": 12,
    "critical_risk_users": 5,
    "department_stats": [
      {"department": "Engineering", "click_rate": 12.5, "report_rate": 8.2},
      {"department": "Sales", "click_rate": 25.3, "report_rate": 3.1},
      {"department": "HR", "click_rate": 8.7, "report_rate": 12.4},
      {"department": "Finance", "click_rate": 15.2, "report_rate": 6.8},
      {"department": "IT", "click_rate": 5.4, "report_rate": 15.2}
    ]
  },
  "risk_distribution": {
    "low": 45,
    "medium": 23,
    "high": 12,
    "critical": 5
  },
  "at_risk_users": [
    {"email": "user1@company.com", "score": 85, "times_clicked": 3, "times_opened": 5, "times_reported": 0, "risk_level": "critical"},
    {"email": "user2@company.com", "score": 70, "times_clicked": 2, "times_opened": 4, "times_reported": 0, "risk_level": "high"},
    {"email": "user3@company.com", "score": 55, "times_clicked": 1, "times_opened": 6, "times_reported": 1, "risk_level": "medium"},
    {"email": "user4@company.com", "score": 45, "times_clicked": 1, "times_opened": 3, "times_reported": 0, "risk_level": "medium"},
    {"email": "user5@company.com", "score": 30, "times_clicked": 0, "times_opened": 5, "times_reported": 2, "risk_level": "low"}
  ]
}
EOF

# 3. Iniciar servidor mock
json-server --watch db.json --port 3000

# 4. Probar en navegador
# Abrir: http://localhost:3333/dashboard
```

---

## ✅ Checklist de Validación Visual

### 1. Página carga correctamente
- [ ] Sin errores en consola
- [ ] Loading spinner aparece y desaparece
- [ ] Todos los componentes se renderizan

### 2. Cards de Overview
- [ ] Total Campaigns muestra número
- [ ] Active Campaigns muestra número
- [ ] Emails Sent muestra número
- [ ] Total Clicks muestra número

### 3. Rate Cards
- [ ] Click Rate muestra porcentaje con color
- [ ] Open Rate muestra porcentaje con color
- [ ] Report Rate muestra porcentaje con color
- [ ] Badges muestran GOOD/MEDIUM/HIGH correctamente

### 4. Gráficos
- [ ] Risk Distribution Donut Chart renderiza
- [ ] Department Bar Chart renderiza
- [ ] Leyendas son visibles

### 5. Tabla de At-Risk Users
- [ ] 5 usuarios muestran
- [ ] Scores son correctos
- [ ] Progress bars funcionan
- [ ] Risk levels tienen colores correctos

---

## 🔍 Validación de API (Sin Servidor Real)

### Probar Endpoints con curl

```bash
# Simular respuestas de API
curl -X GET http://localhost:3000/dashboard
curl -X GET http://localhost:3000/risk_distribution
curl -X GET http://localhost:3000/at_risk_users
```

### Verificar Estructura JSON

```bash
# Dashboard debe tener:
jq '.dashboard | keys' db.json

# Salida esperada:
# [
#   "active_campaigns",
#   "completed_campaigns", 
#   "global_click_rate",
#   "global_open_rate",
#   ...
# ]
```

---

## 🐛 Troubleshooting

### Error: "ModernDashboard is not defined"
→ Verificar que el archivo JS está incluido en el HTML

### Error: "Chart is not defined"
→ Verificar que Chart.js está cargado antes del script

### Gráficos no aparecen
→ Verificar que el canvas tiene ID correcto

### Datos no cargan
→ Verificar console del navegador para errores de red

---

## 📱 Capturas de Pantalla Esperadas

### Dashboard Vacío (sin datos)
- Todos los valores en 0
- Gráficos vacíos o con 0
- Tabla muestra "No data available"

### Dashboard con Datos
- Cards con números reales
- Donut chart con 4 secciones coloreadas
- Bar chart con 5 departamentos
- Tabla con 5 usuarios y barras de progreso

---

## 🎨 CSS Validation

```bash
# Verificar que el CSS carga
# En navegador: verificar en Network tab que modern_dashboard.css se carga
# Sin errores 404
```

## 📝 Nota para Prueba Real

Para probar con datos reales:
1. Ejecutar Gophish con las nuevas migraciones
2. Crear una campaña de prueba
3. Generar interacciones (clicks, opens, reports)
4. Visitar /dashboard para ver analytics
