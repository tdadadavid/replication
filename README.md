# The Goal for this is to understand replication

Check out the readme
- [async](./internal/dbasync/README.md)
- [sync](./internal/dbsync/README.md)

## Replication Types

1. Synchronous Replication
The flow includes
   1. send a `POST` request to master on endpoint `sync-users` ensure that write happens
   2. Next, send it to all followers within the same response before responding to the client
   3. Prometheus metrics are registered to track important information


1. Asynchronous Replication
The asynchronous flow is
   1. send a `POST` request to master on endpoint `async-users` ensure that write happens, then respond to a client immediately
   2. watch how replication is carried out in postgres. Necessary configuration is [here](./internal/dbasync/README.md) 

## NOTE
There are several ways of handling read-replicas
1. Application Layer, this is where you manually route the db requests to the appropriate replica in the application logic. This issue with this is that its application routing logic can be complex and error-prone.
2. Database Layer, this is where you configure the database to route read requests to the appropriate replica. You can use tools like pgpool-II, pgBouncer or Vitess to achieve this. The issue with this is that it can be complex to set up and maintain, and it may not be suitable for all use cases.

---
Something I saw online, just dumping it here.
<br>
```sql
--- Verifying replication lag
SELECT
    application_name,
    client_addr,
    state,
    sent_lsn,
    write_lsn,
    flush_lsn,
    replay_lsn,
    pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn)) AS replication_lag
FROM
    pg_stat_replication;
```

## Resources
1. [PostgreSQL Replication - TigerData](https://www.tigerdata.com/learn/best-practices-for-postgres-database-replication)
2. [PostgreSQL Replication - Reddit](https://www.reddit.com/r/PostgreSQL/comments/1krsrnb/comment/mtg4gpn/?force-legacy-sct=1)
3. [PostgreSQL Replication - Linkedin](https://www.linkedin.com/pulse/streaming-replication-postgresql-17-ubuntu-2404-lts-step-mahto-ww0wf/)
