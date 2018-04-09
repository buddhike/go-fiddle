FROM node:latest

WORKDIR /app/web

COPY package.json yarn.lock ./

RUN yarn install

COPY ./ ./

CMD [ "npm", "start" ]
