Run unit test with docker

Option #1:

``
docker run --cpus=1 --rm -v $(pwd):/app --workdir /app golang:1.21 go test ./pubsub/gochannel -run ^TestSubscriber_race_condition_when_closing$ -v
``

Option #2:

``
docker build -t gochanneltest .
``

``
docker run --rm gochanneltest:latest go test ./pubsub/gochannel -run ^TestSubscriber_race_condition_when_closing$ -v
``