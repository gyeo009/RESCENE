package com.example.rescene.core.scenario.scenariogoal.dto;

public record ScenarioGoalResponse(
        Long id,
        String title,
        String summary,
        String description,
        String category,
        String difficulty) {}
