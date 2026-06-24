#!/bin/bash
# Configure Gitea for El Foundation
# Run after initial Gitea setup

set -e

GITEA_URL="http://localhost:3000"
GITEA_TOKEN="${GITEA_TOKEN:-your-admin-token}"

echo "=== Configuring Gitea for El Foundation ==="

# Create organization
echo "Creating El Foundation organization..."
curl -X POST "$GITEA_URL/api/v1/orgs" \
  -H "Authorization: token $GITEA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "el-foundation",
    "full_name": "El Foundation",
    "description": "El Foundation - Open Source Projects",
    "website": "",
    "location": "",
    "visibility": "public"
  }'

# Create hatch-surf repository
echo "Creating hatch-surf repository..."
curl -X POST "$GITEA_URL/api/v1/orgs/el-foundation/repos" \
  -H "Authorization: token $GITEA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "hatch-surf",
    "description": "Hatch.surf web application",
    "private": false,
    "auto_init": false,
    "default_branch": "main"
  }'

# Add remote to local repo
echo "Adding Gitea remote to local repo..."
cd /home/nara/apps/web
git remote add origin http://localhost:3000/el-foundation/hatch-surf.git
git push -u origin main

echo ""
echo "=== Configuration Complete ==="
echo "Repository: $GITEA_URL/el-foundation/hatch-surf"
echo ""
echo "To enable Gitea Actions (CI/CD):"
echo "1. Go to repository Settings > Actions"
echo "2. Enable Actions"
echo "3. Add secrets: SERVER_HOST, SERVER_USER, SERVER_SSH_KEY"
echo "4. Push to main branch to trigger deployment"
