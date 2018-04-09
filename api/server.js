import restify from 'restify';
import kafka from 'no-kafka';
import WebSocket from 'ws'
import config from './config';
import { getMessagesHandler, getMessageDetailsHandler } from './controllers/messages';

export async function createServer() {
  const server = restify.createServer({});
  const wss = new WebSocket.Server({ server: server.server });
  const sockets = [];

  const consumer = new kafka.SimpleConsumer({
    connectionString: config.KAFKA_SERVERS,
    groupId: 'kafka-client',
  });
  await consumer.init();

  const messageHandler = function messageHandler(messages, topic, partition) {
    messages.forEach((m) => {
      sockets.forEach(ws => {
        ws.send(m.message.value.toString('utf8'));
      });
    });
  };

  consumer.subscribe('requestsummary', 0, messageHandler);
  consumer.subscribe('responsesummary', 0, messageHandler);

  server.use((req, res, next) => {
    res.header('Access-Control-Allow-Origin', '*');
    res.header('Access-Control-Allow-Headers', 'X-Requested-With');
    return next();
  });

  server.get('/messages', getMessagesHandler);
  server.get('/messages/:id', getMessageDetailsHandler);
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
