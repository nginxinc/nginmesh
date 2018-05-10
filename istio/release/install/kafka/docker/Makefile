TAG = 4.1.0
NAME = nginmesh/kafka-connect


all:    build push

build:
	docker build -t $(NAME):$(TAG) --build-arg CONFLUENT_VERSION=$(TAG) .

push:
	docker push $(NAME):$(TAG)
