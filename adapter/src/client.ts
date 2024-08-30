import NatsMessage from "./message";

// only for clearfy how nats server works and test the methods, dont use in prod
async function main() {
  const nc = new NatsMessage();

  await nc.connect("nats://localhost:4222");

  const id = await nc.request(
    "adapter:docker:create-container",
    JSON.stringify({
      image_name: "redis",
      container_name: "myredis",
    }),
  );

  nc.publish(
    "adapter:docker:run-container",
    JSON.stringify({
      id: id,
    }),
  );
}

main();
