package com.example.rescene.core.sandbox.domain;

import com.example.rescene.common.jpa.entity.AuditableEntity;
import com.example.rescene.core.scenario.scenariorun.domain.ScenarioRun;
import com.example.rescene.core.scenario.scenariospec.domain.ScenarioSpec;
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
import jakarta.persistence.OneToOne;
import jakarta.persistence.Table;
import lombok.AccessLevel;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Entity
@Getter
@NoArgsConstructor(access = AccessLevel.PROTECTED)
@Table(name = "sandboxes")
public class Sandbox extends AuditableEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @OneToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "scenario_run_id", nullable = false, unique = true)
    private ScenarioRun scenarioRun;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "scenario_spec_id", nullable = false)
    private ScenarioSpec scenarioSpec;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 50)
    private SandboxStatus status;

    @Column(name = "workspace_path", nullable = false, length = 500)
    private String workspacePath;

    private Sandbox(ScenarioRun scenarioRun, ScenarioSpec scenarioSpec, String workspacePath) {
        validateRequired(scenarioRun, "scenarioRun");
        validateRequired(scenarioSpec, "scenarioSpec");
        validateRequired(workspacePath, "workspacePath");

        this.scenarioRun = scenarioRun;
        this.scenarioSpec = scenarioSpec;
        this.status = SandboxStatus.PROVISIONING;
        this.workspacePath = workspacePath;
    }

    public static Sandbox provision(
            ScenarioRun scenarioRun, ScenarioSpec scenarioSpec, String workspacePath) {
        return new Sandbox(scenarioRun, scenarioSpec, workspacePath);
    }

    private static void validateRequired(Object value, String fieldName) {
        if (value == null) {
            throw new IllegalArgumentException("Sandbox " + fieldName + " is required");
        }
    }

    private static void validateRequired(String value, String fieldName) {
        if (value == null || value.isBlank()) {
            throw new IllegalArgumentException("Sandbox " + fieldName + " is required");
        }
    }
}
