.PHONY: install-tools
install-tools:
	@./scripts/install-tools.sh

.PHONY: run-catalogs-write-service
run-catalogs-write-service:
	@./scripts/run.sh  catalog_write_service

.PHONY: run-catalog-read-service
run-catalog-read-service:
	@./scripts/run.sh  catalog_read_service

.PHONY: run-order-service
run-order-service:
	@./scripts/run.sh  order_service

.PHONY: build
build:
	@./scripts/build.sh  pkg
	@./scripts/build.sh  catalog_write_service
	@./scripts/build.sh  catalog_read_service
	@./scripts/build.sh  order_service

.PHONY: install-dependencies
install-dependencies:
	@./scripts/install-dependencies.sh  pkg
	@./scripts/install-dependencies.sh  catalog_write_service
	@./scripts/install-dependencies.sh  catalog_read_service
	@./scripts/install-dependencies.sh  order_service

.PHONY: docker-compose-infra-up
docker-compose-infra-up:
	@docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml up --build -d

docker-compose-infra-down:
	@docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml down

.PHONY: openapi
openapi:
	@./scripts/openapi.sh catalog_write_service
	@./scripts/openapi.sh catalog_read_service
	@./scripts/openapi.sh order_service

# https://stackoverflow.com/questions/13616033/install-protocol-buffers-on-windows
.PHONY: proto
proto:
	@./scripts/proto.sh catalog_write_service
	@./scripts/proto.sh order_service

.PHONY: unit-test
unit-test:
	@./scripts/test.sh catalog_read_service unit
	@./scripts/test.sh catalog_write_service unit
	@./scripts/test.sh  order_service unit

.PHONY: integration-test
integration-test:
	@./scripts/test.sh catalog_read_service integration
	@./scripts/test.sh catalog_write_service integration
	@./scripts/test.sh  order_service integration

.PHONY: e2e-test
e2e-test:
	@./scripts/test.sh catalog_read_service e2e
	@./scripts/test.sh catalog_write_service e2e
	@./scripts/test.sh  order_service e2e

#.PHONY: load-test
#load-test:
#	@./scripts/test.sh catalogs_write load-test
#	@./scripts/test.sh catalogs_read load-test
#	@./scripts/test.sh  orders load-test

.PHONY: format
format:
	@./scripts/format.sh catalog_write_service
	@./scripts/format.sh catalog_read_service
	@./scripts/format.sh order_service
	@./scripts/format.sh pkg

.PHONY: lint
lint:
	@./scripts/lint.sh catalog_write_service
	@./scripts/lint.sh catalog_read_service
	@./scripts/lint.sh order_service
	@./scripts/lint.sh pkg

# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/database/postgres/TUTORIAL.md
# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/GETTING_STARTED.md
# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/MIGRATIONS.md
# https://github.com/golang-migrate/migrate/tree/856ea12df9d230b0145e23d951b7dbd6b86621cb/cmd/migrate#usage
.PHONY: go-migrate
go-migrate:
	@./scripts/go-migrate.sh -p ./internal/services/catalog_write_service/db/migrations/go-migrate -c create -n create_product_table
	@./scripts/go-migrate.sh -p ./internal/services/catalog_write_service/db/migrations/go-migrate -c up -o postgres://postgres:postgres@localhost:5432/catalogs_service?sslmode=disable
	@./scripts/go-migrate.sh -p ./internal/services/catalog_write_service/db/migrations/go-migrate -c down -o postgres://postgres:postgres@localhost:5432/catalogs_service?sslmode=disable

# https://github.com/pressly/goose#usage
.PHONY: goose-migrate
goose-migrate:
	@./scripts/goose-migrate.sh -p ./internal/services/catalog_write_service/db/migrations/goose-migrate -c create -n create_product_table
	@./scripts/goose-migrate.sh -p ./internal/services/catalog_write_service/db/migrations/goose-migrate -c up -o "user=postgres password=postgres dbname=catalogs_service sslmode=disable"
	@./scripts/goose-migrate.sh -p ./internal/services/catalog_write_service/db/migrations/goose-migrate -c down -o "user=postgres password=postgres dbname=catalogs_service sslmode=disable"

# https://atlasgo.io/guides/orms/gorm
.PHONY: atlas
atlas:
	@./scripts/atlas-migrate.sh -c gorm-sync -p "./internal/services/catalog_write_service"
	@./scripts/atlas-migrate.sh -c apply -p "./internal/services/catalog_write_service" -o "postgres://postgres:postgres@localhost:5432/catalogs_service?sslmode=disable"

#.PHONY: c4
#c4:
#	cd tools/c4 && go mod tidy && sh generate.sh

# https://medium.com/yemeksepeti-teknoloji/mocking-an-interface-using-mockery-in-go-afbcb83cc773
# https://vektra.github.io/mockery/latest/running/
# https://amitshekhar.me/blog/test-with-testify-and-mockery-in-go
.PHONY: pkg-mocks
pkg-mocks:
	cd internal/pkg/messaging && mockery --output mocks --all
	cd internal/pkg/es && mockery --output mocks --all
	cd internal/pkg/core && mockery --output mocks --all

.PHONY: services-mocks
services-mocks:
	cd internal/services/catalog_write_service && mockery --output mocks --all
	cd internal/services/catalog_read_service && mockery --output mocks --all
	cd internal/services/order_service && mockery --output mocks --all
