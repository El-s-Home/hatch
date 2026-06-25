#!/usr/bin/env python3
"""
Gitea webhook server for hatch.surf deployment
Listens on port 9000 for POST requests from Gitea
"""

import http.server
import json
import subprocess
import logging
from pathlib import Path

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('/home/nara/logs/gitea-webhook.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

DEPLOY_SCRIPT = Path(__file__).parent / "deploy.sh"

class WebhookHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        """Handle POST request from Gitea webhook"""
        content_length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_length) if content_length > 0 else b''
        
        try:
            payload = json.loads(body) if body else {}
            event = self.headers.get('X-Gitea-Event', 'unknown')
            
            logger.info(f"Received webhook: {event}")
            logger.info(f"Payload: {json.dumps(payload, indent=2)}")
            
            # Only deploy on push to main branch
            if event == 'push':
                ref = payload.get('ref', '')
                if ref == 'refs/heads/main':
                    logger.info("Push to main branch - triggering deployment...")
                    self._deploy()
                else:
                    logger.info(f"Ignoring push to {ref}")
            else:
                logger.info(f"Ignoring event: {event}")
            
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({"status": "ok"}).encode())
            
        except Exception as e:
            logger.error(f"Error processing webhook: {e}")
            self.send_response(500)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({"status": "error", "message": str(e)}).encode())
    
    def _deploy(self):
        """Run the deployment script"""
        try:
            logger.info(f"Running deploy script: {DEPLOY_SCRIPT}")
            result = subprocess.run(
                [str(DEPLOY_SCRIPT)],
                capture_output=True,
                text=True,
                timeout=300  # 5 minute timeout
            )
            
            if result.returncode == 0:
                logger.info("Deployment successful!")
                logger.info(f"Output: {result.stdout}")
            else:
                logger.error(f"Deployment failed with code {result.returncode}")
                logger.error(f"Stdout: {result.stdout}")
                logger.error(f"Stderr: {result.stderr}")
                
        except subprocess.TimeoutExpired:
            logger.error("Deployment timed out after 5 minutes")
        except Exception as e:
            logger.error(f"Error running deployment: {e}")
    
    def do_GET(self):
        """Health check endpoint"""
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps({"status": "healthy"}).encode())
    
    def log_message(self, format, *args):
        """Override to use our logger"""
        logger.info(format % args)

def main():
    server = http.server.HTTPServer(('0.0.0.0', 9000), WebhookHandler)
    logger.info("Starting webhook server on port 9000...")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logger.info("Shutting down webhook server...")
        server.shutdown()

if __name__ == '__main__':
    main()
