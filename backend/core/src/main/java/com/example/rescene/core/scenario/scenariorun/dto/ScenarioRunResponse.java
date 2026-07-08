package com.example.rescene.core.scenario.scenariorun.dto;

import java.time.Instant;

public record ScenarioRunResponse(
        Long id,
        Long userId,
        Long scenarioGoalId,
        String status,
        Instant startedAt,
        Instant completedAt,
        Instant lastAccessedAt) {}
