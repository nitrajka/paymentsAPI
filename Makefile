up:
	docker-compose up -d

updb:
	docker-compose up db -d

createdb:
	docker-compose exec -u postgres db psql -h localhost -p 5432 -U postgres -d dev -f /createDB.sql

down:
	docker-compose down