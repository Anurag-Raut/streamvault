FROM nginx:latest

COPY nginx.conf /etc/nginx/nginx.conf
COPY nginx-selfsigned.crt /etc/nginx/ssl/
COPY nginx-selfsigned.key /etc/nginx/ssl/

CMD ["nginx", "-g", "daemon off;"]