package com.example.rescene.core.scenario.scenariospec.repository;

import com.example.rescene.core.scenario.scenariospec.domain.ScenarioSpec;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ScenarioSpecRepository extends JpaRepository<ScenarioSpec, Long> {
}
