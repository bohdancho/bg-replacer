FROM node:20-slim AS build
RUN npm i -g pnpm@8.6.3
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
WORKDIR /app
COPY . .

RUN pnpm install --frozen-lockfile
RUN pnpm install -g @angular/cli
RUN ng build --configuration production

FROM nginx
COPY --from=build /app/dist/browser /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
