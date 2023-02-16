FROM node:18-14-alpine

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY . /app
RUN npm ci

CMD ["npm", "run", "start"]
