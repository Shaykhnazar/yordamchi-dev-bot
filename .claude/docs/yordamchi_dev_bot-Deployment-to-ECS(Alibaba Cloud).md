# yordamchi_dev_bot Deployment to ECS(Alibaba Cloud)

Created: August 24, 2025 4:43 PM
Last Edited Time: August 24, 2025 5:37 PM
Created By: mr madaminov
Last Edited By: mr madaminov

# Installation steps

## 📋 Method 1: Lightweight Direct Installation (Recommended)

### Step 1: ECS Instance Setup

```bash
*# Choose minimal instance*
Instance Type: ecs.t6-c1m1.small (1 vCPU, 1GB RAM)
Image: Ubuntu 20.04 LTS (minimal)
Storage: 40GB SSD
```

### Step 2: System Optimization

```bash
*# Update system*
apt update && apt upgrade -y

*# Install essential packages only*
apt install -y curl wget git nano htop

*# Optimize memory usage*
echo "vm.swappiness=10" >> /etc/sysctl.conf
echo "vm.vfs_cache_pressure=50" >> /etc/sysctl.conf
sysctl -p

*# Create swap file (important for 1GB RAM)*
fallocate -l 2G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
echo '/swapfile none swap sw 0 0' >> /etc/fstab
```

### Step 3: Install Go

```bash
*# Install Go*
cd /tmp
wget https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz

*# Set environment*
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

*# Verify installation*
go version
```

### Step 4: Install PostgreSQL (Lightweight)

```bash
apt install -y postgresql postgresql-contrib
```

**PostgreSQL Configuration:**

```bash
sudo nano /etc/postgresql/16/main/postgresql.conf
```

and change these parameter`s values accordingly but without duplications!

`# Memory settings for 1GB RAM
shared_buffers = 128MB
effective_cache_size = 512MB
work_mem = 4MB
maintenance_work_mem = 64MB
max_connections = 20

# Logging (minimal)
log_statement = 'none'
log_duration = off`

```bash
# Restart PostgreSQL
systemctl restart postgresql
systemctl enable postgresql
```

**🔐 Configure Database and User:**

Once PostgreSQL is running:

```bash
*# Connect as postgres user*
sudo -u postgres psql
```

**Inside PostgreSQL shell:**

```sql
*- Check current databases*
\

*- Create database*
CREATE DATABASE yordamchi_bot;

*- Create user*
CREATE USER yordamchi_user WITH PASSWORD 'StrongPassword123!';

*- Grant privileges*
GRANT ALL PRIVILEGES ON DATABASE yordamchi_bot TO yordamchi_user;

-- Connect to your bot database
\c yordamchi_bot

-- Grant all privileges on the public schema
GRANT ALL ON SCHEMA public TO yordamchi_user;

-- Grant create privileges specifically
GRANT CREATE ON SCHEMA public TO yordamchi_user;

-- Grant usage on the database
GRANT ALL PRIVILEGES ON DATABASE yordamchi_bot TO yordamchi_user;

-- For good measure, make yordamchi_user the owner of the database
ALTER DATABASE yordamchi_bot OWNER TO yordamchi_user;

-- Grant default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO yordamchi_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO yordamchi_user;

-- Exit
\q
```

**🎯 Quick Status Check:**

After starting PostgreSQL, you should see:

```sql
$ sudo systemctl status postgresql
● postgresql.service - PostgreSQL RDBMS
     Active: active (exited)

$ ss -tlnp | grep :5432
LISTEN 0  128  127.0.0.1:5432  0.0.0.0:*  users:(("postgres",pid=XXXX))

$ sudo -u postgres psql -c "SELECT version();"
 PostgreSQL 16.x on x86_64-pc-linux-gnu
```

### Step 5: Install Build Dependencies

```bash
*# Install C compiler for SQLite driver*
apt install -y build-essential

*# Install git for cloning repository*
apt install -y git
```

**Note:** Your bot uses SQLite and PostgreSQL drivers that require CGO compilation, so build-essential is needed for the compilation process.

### Step 6: Install Nginx (Minimal)

```bash
*# Install Nginx*
apt install -y nginx

*# Remove default configs*
rm /etc/nginx/sites-enabled/default

