package com.rescene.troubleshooting.scenario;

import java.util.List;
import java.util.Optional;

import org.springframework.stereotype.Component;

@Component
public class ScenarioCatalog {

	private final List<ScenarioSummary> scenarios = List.of(
			new ScenarioSummary(
					"db-connection-pool-exhaustion",
					"DB connection pool exhaustion",
					Difficulty.EASY,
					"database",
					30,
					"Requests slow down after leaked database connections consume the pool.",
					List.of("HikariCP", "connection lifecycle", "thread dump", "Actuator metrics")
			),
			new ScenarioSummary(
					"queue-backlog",
					"Async job queue backlog",
					Difficulty.MEDIUM,
					"messaging",
					45,
					"Background jobs pile up when a consumer is slower than the producer.",
					List.of("RabbitMQ", "consumer concurrency", "retry", "dead-letter queue")
			),
			new ScenarioSummary(
					"slow-query",
					"Slow query under load",
					Difficulty.MEDIUM,
					"database",
					45,
					"An endpoint looks fine locally but collapses under realistic data volume.",
					List.of("JPA", "index", "query plan", "load test")
			)
	);

	public List<ScenarioSummary> findAll() {
		return scenarios;
	}

	public Optional<ScenarioSummary> findById(String id) {
		return scenarios.stream()
				.filter(scenario -> scenario.id().equals(id))
				.findFirst();
	}
}
