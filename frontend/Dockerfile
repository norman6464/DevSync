FROM node:20-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ENV VITE_BACKEND_URL=http://backend:8080

EXPOSE 5173

CMD ["npm", "run", "dev"]
