const { createDatabase, deleteDatabase, listDatabases } = require('./docker-service');

const dbRegistry = new Map();

async function sendLog(event, userId, dbId, status, detail, errorMsg = '') {
  try {
    await fetch(process.env.LOGGER_URL || 'http://localhost:8081/log', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        Event: event,
        UserID: userId,
        DbID: dbId,
        Status: status,
        Detail: detail,
        ErrorMsg: errorMsg,
      }),
    });
  } catch (err) {
    console.warn('[WARN] Logger service tidak dapat dihubungi:', err.message);
  }
}

async function createDb(req, res, next) {
  try {
    const { db_name, package: packageType, user_id } = req.body;
    if (!db_name || !packageType || !user_id) {
      return res.status(400).json({
        success: false,
        error: 'Field wajib: db_name, package, user_id',
      });
    }

    const result = await createDatabase(db_name, packageType, user_id);

    dbRegistry.set(result.db_id, {
      ...result,
      user_id,
    });

    await sendLog(
      'DB_CREATED',
      user_id,
      result.db_id,
      'success',
      `Container spawned on port ${result.port}, image: postgres:15, package: ${packageType}`
    );

    res.status(201).json({
      success: true,
      data: result,
    });
  } catch (err) {
    await sendLog('PROVISION_FAILED', req.body.user_id || 'unknown', '-', 'failed', '', err.message);
    next(err);
  }
}

async function deleteDb(req, res, next) {
  try {
    const { id } = req.params;
    const dbInfo = dbRegistry.get(id);

    if (!dbInfo) {
      return res.status(404).json({
        success: false,
        error: `Database dengan id "${id}" tidak ditemukan`,
      });
    }

    await deleteDatabase(dbInfo.container_id);
    dbRegistry.delete(id);

    await sendLog(
      'DB_DELETED',
      dbInfo.user_id,
      id,
      'success',
      `Database ${dbInfo.db_name} berhasil dihapus`
    );

    res.json({
      success: true,
      message: `Database ${id} berhasil dihapus`,
    });
  } catch (err) {
    next(err);
  }
}

async function getDatabases(req, res, next) {
  try {
    const containers = await listDatabases();
    res.json({
      success: true,
      total: containers.length,
      data: containers.map(c => ({
        container_id: c.Id,
        name: c.Names[0].replace('/', ''),
        status: c.Status,
        ports: c.Ports,
      })),
    });
  } catch (err) {
    next(err);
  }
}

async function getDbStatus(req, res, next) {
  try {
    const { id } = req.params;
    const dbInfo = dbRegistry.get(id);

    if (!dbInfo) {
      return res.status(404).json({
        success: false,
        error: `Database "${id}" tidak ditemukan`,
      });
    }

    res.json({
      success: true,
      data: {
        db_id: id,
        db_name: dbInfo.db_name,
        status: 'running',
        port: dbInfo.port,
        package: dbInfo.package,
        created_at: dbInfo.created_at,
      },
    });
  } catch (err) {
    next(err);
  }
}

module.exports = { createDb, deleteDb, getDatabases, getDbStatus };