# Base stage
FROM node:24-alpine AS base

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

#Development stage
FROM base AS development

EXPOSE 80

CMD ["npm", "run", "dev"]

#Prod build
FROM base AS build

RUN npm run build

#Prod
FROM nginx:alpine AS production

RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown nginx:nginx /var/run/nginx.pid

COPY --from=build --chown=nginx:nginx /app/dist /usr/share/nginx/html
COPY --chown=nginx:nginx nginx.conf /etc/nginx/conf.d/default.conf

USER nginx

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
