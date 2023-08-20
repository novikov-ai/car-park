APP_NAME = car-park
DB_NAME = car-park-db

.PHONY: build run clean compose-build compose-run compose-clean

build:
	docker build -t $(APP_NAME) .

run:
	docker run -d --name $(DB_NAME) \
		-e POSTGRES_PASSWORD=labuser POSTGRES_PASSWORD=labpassword POSTGRES_PASSWORD=labdb postgres
	docker run -d --name $(APP_NAME) -p 8080:8080 $(APP_NAME)

clean:
	docker stop $(APP_NAME) $(DB_NAME)
	docker rm $(APP_NAME) $(DB_NAME)
	docker network rm $(NETWORK_NAME)

compose-build:
	docker-compose build

compose-run:
	docker-compose up -d

compose-clean:
	docker-compose down