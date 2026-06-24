#!/bin/bash
# Webhook endpoint for Gitea - triggers deployment
# Add this to Gitea repository webhooks: http://localhost:9000/deploy

PORT=9000
SCRIPT_DIR="/home/nara/apps/web"

echo "Starting webhook server on port $PORT..."

while true; do
    # Wait for HTTP request
    REQUEST=$(nc -l -p $PORT -q 1 2>/dev/null)
    
    if echo "$REQUEST" | grep -q "POST /deploy"; then
        echo "Webhook received - deploying..."
        $SCRIPT_DIR/deploy.sh >> /tmp/hatch-deploy.log 2>&1
        echo "Deploy complete"
    fi
done
