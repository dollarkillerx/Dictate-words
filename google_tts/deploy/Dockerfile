FROM node:18-alpine
WORKDIR /app
COPY . /app
RUN npm config set registry https://registry.npm.taobao.org
RUN npm install
CMD ["node","src/main.js"]
EXPOSE 3030