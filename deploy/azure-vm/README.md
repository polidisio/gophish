# Gophish Deployment for Azure VM

## Overview

This directory contains the deployment scripts and configuration files for installing Gophish with Microsoft Entra ID SSO on an Azure VM running Ubuntu.

## Domain

- **Domain:** `gophish.savanamed.com`
- **Port:** 443 (HTTPS)

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  gophish.savanamed.com                      │
│                    (Azure VM Ubuntu)                        │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐ │
│  │                   Nginx (Port 443)                    │ │
│  │   • SSL/TLS Termination                              │ │
│  │   • Reverse Proxy                                    │ │
│  └─────────────────────────────────────────────────────┘ │
│         │                         │                        │
│         ▼                         ▼                        │
│  ┌─────────────┐          ┌─────────────┐                 │
│  │   Gophish   │          │   Phishing  │                 │
│  │   Admin UI  │          │   Server    │                 │
│  │  (3333)     │          │   (80)      │                 │
│  └─────────────┘          └─────────────┘                 │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐ │
│  │              PostgreSQL (5432)                       │ │
│  └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Prerequisites

1. **Azure VM** with Ubuntu 20.04+ installed
2. **SSL Certificate** for `gophish.savanamed.com`
3. **Azure AD Application** for SSO (see below)
4. **SSH access** to the VM

## Step 1: Prepare Your Local Machine

### Option A: Build on Linux VM (Recommended)

Simply upload the source code to the VM and build there:

```bash
# Upload source code
scp -r /Users/clot/Desktop/gophish azureuser@gophish.savanamed.com:/tmp/

# On VM:
cd /tmp/gophish
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o gophish .
```

### Option B: Build using Docker

```bash
cd deploy/azure-vm
docker run --rm -v $(pwd)/../../:/workspace -w /workspace golang:1.24 bash -c "GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o gophish ."

# This will create a gophish binary in the project root
```

### Option C: Build on Mac (may not work with CGO)

```bash
cd deploy/azure-vm
chmod +x build.sh
./build.sh
```

This will create a `gophish` binary in the project root.

## Step 2: Prepare Azure VM

### 2.1 Create the VM (if not already created)

```bash
# Example using Azure CLI
az vm create \
  --resource-group gophish-rg \
  --name gophish-vm \
  --image UbuntuLTS \
  --size Standard_B2s \
  --admin-username azureuser \
  --ssh-key-value ~/.ssh/id_rsa.pub
```

### 2.2 Open ports in NSG

```bash
az network nsg rule create \
  --resource-group gophish-rg \
  --nsg-name gophish-vmNSG \
  --name AllowHTTPS \
  --priority 100 \
  --destination-port-range 443 \
  --access Allow \
  --protocol TCP

az network nsg rule create \
  --resource-group gophish-rg \
  --nsg-name gophish-vmNSG \
  --name AllowHTTP \
  --priority 110 \
  --destination-port-range 80 \
  --access Allow \
  --protocol TCP
```

### 2.3 Upload files to VM

```bash
# Upload gophish binary
scp gophish azureuser@gophish.savanamed.com:/tmp/

# Upload install script
scp deploy/azure-vm/install.sh azureuser@gophish.savanamed.com:/tmp/
```

## Step 3: Install Gophish

SSH to your VM and run:

```bash
ssh azureuser@gophish.savanamed.com

# Copy files to home directory
cp /tmp/gophish ~/
cp /tmp/install.sh ~/

# Make executable
chmod +x install.sh

# Run installation (as root)
sudo bash install.sh
```

### SSL Certificates

The installation script will look for SSL certificates at:
- `/etc/ssl/certs/gophish.savanamed.com.crt`
- `/etc/ssl/private/gophish.savanamed.com.key`

If you don't have certificates yet, the script will create a self-signed certificate for testing.

## Step 4: Configure Microsoft Entra ID SSO

### 4.1 Create Azure AD Application

1. Go to [Azure Portal](https://portal.azure.com)
2. Navigate to **Microsoft Entra ID** > **App registrations**
3. Click **New registration**
4. Name: `Gophish SSO`
5. Supported account types: Single tenant
6. Redirect URI: `https://gophish.savanamed.com/login/entra/callback`
7. Click Register

### 4.2 Get Application Credentials

1. Go to **Certificates & secrets**
2. Create a new client secret
3. Note the **Client ID** (Application ID)
4. Note the **Client Secret** value
5. Go to **Overview** > note the **Directory (tenant) ID**

### 4.3 Configure SSO in Gophish

1. Open browser to `https://gophish.savanamed.com`
2. Login with admin credentials (check logs: `sudo journalctl -u gophish -n 50`)
3. Go to **Settings** > **Microsoft Entra ID SSO**
4. Fill in:
   - Enable Microsoft Entra ID SSO: ✅
   - Client ID: (from Azure)
   - Client Secret: (from Azure)
   - Tenant ID: (from Azure)
   - Redirect URI: `https://gophish.savanamed.com/login/entra/callback`
   - Auto-create users: ✅ (optional)
   - Sync departments: ✅ (optional)
5. Click **Save**

## Step 5: Test SSO

1. Open browser to `https://gophish.savanamed.com/login`
2. Click **Sign in with Microsoft**
3. Complete Microsoft authentication
4. You should be redirected to the Gophish dashboard

## Troubleshooting

### Check Gophish Status
```bash
sudo systemctl status gophish
sudo journalctl -u gophish -f
```

### Check Nginx Status
```bash
sudo systemctl status nginx
sudo nginx -t
```

### Check PostgreSQL
```bash
sudo systemctl status postgresql
sudo -u postgres psql -l
```

### Common Issues

1. **404 on login page**: Check Nginx proxy_pass settings
2. **SSO not working**: Verify redirect URI matches Azure app registration
3. **Database connection error**: Check PostgreSQL credentials in config.json

## Commands Reference

```bash
# Start Gophish
sudo systemctl start gophish

# Stop Gophish
sudo systemctl stop gophish

# Restart Gophish
sudo systemctl restart gophish

# View logs
sudo journalctl -u gophish -n 100

# Update Gophish (after new build)
sudo systemctl stop gophish
sudo cp /path/to/new/gophish /opt/gophish/gophish
sudo systemctl start gophish
```

## Files Included

| File | Description |
|------|-------------|
| `install.sh` | Main installation script |
| `build.sh` | Build script for Linux AMD64 |
| `gophish.service` | Systemd service file |
| `nginx-ssl.conf` | Nginx configuration with SSL |
| `.env.example` | Environment variables template |
| `README.md` | This file |

## Security Notes

- Change default database password in `config.json`
- Use strong SSL certificates (not self-signed for production)
- Enable firewall: `sudo ufw enable`
- Configure fail2ban for brute force protection
- Regular backups of PostgreSQL database
