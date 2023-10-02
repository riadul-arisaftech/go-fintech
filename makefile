create_sql:
	dbml2sql --postgres -o doc/db_schema.sql doc/db.dbml

create_migration:
	# create new migration
	migrate create -ext sql -dir db/migrations -seq ${name}

docker_up:
	# postgres up - create postgres server
	docker-compose up -d

docker_down:
	# postgres down - delete postgres server
	docker-compose down

db_up:
	docker exec -it fintech_postgres createdb --username=root --owner=root fintech_db
	docker exec -it fintech_postgres_live createdb --username=root --owner=root fintech_db

db_down:
	docker exec -it fintech_postgres dropdb --username=root fintech_db
	docker exec -it fintech_postgres_live dropdb --username=root fintech_db

m_up:
	# run migrate up
	migrate -path db/migrations -database "postgres://root:2654@localhost:5432/fintech_db?sslmode=disable" up
	migrate -path db/migrations -database "postgres://root:2654@localhost:5433/fintech_db?sslmode=disable" up

m_down:
	# run migrate down
	migrate -path db/migrations -database "postgres://root:2654@localhost:5432/fintech_db?sslmode=disable" down
	migrate -path db/migrations -database "postgres://root:2654@localhost:5433/fintech_db?sslmode=disable" down

sqlc:
	sqlc generate

start:
	CompileDaemon -command="./fintech_backend"

test:
	go test -v -cover ./...

.PHONY: create_sql create_migration docker_up docker_down db_up db_down m_up m_down sqlc start test