export default {
  restApi: process.env['REACT_APP_REST_API'] || 'http://localhost:8888/',
  websocket: process.env['REACT_APP_WEBSOCKET'] || 'ws://localhost:8888/ws',
};
