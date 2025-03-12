# running
run:
	docker compose up -d

run-build:
	docker compose up --build -d

stop:
	docker compose down

# linting
lint:
	golangci-lint run --enable-all

test:
	go test -v

# building
build:
	go build -o main main.go 

# clean
clean:
	rm ./main

# builds a docker image
docker-build:
	docker build -t ${IMAGE_TAG} -f ./docker/Dockerfile .

# builds a docker image and loads into minikube
minikube: docker-build
	minikube image rm ${IMAGE_TAG} 2> /dev/null
	minikube image load ${IMAGE_TAG}

