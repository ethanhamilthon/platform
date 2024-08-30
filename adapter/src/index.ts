import NatsMessage from "./message";
import { DockerAdapter } from "./services/docker";

async function main() {
  const docker = new DockerAdapter();
  const nc = new NatsMessage();
  await nc.connect("nats://localhost:4222");

  nc.addSub(
    "adapter:docker:create-container",
    nc.errorHandler(docker.createContainer),
  );
  nc.addSub(
    "adapter:docker:run-container",
    nc.errorHandler(docker.runContainer),
  );
}

main();
