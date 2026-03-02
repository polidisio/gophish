# 🧪 Guía Completa: Probar SSO Azure AD en Gophish

---

## 📋 Parte 1: Prueba Local

### Requisitos Previos
- [ ] Cuenta de Azure (puede ser trial gratuito)
- [ ] Go 1.26+ instalado
- [ ] Gophish compilado

---

### Paso 1: Crear App en Azure AD

#### 1.1 Ir al Portal de Azure
```
https://portal.azure.com/#view/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps
```

#### 1.2 Nueva Registración
- Click en **"New registration"**
- Name: `Gophish SSO`
- Supported account types: **"Accounts in this organizational directory only"**
- Redirect URI: **Web** → `http://localhost:3333/oauth/callback`
- Click **Register**

#### 1.3 Copiar IDs
- **Application (client) ID** → Guardar
- **Directory (tenant) ID** → Guardar

#### 1.4 Crear Client Secret
- Ir a **"Certificates & secrets"**
- Click **"New client secret"**
- Description: `Gophish Production`
- Expires: 12 o 24 meses
- **Copiar el valor** (solo se muestra una vez)

#### 1.5 Configurar Redirect URIs
- Ir a **"Authentication"**
- Click **"Add a platform"** → **Web**
- Redirect URIs: `http://localhost:3333/oauth/callback`
- Implicit grant: ✅ ID tokens

#### 1.6 Permisos API
- Ir a **"API permissions"**
- Click **"Add a permission"** → **Microsoft Graph**
- **Delegated permissions**:
  - [x] User.Read
  - [x] Profile
  - [x] Email
- **Application permissions**:
  - [x] Directory.Read.All
- Click **"Grant admin consent"** (si tienes permisos)

---

### Paso 2: Configurar Gophish Local

#### 2.1 Compilar
```bash
cd ~/.openclaw/workspace/gophish
go build -o gophish
```

#### 2.2 Ejecutar Migraciones
```bash
# SQLite
sqlite3 gophish.db < db/db_sqlite3/migrations/20260302000000_modern_analytics.sql
sqlite3 gophish.db < db/db_sqlite3/migrations/20260302000001_azure_sso.sql
```

#### 2.3 Iniciar Gophish
```bash
./gophish
```

#### 2.4 Acceder
- Abrir: http://localhost:3333
- Login: `admin` / contraseña del log

#### 2.5 Configurar SSO
1. Ir a **Settings** → **SSO Configuration**
2. Rellenar:
   - **Display Name**: Azure AD
   - **Client ID**: (copiado de Azure)
   - **Client Secret**: (copiado de Azure)
   - **Tenant ID**: (copiado de Azure)
   - **Redirect URL**: `http://localhost:3333/oauth/callback`
3. ✅ Enable SSO Login
4. Click **Save**

#### 2.6 Probar Login SSO
1. Hacer logout
2. Verificar que aparece botón **"Login with Microsoft"**
3. Click → debería redirigir a Microsoft
4. Login con cuenta de Azure AD
5. Redirigir a dashboard de Gophish

---

### Paso 3: Probar Sync de Usuarios

#### 3.1 En la UI de SSO Settings
1. Click **"Sync Users from Azure AD"**
2. Verificar que aparecen usuarios
3. Revisar tabla `azure_ad_users` en DB

#### 3.2 Verificar con API
```bash
# Obtener API key (del profile)
API_KEY="tu_api_key_aqui"
curl -H "Authorization: Bearer $API_KEY" \
     http://localhost:3333/api/azuread/users
```

---

## ☁️ Parte 2: Instalar en Azure VM

### Requisitos Azure
- [ ] Suscripción de Azure
- [ ] VM con Ubuntu 22.04 LTS ( mínimo 2 vCPU, 4GB RAM)

---

### Paso 1: Crear VM en Azure

#### 1.1 Desde Portal Azure
```
Virtual machines → Create → Ubuntu Server 22.04 LTS
```

#### 1.2 Configuración
- **Resource Group**: `gophish-rg`
- **VM name**: `gophish-vm`
- **Size**: `Standard_B2s` (2 vCPU, 4GB RAM) - ~$10/mes
- **Authentication**: SSH key o password
- **Inbound ports**: 22 (SSH), 80 (HTTP), 443 (HTTPS), 3333 (Gophish), 8080 (Phish Server)

