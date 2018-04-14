import CertificateRoutes from './routes/certificate';
import MessagesRoutes from './routes/messages';
import WebSocket from 'ws'
import config from './config';
import kafka from './kafka';
import restify from 'restify';

export async function createServer() {
  const server = restify.createServer({});
  const wss = new WebSocket.Server({ server: server.server });
  const sockets = [];

  const messageHandler = function messageHandler(messages, topic, partition) {
    messages.forEach((m) => {
      sockets.forEach(ws => {
        ws.send(m.message.value.toString('utf8'));
      });
    });
  };

  kafka.createConsumer(messageHandler);

  server.use((req, res, next) => {
    res.header('Access-Control-Allow-Origin', '*');
    res.header('Access-Control-Allow-Headers', 'X-Requested-With');
    return next();
  });

  MessagesRoutes.register(server);
  CertificateRoutes.register(server);

  server.on('uncaughtException', (req, res, route, err) => {
    console.err(err);
  });

  wss.on('connection', (ws, req) => {
    sockets.push(ws);

    ws.on('close', () => {
      const index = sockets.indexOf(ws);
      if (index !== -1) {
        sockets.splice(index, 1);
      }
    });
  });

  server.listen(config.PORT, () => {
    console.log(`${server.name} listening at ${server.url}`);
  });

  return server;
}
