FROM node:20

WORKDIR /usr/app

COPY package.json .
COPY src/ ./src/
COPY .env .
COPY dist/ ./dist/

RUN npm i 
RUN npm i -g --save typescript