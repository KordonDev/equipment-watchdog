
# Build the app
FROM node:19-alpine AS build

WORKDIR /var/frontend

COPY package.json /var/frontend

COPY package-lock.json /var/frontend

RUN npm ci

ADD . .

RUN npm run build


# Run the app
FROM 2.10-alpine

# if not working https://medium.com/@jonsbun/why-need-to-be-careful-when-mounting-single-files-into-a-docker-container-4f929340834
ADD docker/Caddyfile /etc/caddy/Caddyfile

WORKDIR /var/www

COPY --from=build /var/frontend/build /var/www
ADD docker/entrypoint.sh /entrypoint.sh
RUN chmod -R 777 /entrypoint.sh

CMD ["/entrypoint.sh"]
