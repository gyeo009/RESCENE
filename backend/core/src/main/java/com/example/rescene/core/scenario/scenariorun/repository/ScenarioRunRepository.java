package com.example.rescene.core.scenario.scenariorun.repository;

import com.example.rescene.core.scenario.scenariorun.domain.ScenarioRun;
import com.example.rescene.core.scenario.scenariorun.domain.ScenarioRunStatus;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ScenarioRunRepository extends JpaRepository<ScenarioRun, Long> {
    boolean existsByUser_IdAndScenarioGoal_IdAndStatus(
            Long userId, Long scenarioGoalId, ScenarioRunStatus status);
}
