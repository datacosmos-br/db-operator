# Task board — db-operator

Shared coordination board. See [AGENTS.md](../../AGENTS.md) §"Multi-Agent
Coordination Protocol" for the rules. Timestamps are UTC ISO-8601.

| id | title | owner | status | claimed_at | lease_expires_at | updated_at | notes |
|----|-------|-------|--------|------------|------------------|------------|-------|
| ch-feature | Advanced ClickHouse support (cluster, Replicated, RBAC, driver v2) | - | review | 2026-05-16T14:00:00Z | 2026-05-16T19:00:00Z | 2026-05-16T18:50:00Z | Phases 0-6 committed (driver v2, ON CLUSTER+Replicated, AccessType grants, quotas/profiles, webhook validation, Altinity-operator test charts). Single-node validated on CH 24.3/26.3; cluster validated on local Docker 26.3. Pending: kind cluster integration green (dbo_admin auth nuance under investigation), `task clickhouse-test` single-node in kind, plan file refresh. |
| ch-cluster-kind | Validate clickhouse-cluster chart end-to-end in kind | - | blocked | - | - | 2026-05-16T18:50:00Z | dbo_admin user created by Altinity (correct sha256) but client auth fails inside the kind pods; likely an Altinity users.d reload / config nuance. Needs a surgical fix to test/charts/clickhouse-cluster or the CHI users block. |