*# Create bot config*
nano /etc/nginx/sites-available/helperdevbot
```

**Nginx Configuration:**

> "Replace `domain_name` with your actual domain name that will be used for the Telegram bot webhook."
> 

```bash
server {
    server_name domain_name _;

    # Bot webhook endpoint
    location /webhook {
        proxy_pass http://127.0.0.1:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        return 200 "GoBot Server Working!\nNginx listening on port 80\nTime: $time_local";
        add_header Content-Type text/plain;
    }

    location /health {
        return 200 "Health: OK";
        add_header Content-Type text/plain;
    }

    # THIS part should be appears after generating SSL for your domain automatically by Certbot so you can ignore it in this step
    #
    # listen [::]:443 ssl ipv6only=on; # managed by Certbot
    # listen 443 ssl; # managed by Certbot
    # ssl_certificate /etc/letsencrypt/live/domain_name/fullchain.pem; # managed by Certbot
    # ssl_certificate_key /etc/letsencrypt/live/domain_name/privkey.pem; # managed by Certbot
    # include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    # ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot
}
server {
    if ($host = domain_name) {
        return 301 https://$host$request_uri;
    } # managed by Certbot

    listen 80 default_server;
    listen [::]:80 default_server;
    server_name domain_name _;
    return 404; # managed by Certbot

}
```

```bash
*# Enable site*
ln -s /etc/nginx/sites-available/helperdevbot /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
systemctl enable nginx
```

### Step 7: Deploy Go Application

```bash
*# Create application directory*
mkdir -p /opt/yordamchi-bot
cd /opt/yordamchi-bot
sudo chown -R ecs-user:ecs-user /opt/yordamchi-bot

*# Clone or upload your project*
git clone https://github.com/yourusername/yordamchi-dev-bot.git .

*# Create production config*
nano .env
```

**Environment Configuration:**

```bash
BOT_TOKEN=your_actual_bot_token_from_botfather
BOT_MODE=webhook
DATABASE_URL=postgres://yordamchi_user:StrongPassword123!@localhost/yordamchi_bot?sslmode=disable
DB_TYPE=postgres
APP_PORT=8080
```

**Build application**

```bash
*# Install Go dependencies*
go mod tidy

*# Build application*
CGO_ENABLED=1 go build -ldflags="-s -w" -o yordamchi-bot .

*# Make executable*
chmod +x yordamchi-bot

*# Create systemd service*
sudo nano /etc/systemd/system/yordamchi-bot.service
```

**Note:** `CGO_ENABLED=1` is required for SQLite driver compilation.

**Verify Database Connection:**

```bash
*# Test database connection*
sudo -u postgres psql -d yordamchi_bot -c "SELECT 'Connection successful';"

*# Check if tables will be created automatically*
# Your bot will create tables automatically on first run
```

**Systemd Service:**

```bash
[Unit]
Description=Yordamchi Dev Bot
After=network.target

[Service]
Type=simple
User=ecs-user
Group=ecs-user
WorkingDirectory=/opt/yordamchi-bot
ExecStart=/opt/yordamchi-bot/yordamchi-bot
Restart=always
RestartSec=5

# Environment variables
Environment=BOT_TOKEN=your_actual_bot_token_from_botfather
Environment=BOT_MODE=webhook
Environment=DATABASE_URL=postgres://yordamchi_user:StrongPassword123!@localhost/yordamchi_bot?sslmode=disable
Environment=DB_TYPE=postgres
Environment=APP_PORT=8090

# Security settings
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ReadWritePaths=/opt/yordamchi-bot
ProtectHome=yes

[Install]
WantedBy=multi-user.target
```

```bash
# Start service
systemctl daemon-reload
systemctl enable yordamchi-bot
systemctl start yordamchi-bot

# Check status
systemctl status yordamchi-bot
```

### Step 8: SSL Setup (Lightweight)

```bash
*# Install Certbot*
apt install -y certbot python3-certbot-nginx

*# Get certificate*
certbot --nginx -d your-domain.com

*# Auto-renewal*
echo "0 12 * * * /usr/bin/certbot renew --quiet" | crontab -
```

### After change something locally and push the changes into server run these commands!

```bash
cd /opt/yordamchi-bot
git pull
go mod tidy
CGO_ENABLED=1 go build -ldflags="-s -w" -o yordamchi-bot .
sudo systemctl restart yordamchi-bot
sudo systemctl status yordamchi-bot
```