# Common
This folder includes the application service hub which is the main facade for all the middleware services.
The folder also includes the factory methods for all the middleware services implementation based on the schema

## Service hub services

#### Database
Facade of configuration database implementing the `database.IDatabase` interface.
This middleware is used by the application to read/write persistent, transactional configuration data.
The concrete implementation base on Postgresql db using `go-yaaf/yaaf-common-postgresql` package
Alternative implementations (for testing) may include:
* In-memory database using `go-yaaf/yaaf-common/database` package

### Datastore
Facade of big data (usually No SQL document store) implementing the `database.IDatastore` interface.
This middleware is used by the application to read/write data to big data store or use aggregation functions (the main use case)
The concrete implementation base on Elasticsearch db using `go-yaaf/yaaf-common-elastic` package
Alternative implementations (for testing) may include:
* In-memory datastore using `go-yaaf/yaaf-common/database` package

### DataCache
Facade of distributed cache implementing the `database.IDataCache` interface.
This middleware is used by the application to read/write key-value pairs (or advanced data structures like lists, maps, queues) from fast distributed cache
The concrete implementation base on Redis using `go-yaaf/yaaf-common-redis` package
Alternative implementations (for testing) may include:
* In-memory cache using `go-yaaf/yaaf-common/database` package

### MessageBus
Facade of real-time messaging infrastructure implementing the `messaging.IMessageBus` interface.
This middleware is used by the application to exchange command and control messages between services
The concrete implementation base on Redis pub/sub using `go-yaaf/yaaf-common-redis` package
Alternative implementations (for testing) may include:
* In-memory message bus using `go-yaaf/yaaf-common/messaging` package

### Streaming
Facade of durable stream processing infrastructure implementing the `messaging.IMessageBus` interface.
This middleware is used by the application to streamline applicative data between services
The concrete implementation base on Google PubSub / PubSub Lite using `go-yaaf/yaaf-common-pubsub` package
Alternative implementations may include:
* Kafka message bus using `go-yaaf/yaaf-common/kafka` package (for use in on-prem environments)
* Redis pub/sub using `go-yaaf/yaaf-common-redis` package (for small scale POCs)
* In-memory message bus using `go-yaaf/yaaf-common/messaging` package (for testing)
