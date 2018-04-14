import config from './config';
import kafka from 'no-kafka';

async function retry(func, period, limit) {
  let count = 0;
  while (true) {
    if (count > limit) {
      throw new Error(`Retry limit of ${limit} exceeded`);
    }

    try {
      return await func();
    } catch (err) {
      console.log(`Error, retrying: ${err.message}`);
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }

    count++;
  }
}

export async function createConsumer(handler) {
  const consumer = new kafka.SimpleConsumer({
    connectionString: config.KAFKA_SERVERS,
    groupId: 'kafka-client',
  });

  await retry(() => consumer.init(), 1000, 30);
  await retry(() => consumer.subscribe('request', 0, handler), 1000, 30);
  await retry(() => consumer.subscribe('response', 0, handler), 1000, 30);

  return consumer;
}

export default {
  createConsumer,
}
