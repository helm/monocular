FROM bitnami/nginx

# built by running:
# $ docker-compose run --rm ui ng build -t production -o rootfs/dist
COPY dist/ /app
COPY vhost.conf /opt/bitnami/nginx/conf/vhosts/
