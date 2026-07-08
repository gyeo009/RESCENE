package com.example.rescene.core.scenario.scenariorun.service;

import com.example.rescene.core.common.exception.DuplicateActiveScenarioRunException;
import com.example.rescene.core.common.exception.ResourceNotFoundException;
import com.example.rescene.core.scenario.scenariogoal.domain.ScenarioGoal;
import com.example.rescene.core.scenario.scenariogoal.repository.ScenarioGoalRepository;
import com.example.rescene.core.scenario.scenariorun.domain.ScenarioRun;
import com.example.rescene.core.scenario.scenariorun.domain.ScenarioRunStatus;
import com.example.rescene.core.scenario.scenariorun.dto.ScenarioRunResponse;
import com.example.rescene.core.scenario.scenariorun.repository.ScenarioRunRepository;
import com.example.rescene.core.user.domain.UserEntity;
import com.example.rescene.core.user.repository.UserRepository;
import java.time.Instant;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class ScenarioRunService {
    private final UserRepository userRepository;
    private final ScenarioGoalRepository scenarioGoalRepository;
    private final ScenarioRunRepository scenarioRunRepository;

    public ScenarioRunService(
            UserRepository userRepository,
            ScenarioGoalRepository scenarioGoalRepository,
            ScenarioRunRepository scenarioRunRepository) {
        this.userRepository = userRepository;
        this.scenarioGoalRepository = scenarioGoalRepository;
        this.scenarioRunRepository = scenarioRunRepository;
    }

    @Transactional
    public ScenarioRunResponse startScenarioRun(Long userId, Long scenarioGoalId) {
        validateRequired(userId, "userId");
        validateRequired(scenarioGoalId, "scenarioGoalId");

        UserEntity user =
                userRepository
                        .findById(userId)
                        .orElseThrow(
                                () -> new ResourceNotFoundException("User not found: " + userId));
        ScenarioGoal scenarioGoal =
                scenarioGoalRepository
                        .findById(scenarioGoalId)
                        .orElseThrow(
                                () ->
                                        new ResourceNotFoundException(
                                                "Scenario goal not found: " + scenarioGoalId));

        if (scenarioRunRepository.existsByUser_IdAndScenarioGoal_IdAndStatus(
                userId, scenarioGoalId, ScenarioRunStatus.ACTIVE)) {
            throw new DuplicateActiveScenarioRunException(
                    "Active scenario run already exists for user "
                            + userId
                            + " and scenario goal "
                            + scenarioGoalId);
        }

        ScenarioRun scenarioRun = ScenarioRun.start(user, scenarioGoal, Instant.now());
        return toResponse(scenarioRunRepository.save(scenarioRun));
    }

    @Transactional(readOnly = true)
    public ScenarioRunResponse getScenarioRun(Long id) {
        return scenarioRunRepository
                .findById(id)
                .map(this::toResponse)
                .orElseThrow(() -> new ResourceNotFoundException("Scenario run not found: " + id));
    }

    private ScenarioRunResponse toResponse(ScenarioRun scenarioRun) {
        return new ScenarioRunResponse(
                scenarioRun.getId(),
                scenarioRun.getUser().getId(),
                scenarioRun.getScenarioGoal().getId(),
                scenarioRun.getStatus().name(),
                scenarioRun.getStartedAt(),
                scenarioRun.getCompletedAt(),
                scenarioRun.getLastAccessedAt());
    }

    private static void validateRequired(Object value, String fieldName) {
        if (value == null) {
            throw new IllegalArgumentException("Scenario run " + fieldName + " is required");
        }
    }
}
