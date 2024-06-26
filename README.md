<details>
<summary>
Notes #8
</summary>

```sql
-- Tx1: Transfer $10 from account 1 to 2
BEGIN;

UPDATE accounts SET balance = balance - 10 where id = 1 RETURNING *;
UPDATE accounts SET balance = balance + 10 where id = 2 RETURNING *;

ROLLBACK;

-- Tx2: Transer $10 from account 2 to 1
BEGIN;

UPDATE accounts SET balance = balance - 10 where id = 2 RETURNING *;
UPDATE accounts SET balance = balance + 10 where id = 1 RETURNING *;

ROLLBACK;
```

Use SQL to debug from Notes #7 to see how txn2's exclusive lock blocks txn1 2nd update command
and continuing to run tx2's second command would end up in a dead lock

Solution:
change txn2 to update id = 1 first

```sql
UPDATE accounts SET balance = balance + 10 where id = 1 RETURNING *;
UPDATE accounts SET balance = balance - 10 where id = 2 RETURNING *;
```

This way when tx1 is initiated it holds a exclusive lock for id 1 and tx2 is put on hold from the first command itself.
This way tx1 can complete with its second command without a hold by tx2 and thus dead lock is resolved

To summarize: When we follow an order in updating the accounts we can avoid a deadlock, here in our case account with smaller id is updated in any txn

</details>

<details>
<summary>
Notes #7
</summary>

```sql
-- SQL to create deadlock
BEGIN;

INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES (1, 2, 10) RETURNING \*;

INSERT INTO entries (account*id, amount) VALUES (1, -10) RETURNING *;
INSERT INTO entries (account*id, amount) VALUES (2, 10) RETURNING *;

SELECT * from accounts where id = 1 FOR UPDATE;
UPDATE accounts SET balance = 90 where id = 1 RETURNING _;

SELECT * from accounts where id = 2 FOR UPDATE;
UPDATE accounts SET balance = 110 where id = 2 RETURNING;

ROLLBACK;
```

<details>
<summary>SQL to debug deadlock: </summary>

[Reference](https://wiki.postgresql.org/wiki/Lock_Monitoring)

```sql
SELECT
a.application_name,
l.relation::regclass,
l.transactionid,
l.mode,
l.locktype,
l.GRANTED,
a.usename,
a.query,
a.pid
FROM pg_stat_activity a
JOIN pg_locks l ON l.pid = a.pid
where application_name='psql'
ORDER BY a.pid;
```

<details>
<summary>Debug More</summary>

```sql
SELECT blocked_locks.pid AS blocked_pid,
blocked_activity.usename AS blocked_user,
blocking_locks.pid AS blocking_pid,
blocking_activity.usename AS blocking_user,
blocked_activity.query AS blocked_statement,
blocking_activity.query AS current_statement_in_blocking_process
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks
ON blocking_locks.locktype = blocked_locks.locktype
AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
AND blocking_locks.pid != blocked_locks.pid

    JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid

WHERE NOT blocked_locks.granted;
```

</details>

</details>

</details>
