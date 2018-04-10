FROM node:latest AS build

WORKDIR /app/web

COPY package.json yarn.lock ./

RUN yarn install

COPY ./ ./

RUN npm run build


FROM nginx:alpine
COPY --from=build /app/web/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
