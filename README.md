
### Goal for this is to understand replication
1. synchronous replication
The synchronous flow is
1. send api request to master
2. ensure write happens, then send to save within the same response
3. the send response to client.

2. asynchronous replication
The asynchronous flow is
1. send api request to master 
2. then send response to client.
3. watch how replication is carried out in postgres

A docker compose is used for this.
