FROM node:alpine3.15

WORKDIR /app

COPY ./package.json ./package-lock.json ./index.js ./

RUN npm install

EXPOSE 3000

CMD ["./index.js"]

# #  for testing
# ENTRYPOINT ["tail", "-f" ,"/dev/null"]