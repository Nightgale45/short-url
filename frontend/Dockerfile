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
FROM nginx:alpine as production

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80


