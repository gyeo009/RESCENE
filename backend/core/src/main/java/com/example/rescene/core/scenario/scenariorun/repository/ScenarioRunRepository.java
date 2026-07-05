package com.example.rescene.core.scenario.scenariorun.repository;

import com.example.rescene.core.scenario.scenariorun.domain.ScenarioRun;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ScenarioRunRepository extends JpaRepository<ScenarioRun, Long> {
}
