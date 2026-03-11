#!/bin/bash
set -e

# ============================================
# Gophish Installation Script for Azure VM
# ============================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DOMAIN="gophish.savanamed.com"
GOPHISH_DIR="/opt/gophish"
DB_NAME="gophish"
DB_TYPE="postgres"
DB_USER="gophish"
DB_PASS="ChangeMe123!"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Gophish Installation Script${NC}"
echo -e "${GREEN}  Domain: ${DOMAIN}${NC}"
echo -e "${GREEN}========================================${NC}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root (sudo)${NC}"
    exit 1
fi

# Install Go if not present
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Installing Go...${NC}"
    cd /tmp
    wget -q https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
fi

# Verify Go is available
export PATH=$PATH:/usr/local/go/bin
go version

# Update system
echo -e "${YELLOW}[1/8] Updating system...${NC}"
apt update && apt upgrade -y

# Install PostgreSQL
echo -e "${YELLOW}[2/8] Installing PostgreSQL...${NC}"
apt install -y postgresql postgresql-contrib nginx certbot python3-certbot-nginx

# Configure PostgreSQL
echo -e "${YELLOW}[3/8] Configuring PostgreSQL...${NC}"
su - postgres -c "psql -c \"CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASS}';\"" || true
su - postgres -c "psql -c \"CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};\"" || true
su - postgres -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};\"" || true
su - postgres -c "psql -d ${DB_NAME} -c 'CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";'" || true

# Create Gophish directory
echo -e "${YELLOW}[4/8] Creating Gophish directory...${NC}"
mkdir -p ${GOPHISH_DIR}
mkdir -p ${GOPHISH_DIR}/static/images
mkdir -p /var/log/gophish

# Copy Gophish files
echo -e "${YELLOW}[5/8] Copying Gophish files...${NC}"

# Find the project root directory (where go.mod is located)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Check multiple possible locations for source code
PROJECT_ROOT=""
if [ -f "${SCRIPT_DIR}/../go.mod" ]; then
    PROJECT_ROOT="${SCRIPT_DIR}/.."
elif [ -f "${SCRIPT_DIR}/../../go.mod" ]; then
    PROJECT_ROOT="${SCRIPT_DIR}/../.."
elif [ -f "./go.mod" ]; then
    PROJECT_ROOT="."
fi

# Always build from source to ensure correct architecture
echo -e "${YELLOW}Building Gophish from source for Linux AMD64...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

if [ -z "$PROJECT_ROOT" ]; then
    echo -e "${RED}Error: Cannot find project root (go.mod not found)${NC}"
    echo -e "${YELLOW}Expected go.mod in one of:${NC}"
    echo -e "  - ${SCRIPT_DIR}/../go.mod"
    echo -e "  - ${SCRIPT_DIR}/../../go.mod"
    echo -e "  - ./go.mod"
    exit 1
fi

# Install CGO dependencies if not present
if ! command -v gcc &> /dev/null; then
    echo -e "${YELLOW}Installing CGO dependencies (gcc, sqlite)...${NC}"
    apt-get update && apt-get install -y gcc libc6-dev libsqlite3-dev
fi

cd ${PROJECT_ROOT}

# Clean and build Gophish for Linux AMD64
rm -f gophish
echo "Building Gophish for Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o gophish .

if [ -f "./gophish" ]; then
    echo -e "${GREEN}Build successful! Verifying architecture...${NC}"
    file ./gophish
    cp ./gophish ${GOPHISH_DIR}/
    cp ./VERSION ${GOPHISH_DIR}/
    chmod +x ${GOPHISH_DIR}/gophish
else
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

# Determine PROJECT_ROOT if not already set
if [ -z "$PROJECT_ROOT" ]; then
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    if [ -f "${SCRIPT_DIR}/../go.mod" ]; then
        PROJECT_ROOT="${SCRIPT_DIR}/.."
    elif [ -f "${SCRIPT_DIR}/../../go.mod" ]; then
        PROJECT_ROOT="${SCRIPT_DIR}/../.."
    elif [ -f "./go.mod" ]; then
        PROJECT_ROOT="."
    fi
fi

