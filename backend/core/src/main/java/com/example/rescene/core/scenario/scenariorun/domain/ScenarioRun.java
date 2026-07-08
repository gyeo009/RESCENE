package com.example.rescene.core.scenario.scenariorun.domain;

import com.example.rescene.common.jpa.entity.AuditableEntity;
import com.example.rescene.core.scenario.scenariogoal.domain.ScenarioGoal;
import com.example.rescene.core.user.domain.UserEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.FetchType;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import java.time.Instant;
import lombok.AccessLevel;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Entity
@Getter
@NoArgsConstructor(access = AccessLevel.PROTECTED)
@Table(name = "scenario_runs")
public class ScenarioRun extends AuditableEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "user_id", nullable = false)
    private UserEntity user;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "scenario_id", nullable = false)
    private ScenarioGoal scenarioGoal;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 50)
    private ScenarioRunStatus status;

    @Column(name = "started_at", nullable = false)
    private Instant startedAt;

    @Column(name = "completed_at")
    private Instant completedAt;

    @Column(name = "last_accessed_at", nullable = false)
    private Instant lastAccessedAt;

    private ScenarioRun(UserEntity user, ScenarioGoal scenarioGoal, Instant startedAt) {
        validateRequired(user, "user");
        validateRequired(scenarioGoal, "scenarioGoal");
        validateRequired(startedAt, "startedAt");

        this.user = user;
        this.scenarioGoal = scenarioGoal;
        this.status = ScenarioRunStatus.ACTIVE;
        this.startedAt = startedAt;
        this.lastAccessedAt = startedAt;
    }

    public static ScenarioRun start(UserEntity user, ScenarioGoal scenarioGoal, Instant startedAt) {
        return new ScenarioRun(user, scenarioGoal, startedAt);
    }

    private static void validateRequired(Object value, String fieldName) {
        if (value == null) {
            throw new IllegalArgumentException("Scenario run " + fieldName + " is required");
        }
    }
}
