FROM confluentinc/cp-kafka:latest

RUN apt-get update && \
  apt-get install -y net-tools

WORKDIR /home/kafka

COPY scripts scripts
# RUN ls -a /etc/confluent/docker
CMD [ "./scripts/start.sh" ]
