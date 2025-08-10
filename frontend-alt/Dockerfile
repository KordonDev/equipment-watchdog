# Build the app
FROM node:19-alpine AS build

WORKDIR /var/frontend

COPY package.json /var/frontend

COPY package-lock.json /var/frontend

RUN npm ci

ADD . .

RUN npm run build


# Run the app
FROM nginx:1.23-alpine

ADD docker/nginx.conf /etc/nginx/nginx.conf

WORKDIR /var/www

COPY --from=build /var/frontend/dist /var/www
ADD docker/entrypoint.sh /entrypoint.sh
RUN chmod -R 777 /entrypoint.sh 

CMD ["/entrypoint.sh"]
