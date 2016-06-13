
### Development Server

```
docker-compose -f ./backend_cassandra.yaml up -d
```

### Database Setup

```
$ bin/cqlsh

Connected to Test Cluster at 127.0.0.1:9042.
[cqlsh 5.0.1 | Cassandra 3.0.6 | CQL spec 3.4.0 | Native protocol v4]
Use HELP for help.

cqlsh> CREATE KEYSPACE semver WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
cqlsh> use semver;
cqlsh:semver> CREATE TABLE db (
          ... id text,
          ... key text,
          ... val text,
          ... PRIMARY KEY (id, key));

```

### Testing

```
cqlsh:semver> INSERT INTO db (id, key, val) VALUES ('04ed14d5-f0dd-4b4c-81b3-635563760da2', 'version', '0.0.1');
cqlsh:semver> INSERT INTO db (id, key, val) VALUES ('04ed14d5-f0dd-4b4c-81b3-635563760da2', 'archive/0.0.2', '0.0.2');
cqlsh:semver> INSERT INTO db (id, key, val) VALUES ('04ed14d5-f0dd-4b4c-81b3-635563760da2', 'archive/0.1.0', '0.1.0');
cqlsh:semver> INSERT INTO db (id, key, val) VALUES ('04ed14d5-f0dd-4b4c-81b3-635563760da2', 'archive/2.1.0', '2.1.0');
cqlsh:semver> SELECT COUNT(*) FROM db WHERE id = '04ed14d5-f0dd-4b4c-81b3-635563760da2'
          ... AND key in ('version', 'archive/0.0.2', 'archive/0.1.0', 'archive/2.1.0');
 count
-------
     4

cqlsh:semver> SELECT * FROM db WHERE id = '04ed14d5-f0dd-4b4c-81b3-635563760da2'
          ... AND key in ('version', 'archive/0.0.2', 'archive/0.1.0', 'archive/2.1.0');

 id                                   | key           | val
--------------------------------------+---------------+-------
 04ed14d5-f0dd-4b4c-81b3-635563760da2 | version       | 0.0.1
 04ed14d5-f0dd-4b4c-81b3-635563760da2 | archive/0.0.2 | 0.0.2
 04ed14d5-f0dd-4b4c-81b3-635563760da2 | archive/0.1.0 | 0.1.0
 04ed14d5-f0dd-4b4c-81b3-635563760da2 | archive/2.1.0 | 2.1.0

```
