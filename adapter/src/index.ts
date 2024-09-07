import NatsMessage from "./message";
import { DockerAdapter } from "./services/docker";
// import { GitHubAdapter } from './services/github';
import { getEnv } from './utils/config';

async function main() {
  const nc = new NatsMessage();
  // const github = new GitHubAdapter();
  const docker = new DockerAdapter(nc);
  await nc.connect(getEnv("NATS_URL", "http://0.0.0.0:4222"));

  nc.addSub(
    getEnv("CREATE_CONTAINER_TOPIC", "adapter:docker:create-container"),
    nc.errorHandler(docker.createContainer),
  )
  nc.addSub(
    getEnv("RUN_CONTAINER_TOPIC", "adapter:docker:run-container"),
    nc.errorHandler(docker.runContainer),
  );
  nc.addSub(
    getEnv("STOP_CONTAINER_TOPIC", "adapter:docker:stop-container"),
    nc.errorHandler(docker.stopContainer),
  );
  nc.addSub(
    getEnv("REMOVE_CONTAINER_TOPIC", "adapter:docker:remove-container"),
    nc.errorHandler(docker.removeContainer),
  );
  nc.addSub(
    getEnv("PULL_IMAGE_TOPIC", "adapter:docker:pull-image"),
    nc.errorHandler(docker.pullImage),
  );
  nc.addSub(
    getEnv("BUILD_IMAGE_TOPIC", "adapter:docker:build-image"),
    nc.errorHandler(docker.buildImage),
  );
  nc.addSub(
    getEnv("REMOVE_IMAGE_TOPIC", "adapter:docker:remove-image"),
    nc.errorHandler(docker.removeImage),
  );
  nc.addSub(
    getEnv("LIST_IMAGES_TOPIC", "adapter:docker:list-images"),
    nc.errorHandler(docker.listImages),
  );
  nc.addSub(
    getEnv("LIST_CONTAINERS_TOPIC", "adapter:docker:list-containers"),
    nc.errorHandler(docker.listContainers),
  );

  // nc.addSub(
  //   getEnv("CLONE_REPOSITORY_TOPIC", "adapter:github:clone-repository"),
  //   nc.errorHandler(github.cloneRepository),
  // );
  // nc.addSub(
  //   getEnv("GET_ALL_REPOSITORIES_TOPIC", "adapter:github:get-all-repositories"),
  //   nc.errorHandler(github.getAllRepositories),
  // );
}

main();