# Copy migration files based on database type
if [ "${DB_TYPE}" = "postgres" ]; then
    # For PostgreSQL, copy migrations from deploy/azure-vm/migrations
    if [ -d "${PROJECT_ROOT}/deploy/azure-vm/migrations" ]; then
        mkdir -p ${GOPHISH_DIR}/migrations
        cp ${PROJECT_ROOT}/deploy/azure-vm/migrations/*.sql ${GOPHISH_DIR}/migrations/
    fi
else
    # For SQLite, copy from db/db_sqlite3/migrations
    if [ -d "${PROJECT_ROOT}/db/db_sqlite3/migrations" ]; then
        cp -r ${PROJECT_ROOT}/db/db_sqlite3/migrations ${GOPHISH_DIR}/migrations
    fi
fi

# Copy static files
if [ -d "${PROJECT_ROOT}/static" ]; then
    cp -r ${PROJECT_ROOT}/static ${GOPHISH_DIR}/
fi

# Copy templates
if [ -d "${PROJECT_ROOT}/templates" ]; then
    cp -r ${PROJECT_ROOT}/templates ${GOPHISH_DIR}/
fi

# Create config.json
echo -e "${YELLOW}[6/8] Creating configuration...${NC}"
cat > ${GOPHISH_DIR}/config.json <<EOF
{
    "admin_server": {
        "listen_url": "127.0.0.1:3333",
        "use_tls": true,
        "cert_path": "/opt/gophish/gophish.crt",
        "key_path": "/opt/gophish/gophish.key"
    },
    "phish_server": {
        "listen_url": "0.0.0.0:80",
        "use_tls": false
    },
    "db_name": "${DB_NAME}",
    "db_path": "",
    "db_type": "postgres",
    "db_host": "localhost",
    "db_port": 5432,
    "db_user": "${DB_USER}",
    "db_password": "${DB_PASS}",
    "migrations_prefix": "db/db_",
    "contact_address": "security@savanamed.com",
    "logging": {
        "level": ""
    }
}
EOF

# Setup SSL Certificates
echo -e "${YELLOW}[7/8] Setting up SSL certificates...${NC}"

# Check if SSL certificates exist
if [ -f "/etc/ssl/certs/${DOMAIN}.crt" ] && [ -f "/etc/ssl/private/${DOMAIN}.key" ]; then
    echo "Using existing SSL certificates..."
    cp /etc/ssl/certs/${DOMAIN}.crt ${GOPHISH_DIR}/gophish.crt
    cp /etc/ssl/private/${DOMAIN}.key ${GOPHISH_DIR}/gophish.key
elif [ -f "/etc/ssl/certs/${DOMAIN}.pem" ] && [ -f "/etc/ssl/private/${DOMAIN}.key" ]; then
    echo "Using existing SSL certificates (pem format)..."
    cp /etc/ssl/certs/${DOMAIN}.pem ${GOPHISH_DIR}/gophish.crt
    cp /etc/ssl/private/${DOMAIN}.key ${GOPHISH_DIR}/gophish.key
else
    echo -e "${RED}SSL certificates not found!${NC}"
    echo -e "${YELLOW}Please copy your SSL certificates to:${NC}"
    echo -e "  - /etc/ssl/certs/${DOMAIN}.crt"
    echo -e "  - /etc/ssl/private/${DOMAIN}.key"
    echo -e "${YELLOW}Or use Let's Encrypt: certbot --nginx -d ${DOMAIN}${NC}"
    
    # Create self-signed certificate for testing
    echo "Creating self-signed certificate for testing..."
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout ${GOPHISH_DIR}/gophish.key \
        -out ${GOPHISH_DIR}/gophish.crt \
        -subj "/C=ES/ST=Madrid/L=Madrid/O=Savanamed/CN=${DOMAIN}"
fi

chmod 600 ${GOPHISH_DIR}/gophish.key

# Configure Nginx
echo -e "${YELLOW}[8/8] Configuring Nginx...${NC}"
cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup

cat > /etc/nginx/sites-available/gophish-ssl <<'EOF'
server {
    listen 80;
    server_name gophish.savanamed.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name gophish.savanamed.com;

    ssl_certificate /opt/gophish/gophish.crt;
    ssl_certificate_key /opt/gophish/gophish.key;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Gophish Admin UI
    location / {
        proxy_pass https://127.0.0.1:3333;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_ssl_verify off;
    }

    # Phishing server (campaign landing pages)
    location /track {
        proxy_pass http://127.0.0.1:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location / {
        proxy_pass http://127.0.0.1:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
EOF

ln -sf /etc/nginx/sites-available/gophish-ssl /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t

# Setup Systemd service
echo -e "${YELLOW}Setting up systemd service...${NC}"
cat > /etc/systemd/system/gophish.service <<'EOF'
[Unit]
Description=Gophish Phishing Framework
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/gophish
ExecStart=/opt/gophish/gophish
Restart=always
RestartSec=10

# Environment
Environment=GO_ENV=production

# Logging
StandardOutput=append:/var/log/gophish/gophish.log
StandardError=append:/var/log/gophish/gophish.log

[Install]
WantedBy=multi-user.target
EOF

# Set permissions
chown -R root:root ${GOPHISH_DIR}
chmod +x ${GOPHISH_DIR}/gophish

# Reload systemd and start services
systemctl daemon-reload
systemctl enable gophish
systemctl enable nginx
systemctl enable postgresql

systemctl restart postgresql
systemctl restart gophish
systemctl restart nginx

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Installation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "Gophish Admin UI: ${GREEN}https://${DOMAIN}${NC}"
echo -e "Phishing Server:  ${GREEN}http://${DOMAIN}${NC}"
echo ""
echo -e "First login credentials:"
echo -e "  Username: ${GREEN}admin${NC}"
echo -e "  Password: ${GREEN}(check logs: journalctl -u gophish -n 50)${NC}"
echo ""
echo -e "To check status:"
echo -e "  systemctl status gophish"
echo -e "  systemctl status nginx"
echo -e "  journalctl -u gophish -f"
echo ""
