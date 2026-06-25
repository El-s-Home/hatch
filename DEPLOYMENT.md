# Hatch Landing Page — Deployment Guide

## Overview

The Hatch landing page (`hatch.surf`) is a **Next.js** app deployed via **Docker** on the production server.

**Key principle**: All changes MUST go through the git repository and automated deployment. Never edit files directly on the server or use `/var/www`.

---

## Repository

**Location**: `/home/nara/apps/web/`

**Structure**:
```
/app/page.tsx          # Main page component
/app/layout.tsx        # Root layout
/app/globals.css       # Global styles
/components/           # React components (Hero, Features, Why, CTA, Footer, Header)
/public/               # Static assets
Dockerfile             # Docker build configuration
docker-compose.yml     # Container orchestration
deploy.sh              # Deployment script
```

---

## Deployment Flow

### Automatic (Recommended)

1. Push changes to `main` branch
2. Gitea webhook triggers `gitea-webhook.py`
3. Script runs `deploy.sh` automatically
4. Docker image rebuilds and container restarts

### Manual

```bash
cd /home/nara/apps/web
./deploy.sh
```

**What `deploy.sh` does**:
1. `git pull origin main` — fetches latest changes
2. `docker build -t hatch-web:latest .` — rebuilds image
3. Stops and removes old container
4. `docker compose up -d` — starts new container
5. Verifies deployment with health check

---

## Docker Setup

**Container**: `hatch-web`  
**Port**: `8080` (mapped to host)  
**Internal**: nginx serving static Next.js export on port 80  
**Network**: `hatch-network`

**Check status**:
```bash
docker ps | grep hatch-web
docker logs hatch-web
```

**Restart without rebuild**:
```bash
docker compose -f /home/nara/apps/web/docker-compose.yml restart
```

---

## Making Changes

### 1. Edit the source code

```bash
cd /home/nara/apps/web

# Edit components
vim components/Hero.tsx
vim components/Features.tsx
# etc.

# Edit styles
vim app/globals.css
```

### 2. Test locally (optional)

```bash
npm run dev
# Visit http://localhost:3000
```

### 3. Commit and push

```bash
git add .
git commit -m "Description of changes"
git push origin main
```

### 4. Deploy

Either wait for webhook or run manually:
```bash
./deploy.sh
```

### 5. Verify

```bash
curl -s http://localhost:8080 | head -20
# Or visit https://hatch.surf
```

---

## Common Mistakes to Avoid

❌ **DON'T**: Edit files in `/var/www/html/`  
❌ **DON'T**: Edit files directly inside the Docker container  
❌ **DON'T**: Work on static HTML/CSS files in a workspace  
❌ **DON'T**: Copy files to the server without using git  

✅ **DO**: Edit files in `/home/nara/apps/web/`  
✅ **DO**: Commit changes to git  
✅ **DO**: Use `deploy.sh` for deployment  
✅ **DO**: Work on the Next.js components (`app/`, `components/`)  

---

## Rollback

If a deployment causes issues:

```bash
cd /home/nara/apps/web

# View recent commits
git log --oneline -10

# Revert to previous version
git revert HEAD
git push origin main

# Or reset to specific commit
git reset --hard <commit-hash>
git push origin main --force

# Redeploy
./deploy.sh
```

---

## Cache Headers

**HTML pages**: `no-cache, must-revalidate` — browsers always revalidate, ensuring fresh content after deploys.

**Static assets** (CSS/JS/images): `public, immutable` with 30-day expiry. Safe because Next.js uses content hashes in filenames (e.g., `2oioldxlp5hov.css`), so new deployments get new URLs.

**Nginx config**: `/etc/nginx/sites-available/hatch.surf.conf`

If users report seeing old content after deployment:
1. Verify the container is running: `docker ps | grep hatch-web`
2. Check the HTML has new classes: `curl -s http://localhost:8080 | grep bg-gradient-to-r`
3. Ask user to hard refresh: Ctrl+Shift+R (Chrome/Firefox) or Cmd+Shift+R (Mac)
4. If still cached, check nginx config has `Cache-Control: no-cache, must-revalidate`

---

## Troubleshooting

**Container not starting**:
```bash
docker logs hatch-web
docker compose -f /home/nara/apps/web/docker-compose.yml up -d
```

**Changes not showing**:
1. Verify commit pushed: `git log --oneline -1`
2. Check container was rebuilt: `docker images | grep hatch-web`
3. Force rebuild: `cd /home/nara/apps/web && docker build -t hatch-web:latest . && docker compose up -d`

**Port conflict**:
```bash
lsof -i :8080
# Kill conflicting process or change port in docker-compose.yml
```

---

## Contact

- **DevOps**: Jing Yang (Sam Lee)
- **CTO**: Jordan Patel
- **Repository**: `/home/nara/apps/web/`
- **Live site**: https://hatch.surf
