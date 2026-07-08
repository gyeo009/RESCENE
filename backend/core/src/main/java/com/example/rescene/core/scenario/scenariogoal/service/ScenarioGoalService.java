package com.example.rescene.core.scenario.scenariogoal.service;

import com.example.rescene.core.common.exception.ResourceNotFoundException;
import com.example.rescene.core.scenario.scenariogoal.domain.ScenarioGoal;
import com.example.rescene.core.scenario.scenariogoal.repository.ScenarioGoalRepository;
import java.util.List;
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
        return scenarioGoalRepository
                .findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Scenario goal not found: " + id));
    }
}
