FROM index.caicloud.io/caicloud/golang:1.6

ENV MONGO_DB_IP 180.101.191.213:31599

COPY ./mongo-client /


CMD ["/mongo-client"]
