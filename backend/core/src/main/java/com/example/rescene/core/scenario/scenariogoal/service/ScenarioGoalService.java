package com.example.rescene.core.scenario.scenariogoal.service;

import java.util.List;

import com.example.rescene.core.scenario.scenariogoal.domain.ScenarioGoal;
import com.example.rescene.core.scenario.scenariogoal.repository.ScenarioGoalRepository;
import org.springframework.stereotype.Service;

@Service
public class ScenarioGoalService {
	private final ScenarioGoalRepository scenarioGoalRepository;

	public ScenarioGoalService(ScenarioGoalRepository scenarioGoalRepository) {
		this.scenarioGoalRepository = scenarioGoalRepository;
	}

	public List<ScenarioGoal> getScenarioGoals() {
		return scenarioGoalRepository.findAll();
	}

	public ScenarioGoal getScenarioGoal(Long id) {
		return scenarioGoalRepository.findById(id)
				.orElseThrow(() -> new IllegalArgumentException("Scenario goal not found: " + id));
	}
}
