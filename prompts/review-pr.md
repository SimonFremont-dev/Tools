# Prompt template: review a pull request

Review the PR diff for {{PROJECT}} (language: {{LANGUAGE}}).

Check, in order:
1. Does it follow the project conventions (layout, AGENTS.md, Makefile targets)?
2. Are all six CI gates addressed (pre-checks, format, lint, security, test, build)?
3. New dependencies: are they justified, and do they have conventions in the
   Tools repo (templates/)? If not, flag them as coverage alerts.
4. Secrets, credentials, or sensitive data in the diff?
5. Error handling and tests for the changed behavior.

Output: a short verdict (approve / request changes), then a numbered list of
findings with file:line references, most severe first.
