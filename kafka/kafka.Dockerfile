FROM confluentinc/cp-kafka:latest

WORKDIR /home/kafka

COPY scripts scripts
# RUN ls -a /etc/confluent/docker
CMD [ "./scripts/start.sh" ]
