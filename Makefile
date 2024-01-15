MAKEFLAGS += --silent

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DATABASE_URL=postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=${DATABASE_SSL_MODE}

# Download dependency
.PHONY: mod
mod:
	go mod tidy -compat=1.17
	go mod vendor

# Database Management
# Create database
.PHONY: createdb
createdb:
	createdb "${DATABASE_NAME}"

# Drop database
.PHONY: dropdb
dropdb:
	dropdb "${DATABASE_NAME}"

# Migrate database
.PHONY: migratedb
migratedb:
	migrate --path=db/migrations/ \
			--database ${DATABASE_URL} up

# Rollback database
.PHONY: rollbackdb
rollbackdb:
	echo "y" | migrate --path=db/migrations/ \
			--database ${DATABASE_URL} down

# Create database migration
migration:
	$(eval timestamp := $(shell date +%s))
	touch db/migrations/$(timestamp)_${name}.up.sql
	touch db/migrations/$(timestamp)_${name}.down.sql
