# Goal for this is to understand replication

## NOTE
There several ways of handling read-replicas
1. Application Layer, this is where you manually route the db requests to the appropraite replica in the application logic. This issue with this is that it application routing logic can be complex and error-prone.


2. Database Layer, this is where you configure the database to route read requests to the appropriate replica. You can use tools like pgpool-II, pgBouncer or Vitess to achieve this. The issue with this is that it can be complex to set up and maintain, and it may not be suitable for all use cases.




## Replication Types

### synchronous replication
The synchronous flow is
1. send api request to master
2. ensure write happens, then send to save within the same response
3. the send response to client.


### asynchronous replication
The asynchronous flow is
1. send api request to master 
2. then send response to client.
3. watch how replication is carried out in postgres

A docker compose is used for this.


---
Something I saw online.
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
