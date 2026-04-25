function errorHandler(err, req, res, next) {
  console.error(`[ERROR] ${err.message}`);
  res.status(500).json({
    success: false,
    error: err.message || 'Terjadi kesalahan internal',
  });
}

module.exports = errorHandler;