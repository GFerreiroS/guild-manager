FROM nginx:1.27-alpine
COPY ./templates /usr/share/nginx/html
COPY ./static /usr/share/nginx/html/static
COPY ./public /usr/share/nginx/html/public
COPY nginx.conf /etc/nginx/conf.d/default.conf