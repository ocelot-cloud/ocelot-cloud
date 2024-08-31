FROM nginx:alpine3.18
COPY default.conf /etc/nginx/conf.d/default.conf
COPY startup.sh /startup.sh
RUN chmod +x /startup.sh
CMD ["/startup.sh"]
