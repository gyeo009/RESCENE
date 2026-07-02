package com.example.rescene.core.scenario.scenariogoal.domain;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.AssertionsForClassTypes.assertThatThrownBy;

public class ScenarioGoalTest {
    @Test
    @DisplayName("시나리오 목표는 유효한 값으로 생성된다")
    void createScenarioGoal() {
        ScenarioGoal goal = ScenarioGoal.create(
                "DB 커넥션 풀 고갈",
                "커넥션 풀이 고갈된 장애를 해결한다",
                "애플리케이션 로그와 DB 커넥션 풀 설정을 확인해 문제를 해결한다",
                "DATABASE",
                "BEGINNER"
        );

        assertThat(goal.getCategory()).isEqualTo("DATABASE");
    }

    @Test
    @DisplayName("시나리오 목표는 카테고리 없이 생성될 수 없다")
    void rejectBlankCategory() {
        assertThatThrownBy(() -> ScenarioGoal.create(
                "DB 커넥션 풀 고갈",
                "커넥션 풀이 고갈된 장애를 해결한다",
                "애플리케이션 로그와 DB 커넥션 풀 설정을 확인해 문제를 해결한다",
                "",
                "BEGINNER"
        )).isInstanceOf(IllegalArgumentException.class);
    }
}
