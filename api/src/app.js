require('dotenv').config();
const express = require('express');
const { createDb, deleteDb, getDatabases, getDbStatus } = require('./db-controller');
const errorHandler = require('./error-handler');

const app = express();

app.use(express.json());

app.post('/create-db', createDb);
app.delete('/delete-db/:id', deleteDb);
app.get('/databases', getDatabases);
app.get('/databases/:id/status', getDbStatus);

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'JackalDB API', timestamp: new Date().toISOString() });
});

app.use(errorHandler);

const PORT = process.env.API_PORT || 3000;
app.listen(PORT, () => {
  console.log(`JackalDB API berjalan di http://localhost:${PORT}`);
});