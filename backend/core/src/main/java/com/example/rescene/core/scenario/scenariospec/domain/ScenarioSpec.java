package com.example.rescene.core.scenario.scenariospec.domain;

import com.example.rescene.common.jpa.entity.AuditableEntity;
import com.example.rescene.core.scenario.scenariogoal.domain.ScenarioGoal;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import lombok.AccessLevel;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Entity
@Getter
@NoArgsConstructor(access = AccessLevel.PROTECTED)
@Table(name = "scenario_specs")
public class ScenarioSpec extends AuditableEntity {
	@Id
	@GeneratedValue(strategy = GenerationType.IDENTITY)
	private Long id;

	@ManyToOne(fetch = FetchType.LAZY, optional = false)
	@JoinColumn(name = "scenario_id", nullable = false)
	private ScenarioGoal scenarioGoal;

	@Column(name = "image_ref", nullable = false, length = 255)
	private String imageRef;

	private ScenarioSpec(ScenarioGoal scenarioGoal, String imageRef) {
		validateRequired(scenarioGoal, "scenarioGoal");
		validateRequired(imageRef, "imageRef");

		this.scenarioGoal = scenarioGoal;
		this.imageRef = imageRef;
	}

	public static ScenarioSpec create(ScenarioGoal scenarioGoal, String imageRef) {
		return new ScenarioSpec(scenarioGoal, imageRef);
	}

	private static void validateRequired(Object value, String fieldName) {
		if (value == null) {
			throw new IllegalArgumentException("Scenario spec " + fieldName + " is required");
		}
	}

	private static void validateRequired(String value, String fieldName) {
		if (value == null || value.isBlank()) {
			throw new IllegalArgumentException("Scenario spec " + fieldName + " is required");
		}
	}
}
