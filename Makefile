up:
	docker-compose up -d

updb:
	docker-compose up -d db

bringtofg:
	docker-compose logs -f

createdb:
	docker-compose exec -u postgres db psql -h localhost -p 5432 -U postgres -d dev -f /schema.sql

down:
	docker-compose down