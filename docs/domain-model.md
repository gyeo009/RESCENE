# RESCENE Domain Model

## Domains

- User
  - RESCENE learner record.
- Scenario
  - Troubleshooting practice template that reproduces an operational failure.
- ScenarioSpec
  - Execution specification for a Scenario.
  - Defines image, source, resource limits, timeout, and validation entrypoint.
- ScenarioRun
  - One attempt/run of a Scenario by a User.
- Sandbox
  - Container-based virtual environment used to run a ScenarioRun.

## Domain Rules

- One User can have many ScenarioRuns.
- One Scenario can have many ScenarioSpecs.
- One Scenario can have many ScenarioRuns.
- One ScenarioRun can have zero or one Sandbox.
- One ScenarioSpec can be used by many Sandboxes.
- One Sandbox uses exactly one ScenarioSpec.
- One User can have only one active ScenarioRun for the same Scenario.

