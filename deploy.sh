#!/bin/bash
# Deploy hatch.surf - triggered by Gitea webhook or manually

set -e

echo "=== Deploying hatch.surf ==="
cd /home/nara/apps/web

# Pull latest changes
git pull origin main

# Build and deploy
docker build -t hatch-web:latest .
docker stop hatch-web 2>/dev/null || true
docker rm hatch-web 2>/dev/null || true
docker compose up -d

# Verify
sleep 5
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080 | grep -q "200"; then
    echo "✅ hatch.surf deployed successfully!"
else
    echo "❌ Deployment failed!"
    docker logs hatch-web
    exit 1
fi
