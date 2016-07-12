FROM index.caicloud.io/caicloud/ubuntu

ENV MONGO_DB_IP mongo-server:27017

COPY ./mongo-client /


CMD ["/mongo-client"]
