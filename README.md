ksqlDB
===========


This is doing basic setup of ksqlDB with the Go client to see if it's usable.
The idea is very cool:

- Use Kafka as the Source Of Truth database
- Use an http client (non-http protocols are a pain when dealing with containers and meshes
- Allow a SQL interface; over kafka
- Allow an evented form of SQL


Once you have docker installed:

```
./runit
```

It's this:

(tutorial)[https://ksqldb.io/quickstart.html#quickstart-content]
