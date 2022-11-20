# Sportech
This is the monorepo that powers Sportech, the application that powers all things sports data.

## Setup
To get started with all the tools needed for development on this project run
```shell
./setup.sh
```

## Migrations

We use [golang-migrate](https://github.com/golang-migrate/migrate) to manage our database migrations.

You can invoke this via a helper script, [./migrate](./migrate), e.g.:

Usage:

To create a new migration called `<name>`:

```shell
./migrate create -ext sql -dir $(pwd)/migrations -seq <name>
```

You will then need to write the forward and reverse migrations yourself. It's best to follow the
[best practices for writing migrations](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md).

The following commands are optional, as the service should handle running migrations for you.

To run migrations:

```bash
./migrate up
```

To undo migrations:

```bash
./migrate down -all
```

### Failed migrations
If a migration fails, you will end up with a "dirty" database, and golang-migrate will refuse to migrate anymore. To fix this, you need to force
a version, by running:

```bash
./migrate force <version>
```
To determine the correct version, you can first find what version you are currently on with `./migrate version`, and then the version you want is very
likely the one before this (if you need to confirm, you can connect to the db to find what state it's in).
