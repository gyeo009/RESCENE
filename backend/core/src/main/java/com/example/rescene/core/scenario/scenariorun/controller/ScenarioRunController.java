package com.example.rescene.core.scenario.scenariorun.controller;

import com.example.rescene.core.scenario.scenariorun.dto.ScenarioRunResponse;
import com.example.rescene.core.scenario.scenariorun.service.ScenarioRunService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ScenarioRunController {
    private final ScenarioRunService scenarioRunService;

    public ScenarioRunController(ScenarioRunService scenarioRunService) {
        this.scenarioRunService = scenarioRunService;
    }

    @GetMapping("/api/scenario-runs/{id}")
    public ResponseEntity<ScenarioRunResponse> getScenarioRun(@PathVariable("id") Long id) {
        return ResponseEntity.ok(scenarioRunService.getScenarioRun(id));
    }
}
