# rabbitmq-resilient-consumer

The consumer was done using lib "github.com/streadway/amqp" in order to handler channel notified when
occurs an error of connection and then keep application up waiting for retrieve connection to do expected work

To test above behaviour you can run bellow commands and after you can kill rabbitmq container and see the application will not killed

- Run worker consumer
    `docker-compose up --build app`

- Run in other tab rabbitmq
    `docker-compose up rabbitmq`

