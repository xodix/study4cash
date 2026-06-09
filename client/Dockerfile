FROM node:26.3-alpine3.23 AS build

WORKDIR /buildfolder
COPY . .
RUN npm install && npm run build


FROM httpd:2.4-alpine3.23 AS final

WORKDIR /usr/local/apache2
COPY --from=build /buildfolder/dist htdocs/
COPY httpd.conf conf/httpd.conf

EXPOSE 80
CMD echo "API_URL=$API_URL" > /usr/local/apache2/htdocs/env && httpd-foreground