#### 1.3 Conectar por SSH
```bash
ssh usuario@ip-publica-vm
```

---

### Paso 2: Instalar Gophish en VM

#### 2.1 Actualizar sistema
```bash
sudo apt update && sudo apt upgrade -y
```

#### 2.2 Instalar dependencias
```bash
# Go (si no está)
wget https://go.dev/dl/go1.26.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# SQLite
sudo apt install -y sqlite3 build-essential
```

#### 2.3 Clonar y compilar
```bash
git clone https://github.com/polidisio/gophish.git
cd gophish
git checkout feature/azure-sso
go build -o gophish
```

#### 2.4 Configurar dominio (opcional)
```bash
# Instalar Nginx
sudo apt install -y nginx

# Configurar SSL con Let's Encrypt
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d tu-dominio.com
```

---

### Paso 3: Configurar Azure AD para Producción

#### 3.1 Actualizar Redirect URIs en Azure
- Cambiar de `http://localhost:3333` a `https://tu-dominio.com/oauth/callback`

#### 3.2 Actualizar Gophish config
```bash
# Editar config.json
{
    "admin_server": {
        "listen_url": "0.0.0.0:3333"
    },
    "phish_server": {
        "listen_url": "0.0.0.0:8080"
    },
    "db_name": "sqlite3",
    "db_path": "/opt/gophish/gophish.db"
}
```

---

### Paso 4: Crear Servicio Systemd

#### 4.1 Crear servicio
```bash
sudo nano /etc/systemd/system/gophish.service
```

```ini
[Unit]
Description=Gophish Phishing Framework
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/opt/gophish
ExecStart=/opt/gophish/gophish
Restart=always

[Install]
WantedBy=multi-user.target
```

#### 4.2 Habilitar servicio
```bash
sudo systemctl daemon-reload
sudo systemctl enable gophish
sudo systemctl start gophish
sudo systemctl status gophish
```

---

### Paso 5: Configurar Firewall Azure

#### 5.1 En Portal Azure
```
VM → Networking → Inbound security rules
```

#### 5.2 Añadir reglas
| Priority | Name | Port | Action |
|----------|------|------|--------|
| 100 | HTTP | 80 | Allow |
| 110 | HTTPS | 443 | Allow |
| 120 | Gophish Admin | 3333 | Allow |
| 130 | Phish Server | 8080 | Allow |
| 1000 | DenyAll | * | Deny |

---

### Paso 6: Verificar Funcionamiento

#### 6.1 Probar endpoints
```bash
# Verificar que responde
curl http://localhost:3333

# Ver logs
sudo journalctl -u gophish -f
```

#### 6.2 Probar desde navegador
- Admin: `https://tu-dominio.com:3333`
- Phishing: `http://tu-dominio.com:8080`

#### 6.3 Probar SSO
1. Hacer logout
2. Login con Microsoft
3. Verificar que funciona

---

## 🔧 Troubleshooting

### Error: "invalid_client"
→ Verificar Client ID y Secret

### Error: "AADSTS50011"
→ Verificar Redirect URI en Azure

### Error: "No access token"
→ Verificar permisos de API en Azure

### Error: "User not found"
→ El usuario de Azure no tiene licencia o no existe

---

## 📝 Comandos Rápidos

```bash
# En VM
sudo systemctl restart gophish     # Reiniciar
sudo journalctl -u gophish -f     # Ver logs
sudo systemctl status gophish      # Estado

# Backup DB
cp /opt/gophish/gophish.db /opt/gophish/backup-$(date +%Y%m%d).db
```

---

## ✅ Checklist Final

- [ ] App creada en Azure
- [ ] Permisos API otorgados
- [ ] VM creada en Azure
- [ ] Gophish compilado
- [ ] Migraciones ejecutadas
- [ ] SSO configurado
- [ ] Login SSO funciona
- [ ] User sync funciona
- [ ] Firewall configurado
- [ ] SSL configurado (producción)

---

## 📞 Soporte

Si tienes errores, revisa:
1. Logs de Gophish: `journalctl -u gophish`
2. Logs de Azure AD: Azure Portal → App Registration → Monitoring
3. Consola del navegador (F12) para errores JS
