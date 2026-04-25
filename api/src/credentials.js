const crypto = require('crypto');

function randomString(length) {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

function generateUsername() {
  return `jkl_user_${randomString(8)}`;
}

function generatePassword() {
  return crypto.randomBytes(32).toString('hex');
}

// id pake jkl_ + random 6 karakter
function generateDbId() {
  return `jkl_${randomString(6)}`;
}

module.exports = { generateUsername, generatePassword, generateDbId };