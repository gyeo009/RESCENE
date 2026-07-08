package com.example.rescene.core.scenario.scenariogoal.controller;

import com.example.rescene.core.scenario.scenariogoal.domain.ScenarioGoal;
import com.example.rescene.core.scenario.scenariogoal.dto.ScenarioGoalResponse;
import com.example.rescene.core.scenario.scenariogoal.service.ScenarioGoalService;
import com.example.rescene.core.scenario.scenariorun.dto.CreateScenarioRunRequest;
import com.example.rescene.core.scenario.scenariorun.dto.ScenarioRunResponse;
import com.example.rescene.core.scenario.scenariorun.service.ScenarioRunService;
import java.net.URI;
import java.util.List;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ScenarioGoalController {
    private final ScenarioGoalService scenarioGoalService;
    private final ScenarioRunService scenarioRunService;

    public ScenarioGoalController(
            ScenarioGoalService scenarioGoalService, ScenarioRunService scenarioRunService) {
        this.scenarioGoalService = scenarioGoalService;
        this.scenarioRunService = scenarioRunService;
    }

    @GetMapping("/api/scenario-goals")
    public ResponseEntity<List<ScenarioGoalResponse>> getScenarioGoals() {
        List<ScenarioGoalResponse> response =
                scenarioGoalService.getScenarioGoals().stream().map(this::toResponse).toList();

        return ResponseEntity.ok(response);
    }

    @GetMapping("/api/scenario-goals/{id}")
    public ResponseEntity<ScenarioGoalResponse> getScenarioGoal(@PathVariable("id") Long id) {
        return ResponseEntity.ok(toResponse(scenarioGoalService.getScenarioGoal(id)));
    }

    @PostMapping("/api/scenario-goals/{id}/runs")
    public ResponseEntity<ScenarioRunResponse> createScenarioRun(
            @PathVariable("id") Long id, @RequestBody CreateScenarioRunRequest request) {
        ScenarioRunResponse response = scenarioRunService.startScenarioRun(request.userId(), id);
        return ResponseEntity.created(URI.create("/api/scenario-runs/" + response.id()))
                .body(response);
    }

    private ScenarioGoalResponse toResponse(ScenarioGoal scenarioGoal) {
        return new ScenarioGoalResponse(
                scenarioGoal.getId(),
                scenarioGoal.getTitle(),
                scenarioGoal.getSummary(),
                scenarioGoal.getDescription(),
                scenarioGoal.getCategory(),
                scenarioGoal.getDifficulty());
    }
}
