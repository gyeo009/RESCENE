package com.rescene.troubleshooting.scenario;

import java.util.List;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/scenarios")
public class ScenarioController {

	private final ScenarioCatalog scenarioCatalog;

	public ScenarioController(ScenarioCatalog scenarioCatalog) {
		this.scenarioCatalog = scenarioCatalog;
	}

	@GetMapping
	public List<ScenarioSummary> listScenarios() {
		return scenarioCatalog.findAll();
	}

	@GetMapping("/{id}")
	public ResponseEntity<ScenarioSummary> getScenario(@PathVariable String id) {
		return scenarioCatalog.findById(id)
				.map(ResponseEntity::ok)
				.orElseGet(() -> ResponseEntity.notFound().build());
	}
}
