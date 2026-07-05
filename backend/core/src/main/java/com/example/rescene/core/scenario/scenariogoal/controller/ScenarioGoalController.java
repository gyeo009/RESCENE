package com.example.rescene.core.scenario.scenariogoal.controller;

import java.util.List;

import com.example.rescene.core.scenario.scenariogoal.domain.ScenarioGoal;
import com.example.rescene.core.scenario.scenariogoal.dto.ScenarioGoalResponse;
import com.example.rescene.core.scenario.scenariogoal.service.ScenarioGoalService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ScenarioGoalController {
	private final ScenarioGoalService scenarioGoalService;

	public ScenarioGoalController(ScenarioGoalService scenarioGoalService) {
		this.scenarioGoalService = scenarioGoalService;
	}

	@GetMapping("/api/scenarios")
	public ResponseEntity<List<ScenarioGoalResponse>> getScenarioGoals() {
		List<ScenarioGoalResponse> response = scenarioGoalService.getScenarioGoals().stream()
				.map(this::toResponse)
				.toList();

		return ResponseEntity.ok(response);
	}

	@GetMapping("/api/scenarios/{id}")
	public ResponseEntity<ScenarioGoalResponse> getScenarioGoal(@PathVariable("id") Long id) {
		return ResponseEntity.ok(toResponse(scenarioGoalService.getScenarioGoal(id)));
	}

	private ScenarioGoalResponse toResponse(ScenarioGoal scenarioGoal) {
		return new ScenarioGoalResponse(
				scenarioGoal.getId(),
				scenarioGoal.getTitle(),
				scenarioGoal.getSummary(),
				scenarioGoal.getDescription(),
				scenarioGoal.getCategory(),
				scenarioGoal.getDifficulty()
		);
	}
}
