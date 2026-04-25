const Docker = require('dockerode');
const net = require('net');
const { generateUsername, generatePassword, generateDbId } = require('./credentials');

const docker = new Docker({ socketPath: process.env.DOCKER_SOCKET || '/var/run/docker.sock' });

const RESOURCE_LIMITS = {
  small: { cpus: 0.5, memoryMB: 256 },
  medium: { cpus: 1, memoryMB: 512 },
  large: { cpus: 2, memoryMB: 1024 },
};

function isPortTaken(port) {
  return new Promise((resolve) => {
    const tester = net.createServer();
    tester.once('error', () => resolve(false));
    tester.once('listening', () => {
      tester.close(() => resolve(true)); // port bisa dibuka = terpake
    });
    tester.listen(port);
  });
}

async function findAvailablePort() {
  let port = 54321;
  while (port < 55000) {
    const available = await isPortTaken(port);
    if (available) return port;
    port++;
  }
  throw new Error('Tidak ada port tersedia di range 54321-55000');
}

async function createDatabase(dbName, packageType, userId) {
  const limits = RESOURCE_LIMITS[packageType];
  if (!limits) {
    throw new Error(`Package tidak valid: ${packageType}. Pilih: small, medium, large`);
  }

  const port = await findAvailablePort();

  const username = generateUsername();
  const password = generatePassword();
  const dbId = generateDbId();

  const container = await docker.createContainer({
    Image: process.env.PG_IMAGE || 'postgres:15',
    name: `jackaldb-${dbId}`,
    Env: [
      `POSTGRES_DB=${dbName}`,
      `POSTGRES_USER=${username}`,
      `POSTGRES_PASSWORD=${password}`,
    ],
    HostConfig: {
      PortBindings: {
        '5432/tcp': [{ HostPort: port.toString() }],
      },
      CpuQuota: Math.round(limits.cpus * 100000),
      Memory: limits.memoryMB * 1024 * 1024,
    },
  });

  await container.start();

  return {
    db_id: dbId,
    container_id: container.id,
    db_name: dbName,
    host: 'localhost',
    port,
    username,
    password,
    package: packageType,
    connection_string: `postgresql://${username}:${password}@localhost:${port}/${dbName}`,
    created_at: new Date().toISOString(),
  };
}

async function deleteDatabase(containerId) {
  const container = docker.getContainer(containerId);
  await container.stop();
  await container.remove();
}

async function listDatabases() {
  const containers = await docker.listContainers({ all: false });
  return containers.filter(c =>
    c.Names.some(name => name.startsWith('/jackaldb-'))
  );
}

module.exports = { createDatabase, deleteDatabase, listDatabases };