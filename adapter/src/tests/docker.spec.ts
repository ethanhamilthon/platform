// testing docker adapter via local docker setup

// pull nats docker image:
// #1 docker pull nats

// run a server with the ports exposed on the 'nats' docker network:
// docker run --name nats --network nats --rm -p 4222:4222 -p 8222:8222 nats --http_port 8222 --cluster_name NATS --cluster nats://0.0.0.0:6222

// run two files
// index.ts first(addition of subs to topics)
// client.ts test of those subscriptions