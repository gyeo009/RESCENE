package com.rescene.troubleshooting.scenario;

import java.util.List;

public record ScenarioSummary(
		String id,
		String title,
		Difficulty difficulty,
		String area,
		int estimatedMinutes,
		String description,
		List<String> concepts
) {
}
