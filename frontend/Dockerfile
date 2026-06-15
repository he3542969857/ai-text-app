# ---- build ----
FROM node:22-alpine AS build
WORKDIR /app
# 使用国内镜像加速依赖安装
RUN npm config set registry https://registry.npmmirror.com
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
RUN npm run build

# ---- runtime ----
FROM nginx:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
