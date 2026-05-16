# AGENTS.md — db-operator

Project guidance for AI agents. Canonical project-level instructions; the
centralized `~/.claude/AGENTS.md` governs cross-project rules.

## Project facts

- Kubernetes operator (kubebuilder v4) provisioning databases: postgres, mysql,
  clickhouse. Module `github.com/db-operator/db-operator/v2`.
- Build system is **Taskfile.yml** (`task`), not `make`. Key targets:
  `build`, `generate`, `manifests`, `lint`, `unit-test`,
  `postgres-test`, `mysql-test`, `clickhouse-test`, `clickhouse-cluster-test`.
- `vendor/` is git-ignored and shared; run `go mod vendor` after dependency
  changes.
- Integration tests run against a kind cluster provisioned by
  `helmfile.yaml.gotmpl` + `test/charts/`. ClickHouse uses the Altinity operator.

---

# Multi-Agent Coordination Protocol (v1)

Multiple agents work in this repository **concurrently and without git
worktrees** — they share one working tree. This protocol keeps them from
colliding, going stale, or wasting tokens. It is mandatory for every agent.

## 1. The task board

The single source of truth for who-is-doing-what is
[`.agents/coordination/tasks.md`](.agents/coordination/tasks.md). It is a
committed file. Every task is one row:

```
| id | title | owner | status | claimed_at | lease_expires_at | updated_at | notes |
```

- **id** — stable short slug (`ch-cluster-test`).
- **owner** — agent id (`agent:<host>-<pid>` or a human-given name); `-` if unclaimed.
- **status** — `todo` | `doing` | `blocked` | `review` | `done`.
- **timestamps** — UTC ISO-8601 (`2026-05-16T18:40:00Z`). Never use relative time.
- **lease_expires_at** — when the claim expires (see §3).

## 2. Lifecycle — every agent, every task

1. **Sync first.** `git pull` (or fetch+rebase) before reading the board.
2. **Pick** the oldest `todo`, or a task whose lease is **expired** (§3).
3. **Claim** — set `owner`, `status=doing`, `claimed_at=now`,
   `lease_expires_at=now+30m`, `updated_at=now`. Commit the board change
   *before* touching code: `git commit -m "[agent:<id>] claim: <id>"`.
4. **Heartbeat** — at least every 20 min of active work, bump
   `lease_expires_at` and `updated_at` and commit. No heartbeat ⇒ stale.
5. **Finish** — set `status=review` or `done`, `updated_at=now`, commit.
6. **Block** — if stuck, set `status=blocked`, write the reason in `notes`,
   release the lease (`owner=-`), commit. Never sit silently on a blocked task.

## 3. Anti-stale (task expiry)

A task is **stale** when `status=doing` and `now > lease_expires_at`. Any agent
may reclaim a stale task: set yourself as `owner`, fresh `lease_expires_at`, and
record `notes: reclaimed from <old owner> (stale)`. Before reclaiming, check
`git log` for that task's commits to recover context. An agent that finds its
own lease expired must assume another agent may have taken over — re-sync.

An **idle** agent (no claimed task) must not stop: it picks the next `todo`,
or, if none, scans for stale tasks, failing PRs, or untested changes.

## 4. Communication standard

- Agents communicate **through git**, not chat. Every claim, heartbeat, block,
  and handoff is a board commit.
- Commit message prefix: `[agent:<id>] <verb>: <task-id> — <summary>`.
  Verbs: `claim`, `heartbeat`, `progress`, `block`, `handoff`, `done`.
- Work commits keep the project's conventional style
  (`feat(...)`, `fix(...)`, `chore(...)`); prepend `[agent:<id>]` only on
  board/coordination commits to keep history readable.
- Leave decisions and non-obvious rationale in the task `notes` or the commit
  body so the next agent needs no out-of-band context.

## 5. Token economy

- **State lives in git, not context.** Commit small and often so a fresh agent
  reconstructs state from `git log` + the board, not from a long transcript.
- Read the board and only the files for your task. Use `grep`/symbol tools, not
  whole-file reads. Delegate broad searches to sub-agents.
- Don't re-read a file you just wrote; don't re-run a green gate.
- Prefer one focused task to completion over many half-done ones.

## 6. Shared-tree discipline (no worktrees)

- `git pull` before editing; `git status` before each edit to detect another
  agent's in-flight changes to the same file — if found, coordinate via the
  board, don't overwrite.
- **Commit with high frequency**, in **surgical, targeted** changes. Never
  leave the tree broken: every commit must `task build` clean.
- Forbidden: `git reset --hard`, `git push --force`, `git clean -fd`,
  discarding another agent's uncommitted work.
- Push after every few commits so other agents and machines stay in sync.

## 7. Definition of done

A task is `done` only when the change builds, the relevant gates pass
(`task lint`, and `task unit-test` / engine integration tests where
applicable), no hacks or silenced errors were introduced, and the board row is
updated. Honest status only — a partially working change stays `review` or
`blocked` with the gap written in `notes`.
