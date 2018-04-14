import * as fs from 'fs'
import config from '../config';

async function getCertificateHandler(req, res, next) {
  res.header('Content-disposition', 'inline; filename=gofiddle-ca.pem');
  res.header('Content-type', 'application/x-pem-file');
  fs.createReadStream(config.CERTIFICATE_FILE).pipe(res);
  next();
}

export function register(server) {
  server.get('/certificate', getCertificateHandler);
}

export default {
  register,
};
