import { execSync } from 'child_process';

async function globalTeardown() {
  console.log('Tearing down test environment...');
  
  // Stop and remove docker containers, networks, and volumes
  execSync('docker compose -p e2e_tests -f docker-compose.test.yml down -v', { stdio: 'inherit' });
  
  console.log('Test environment torn down.');
}

export default globalTeardown;
