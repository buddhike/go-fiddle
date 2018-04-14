export default {
  PORT: process.env['PORT'] || '3001',
  KAFKA_SERVERS: process.env['KAFKA_SERVERS'] || 'localhost:9092',
  MONGODB_SEVER: process.env['MONGODB_SEVER'] || 'mongodb://localhost:27017',
  MONGODB_DATABASE: process.env['MONGODB_DATABASE'] || 'gofiddle',
  CERTIFICATE_FILE: process.env['CERTIFICATE_FILE'] || '../certificates/proxy-ca.pem',
};
