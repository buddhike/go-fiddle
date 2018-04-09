FROM node:latest

RUN apt-get update && \
  apt-get install -y netcat

WORKDIR /app/api

COPY package.json yarn.lock ./

RUN yarn install

COPY ./ ./

CMD [ "sh", "-c", "./wait-for.sh $KAFKA_SERVERS -t 30 -- sleep 5 && npm start" ]
