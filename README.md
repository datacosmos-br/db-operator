# DB Operator

This operator lets you manage databases in a Kubernetes native way, even if they are not deployed to Kubernetes

## Features

DB Operator provides following features:

* Management of **MySQL** and **PostgreSQL** databases in the same way
* Create/Delete databases on the database server running outside/inside Kubernetes by creating `Database` custom resource;
* Create/Delete users on the database server running outside/inside Kubernetes by creating `DbUser` custom resource;
* Creating of custom connection strings using **GO templates**

## Documentation
* [Get Started](./documentation/docs/index.md)
* [Configure instances](./documentation/docs/dbinstances.md)
* [Manage databases](./documentation/docs/database.md)
* [Manage users](./documentation/docs/dbuser.md)
* [A deeper look at templates](./documentation/docs/templates.md)

## Quickstart

### To install DB Operator with helm:

```
$ helm repo add db-operator https://db-operator.github.io/charts/
$ helm install --name my-release db-operator/db-operator
```

To see more options of helm values, [see the chart repo]([https://github.com/db-operator/charts/tree/main/charts/db-operator])

## Datacosmos Fork

This fork adds the following capabilities on top of upstream `db-operator/db-operator`:

- **ClickHouse engine support** — full RBAC, ON CLUSTER DDL, Replicated engine, quotas, settings profiles, host restrictions, and extra privilege grants.
- **PostgreSQL extras** — `owner`, `postInitSQL`, and `crossDatabaseGrants` in the `Database` CR.
- **Operational fixes** — MySQL user SQL host-part fix, CronJob `APIVersion` fix, Postgres/MySQL query timeouts, work-queue starvation prevention, and ArgoCD-compatible ownership tracking.
- **Release pipeline** — multi-arch container images signed with Cosign, published to GHCR on `*-dc*` tags.

Upstream sync is performed weekly via the `upstream-sync` workflow; merges are reviewed manually.
