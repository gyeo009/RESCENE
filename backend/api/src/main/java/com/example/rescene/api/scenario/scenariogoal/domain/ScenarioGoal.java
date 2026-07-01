package com.example.rescene.api.scenario.scenariogoal.domain;

import com.example.rescene.common.jpa.entity.AuditableEntity;
import jakarta.persistence.*;

@Entity
@Table(
        name = "scenario_goal"
)
public class ScenarioGoal extends AuditableEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;


    @Column(name = "category", nullable = false, length = 50)
    private String category;

    // TODO: default category는 어떤 도메인 규약인가? 확실히 정리 필요
    public ScenarioGoal() {
        this.category = "default";
    }

    public Long getId() {
        return id;
    }

    public String getCategory(){
        return category;
    }

    ScenarioGoal(String category){
        if (category == null || category.isBlank()){
            throw new IllegalArgumentException("Scenario Goal category is required");
        }

        this.category = category;
    }

}
