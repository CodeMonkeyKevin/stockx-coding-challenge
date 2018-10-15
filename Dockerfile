FROM golang:1.11.1-stretch AS build

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    postgresql \
    postgresql-contrib \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/stockx
COPY . /go/src/stockx

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    sleep 10 &&\
    createdb -O docker -E UTF8 stockxcc_kp

RUN /etc/init.d/postgresql start &&\
    sleep 10 &&\
    psql -d stockxcc_kp -f database/create.psql

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.6/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/9.6/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

ENV GO111MODULE on
RUN go build -o /bin/stockx

# APP Port
EXPOSE 3000

CMD ["sh", "scripts/startup.sh"]
