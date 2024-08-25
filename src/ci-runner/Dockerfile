FROM alpine:3.18.6
RUN apk add --no-cache docker-cli~25.0 docker-compose~2.17 sqlite~3.41
COPY src/cloud/frontend/dist /opt/ocelot/dist
COPY src/cloud/backend/backend /opt/ocelot/backend
COPY src/cloud/backend/stacks /opt/ocelot/stacks
WORKDIR /opt/ocelot
ENTRYPOINT ["/opt/ocelot/backend"]
