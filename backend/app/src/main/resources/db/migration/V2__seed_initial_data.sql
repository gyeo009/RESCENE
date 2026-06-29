insert into users (id, display_name, created_at, updated_at)
values
	(1, 'test_user', now(), now());

insert into scenarios (id, title, summary, description, category, difficulty, created_at, updated_at)
values
	(
		1,
		'DB connection pool exhaustion',
		'Requests slow down after leaked database connections consume the pool.',
		'Investigate a Spring Boot service where leaked database connections eventually exhaust the HikariCP pool and stall incoming requests.',
		'SPRING_BOOT_RUNTIME',
		'BEGINNER',
		now(),
		now()
	),
	(
		2,
		'Slow query / N+1',
		'A normal-looking endpoint collapses under realistic data volume.',
		'Find and fix a slow persistence path caused by inefficient query behavior under load.',
		'DATABASE',
		'INTERMEDIATE',
		now(),
		now()
	),
	(
		3,
		'Async job queue backlog',
		'Background jobs pile up when the consumer cannot keep up with producers.',
		'Diagnose a queue-backed workflow where job processing falls behind and system health degrades.',
		'QUEUE',
		'INTERMEDIATE',
		now(),
		now()
	);

insert into scenario_specs (id, scenario_id, image_ref, created_at, updated_at)
values
	(1, 1, 'ghcr.io/gyeo009/rescene-scenario-db-pool:latest', now(), now()),
	(2, 2, 'ghcr.io/gyeo009/rescene-scenario-slow-query:latest', now(), now()),
	(3, 3, 'ghcr.io/gyeo009/rescene-scenario-queue-backlog:latest', now(), now());

select setval(pg_get_serial_sequence('users', 'id'), (select max(id) from users));
select setval(pg_get_serial_sequence('scenarios', 'id'), (select max(id) from scenarios));
select setval(pg_get_serial_sequence('scenario_specs', 'id'), (select max(id) from scenario_specs));
