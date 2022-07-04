FROM node:16-alpine

COPY . /usr/src/app

WORKDIR /usr/src/app

RUN apk add git;

RUN npm ci

RUN npx hardhat compile