FROM node:20

WORKDIR /usr/app


COPY tsconfig.json .
COPY package.json .
COPY src/ ./src/
COPY config.js ./

RUN npm i --save @types/node
RUN npm i -g npm typescript
RUN tsc

CMD ["node", "./dist/index.js"]