# Triage Labels

To facilitate collaboration between humans and AI agents, TeleRiHa uses a set of standardized triage labels in `.scratch` issues.

## Available Labels

| Label | Description | Use Case |
|-------|-------------|----------|
| `needs-triage` | Initial state for all new issues. | A new task is identified but not yet analyzed. |
| `needs-info` | Information is missing to proceed. | The agent requires clarification on requirements or a bug report is incomplete. |
| `ready-for-agent` | The task is well-defined and can be handled by an AI. | Implementation tasks with clear scope. |
| `ready-for-human` | Requires human intervention, review, or high-level decision. | Complex architectural choices or final validation. |
| `wontfix` | The issue is invalid or will not be addressed. | Duplicate issues or out-of-scope requests. |

## Applying Labels

Labels should be added to the **Labels** field in the front matter of an issue file:

```markdown
# [ID] Issue Title
**Status**: open
**Labels**: needs-triage, ready-for-agent
```

Agents should actively update these labels as they analyze and progress through issues.
