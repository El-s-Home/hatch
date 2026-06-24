#!/bin/bash
# Gitea Setup Script for El Foundation
# Self-hosted Git + CI/CD

set -e

GITEA_VERSION="1.22"
GITEA_PORT="3000"
GITEA_DATA="/var/lib/gitea"
GITEA_CONFIG="/etc/gitea"

echo "=== Setting up Gitea for El Foundation ==="

# Install Gitea
echo "Installing Gitea..."
wget -O /usr/local/bin/gitea https://dl.gitea.com/gitea/${GITEA_VERSION}.0/gitea-${GITEA_VERSION}.0-linux-amd64
chmod +x /usr/local/bin/gitea

# Create directories
sudo mkdir -p $GITEA_DATA $GITEA_CONFIG
sudo chown -R nara:nara $GITEA_DATA $GITEA_CONFIG

# Create systemd service
cat > /tmp/gitea.service << 'EOF'
[Unit]
Description=Gitea (Git with a cup of tea)
After=network.target

[Service]
Type=simple
User=nara
Group=nara
WorkingDirectory=/var/lib/gitea
ExecStart=/usr/local/bin/gitea web --config /etc/gitea/app.ini
Restart=always
Environment=USER=nara HOME=/home/nara GITEA_WORK_DIR=/var/lib/gitea

[Install]
WantedBy=multi-user.target
EOF

sudo mv /tmp/gitea.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable gitea
sudo systemctl start gitea

echo "Gitea installed and started on port $GITEA_PORT"
echo "Access at http://localhost:$GITEA_PORT"
echo ""
echo "Next steps:"
echo "1. Open http://localhost:$GITEA_PORT in browser"
echo "2. Complete initial setup"
echo "3. Create admin account"
echo "4. Run: ./configure-gitea.sh"
