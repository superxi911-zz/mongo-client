FROM index.caicloud.io/caicloud/golang:1.6

ENV MONGO_DB_IP mongo-server:27017

COPY ./mongo-client /


CMD ["/mongo-client"]
