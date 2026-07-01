package com.example.rescene.api.scenario.scenariogoal.domain;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.AssertionsForClassTypes.assertThatThrownBy;
import static org.junit.jupiter.api.Assertions.assertNotNull;

public class ScenarioGoalTest {
    @Test
    @DisplayName("시나리오 목표는 유효한 카테고리가 지정되어있어야 한다")
    void scenarioGoalRequiresCategory(){
        // Given
        String category = "DB Connection pool exhaustion을 해결한다";

        // When
        ScenarioGoal sg = new ScenarioGoal(
                category
        );

        // Then
        assertNotNull(sg.getCategory());
        assertThatThrownBy(() -> new ScenarioGoal(null))
                .isInstanceOf(IllegalArgumentException.class);
        assertThatThrownBy(() -> new ScenarioGoal(" "))
                .isInstanceOf(IllegalArgumentException.class);
    }
}
