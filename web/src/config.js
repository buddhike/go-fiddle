export default {
  restApi: process.env['REACT_APP_REST_API'] || 'http://localhost:3001/',
  websocket: process.env['REACT_APP_WEBSOCKET'] || 'ws://localhost:3001/ws',
};
