import mongodb from 'mongodb';
import config from '../config';

const absoluteUriExpression = /^https?:\/\//g;

function getHeaderValue(headers, name, defaultValue) {
  const header = headers.filter(h => new RegExp(`^${name}$`, 'i').test(h.name))[0];
  if (header) {
    return header.value;
  }
  return defaultValue
}

function getUri(uri, headers) {
  if (/^https?:\/\//g.test(uri)) {
    return uri;
  }
  const host = getHeaderValue(headers, 'host');
  if (host) {
    return `https://${host}${uri}`;
  }
  return uri;
}

const database = (async () => {
  const client = await mongodb.MongoClient.connect(config.MONGODB_SEVER);
  const db = client.db(config.MONGODB_DATABASE);

  return db;
})();

export async function getMessagesHandler(req, res, next) {
  const db = await database;
  const messages = (
    await db.collection('messages')
      .find({})
      .project({
        '_id': 1,
        'request.method': 1,
        'request.uri': 1,
        'request.headers': 1,
        'response.statuscode': 1,
      })
      .toArray()
    ).map(r => ({
      id: r._id,
      method: r.request.method,
      uri: getUri(r.request.uri, r.request.headers),
      statuscode: (r.response || {}).statuscode,
    }));

  res.json(messages);
  next();
}

export async function getMessageDetailsHandler(req, res, next) {
  const db = await database;
  const message = (await db.collection('messages')
    .find({
      _id: req.params.id,
    })
    .toArray())
    .map(r => ({
      id: r._id,
      request: r.request,
      response: r.response,
    }))[0];

  if (!message) {
    res.status(404);
  } else {
    res.json(message);
  }

  next();
}

export default {
  getMessagesHandler,
  getMessageDetailsHandler,
};
