package com.rescene.troubleshooting.scenario;

import static org.hamcrest.Matchers.hasSize;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.test.web.servlet.MockMvc;

@WebMvcTest(ScenarioController.class)
@Import(ScenarioCatalog.class)
class ScenarioControllerTests {

	@Autowired
	private MockMvc mockMvc;

	@Test
	void listsScenarios() throws Exception {
		mockMvc.perform(get("/api/scenarios"))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$", hasSize(3)))
				.andExpect(jsonPath("$[0].id").value("db-connection-pool-exhaustion"));
	}

	@Test
	void returnsScenarioById() throws Exception {
		mockMvc.perform(get("/api/scenarios/queue-backlog"))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.title").value("Async job queue backlog"));
	}

	@Test
	void returnsNotFoundForUnknownScenario() throws Exception {
		mockMvc.perform(get("/api/scenarios/nope"))
				.andExpect(status().isNotFound());
	}
}
