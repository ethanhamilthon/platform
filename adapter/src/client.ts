import NatsMessage from "./message";

// only for clearfy how nats server works and test the methods, dont use in prod
async function main() {
  const nc = new NatsMessage();

  await nc.connect("nats://localhost:4222");

  const body = await nc.request("balancer:launch:http", "");

  console.log(body);
}

main();
