# Gophish - Enhanced Edition

Fork de [gophish/gophish](https://github.com/gophish/gophish) con mejoras para phishing awareness y analytics.

## 🚀 Roadmap

### Fase 1: Dashboard & Analytics

- [ ] Phishing Susceptibility Score por usuario (0-100)
- [ ] Métricas por departamento
- [ ] Gráficos de tendencia temporal
- [ ] Dashboard UI moderna con Chart.js

### Fase 2: Integración SIEM

- [ ] Webhook estructurado (JSON/CEF)
- [ ] Export a Splunk/Elastic
- [ ] API REST pública

### Fase 3: Templates Modernos

- [ ] SMS/Vishing templates
- [ ] QR Phishing
- [ ] AI-generated phishing templates
- [ ] Import desde URLs reales
- [ ] Azure AD SSO

## 📦 Desarrollo

```bash
# Cl https://github.comonar
git clone/polidisio/gophish.git
cd gophish
git checkout enhancement/modern-dashboard

# Instalar dependencias
npm install
go mod download

# Build
go build -o gophish .

# Ejecutar
./gophish
```

## 🔧 Configuración

Editar `config.json` para configurar el servidor admin y la base de datos.

## 📄 Licencia

MIT - Ver LICENSE
