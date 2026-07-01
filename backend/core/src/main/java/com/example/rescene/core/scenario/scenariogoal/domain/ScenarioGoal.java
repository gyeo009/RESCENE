package com.example.rescene.core.scenario.scenariogoal.domain;

import com.example.rescene.common.jpa.entity.AuditableEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AccessLevel;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Entity
@Getter
@NoArgsConstructor(access = AccessLevel.PROTECTED)
@Table(
        name = "scenario_goal"
)
public class ScenarioGoal extends AuditableEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "category", nullable = false, length = 50)
    private String category;

    ScenarioGoal(String category){
        if (category == null || category.isBlank()){
            throw new IllegalArgumentException("Scenario Goal category is required");
        }

        this.category = category;
    }

}
