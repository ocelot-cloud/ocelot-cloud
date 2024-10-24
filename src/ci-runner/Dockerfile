FROM alpine:3.18.6
RUN apk add --no-cache docker-cli~25.0 docker-compose~2.17 sqlite~3.41
COPY src/frontend/dist /opt/ocelot/dist
COPY src/backend/backend /opt/ocelot/backend
COPY src/backend/assets /opt/ocelot/assets
COPY src/backend/apps_new/sampleuser_nginxdefault /opt/ocelot/sampleuser_nginxdefault
WORKDIR /opt/ocelot
ENTRYPOINT ["/opt/ocelot/backend"]
