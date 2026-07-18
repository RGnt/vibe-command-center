import { execSync } from 'child_process';
import http from 'http';

async function waitForService(url: string, timeout = 120000) {
  const start = Date.now();
  let lastErr = null;
  while (Date.now() - start < timeout) {
    try {
      const ok = await new Promise((resolve) => {
        const req = http.get(url, (res) => {
          resolve(res.statusCode === 200 || res.statusCode === 404); // 404 is okay if API endpoint doesn't exist yet, we just want the server to be up
        });
        req.on('error', (err) => {
          lastErr = err;
          resolve(false);
        });
        req.end();
      });
      if (ok) return true;
    } catch (e) {
      lastErr = e;
    }
    await new Promise(resolve => setTimeout(resolve, 2000));
  }
  throw new Error(`Service at ${url} not ready after ${timeout}ms. Last error: ${lastErr?.message}`);
}

async function globalSetup() {
  console.log('Starting test environment...');
  
  // Start docker containers
  execSync('docker compose -p e2e_tests -f docker-compose.test.yml up -d --build', { stdio: 'inherit' });
  
  console.log('Waiting for frontend and backend to become ready...');
  
  // Wait for the app to be accessible
  await waitForService('http://127.0.0.1:8089');
  await waitForService('http://127.0.0.1:8089/api/workflows'); // Wait for API
  
  console.log('Test environment is ready!');
}

export default globalSetup;
