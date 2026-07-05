package com.example.rescene.core.scenario.scenariogoal.domain;

import com.example.rescene.common.jpa.entity.AuditableEntity;
import jakarta.persistence.*;
import lombok.AccessLevel;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Entity
@Getter
@NoArgsConstructor(access = AccessLevel.PROTECTED)
@Table(name = ScenarioGoal.TABLE_NAME)
public class ScenarioGoal extends AuditableEntity {

    public static final String TABLE_NAME = "scenarios";

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, length = 100)
    private String title;

    @Column(nullable = false, length = 255)
    private String summary;

    @Column(nullable = false, columnDefinition = "text")
    private String description;

    @Column(nullable = false, length = 50)
    private String category;

    @Column(nullable = false, length = 50)
    private String difficulty;

    private ScenarioGoal(String title, String summary, String description, String category, String difficulty) {
        validateRequired(title, "title");
        validateRequired(summary, "summary");
        validateRequired(description, "description");
        validateRequired(category, "category");
        validateRequired(difficulty, "difficulty");

        this.title = title;
        this.summary = summary;
        this.description = description;
        this.category = category;
        this.difficulty = difficulty;
    }

    private static void validateRequired(String value, String fieldName) {
        if(value == null || value.isBlank()){
            throw new IllegalArgumentException("Scenario Goal " + fieldName + " is required");
        }
    }

    public static ScenarioGoal create(String title, String summary, String description, String category, String difficulty) {
        return new ScenarioGoal(title, summary, description, category, difficulty);
    }
}
