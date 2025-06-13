FROM node:18-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN chmod +x gen_front_proto.sh && ./gen_front_proto.sh

RUN npm run build:proto

EXPOSE 3000

CMD ["npm", "start"]