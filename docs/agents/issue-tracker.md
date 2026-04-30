# Local Markdown Issue Tracking

TeleRiHa uses a local, file-based issue tracking system for agentic workflows. This allows AI agents to track progress, blockers, and tasks directly within the repository.

## Structure

All issues are stored in the `.scratch/` directory, organized by feature or epic:

```
.scratch/
└── <feature-name>/
    ├── issue-1-description.md
    ├── issue-2-description.md
    └── metadata.json
```

## Issue File Format

Each issue file should follow this template:

```markdown
# [ID] Issue Title

**Status**: [open | in-progress | blocked | closed]
**Labels**: [needs-triage, ready-for-agent, etc.]
**Assignee**: [Agent Name | Human]

## Description
A clear and concise description of the problem or task.

## Tasks
- [ ] Sub-task 1
- [ ] Sub-task 2

## Notes
Any relevant implementation details or links.
```

## Workflow

### Creating an Issue
Create a new `.md` file in `.scratch/<feature>/` with the next available ID.

### Updating an Issue
Modify the status, check off tasks, or add notes as progress is made.

### Closing an Issue
Change the **Status** to `closed` and move the file to an `archive/` sub-directory if requested, or simply leave it with the closed status.
