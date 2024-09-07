import { mkdirSync, unlink, unlinkSync, writeFileSync } from 'fs';

import NatsMessage from './message';
import { getEnv } from './utils/config';

// only for clearfy how nats server works and test the methods, dont use in prod
async function main() {
  const nc = new NatsMessage();

  await nc.connect(getEnv("NATS_URL", "http://0.0.0.0:4222"));

  // Docker Service Tests
  // Should pull Image
  nc.publish(getEnv("PULL_IMAGE_TOPIC", "adapter:docker:pull-image"), JSON.stringify({
    image_name: 'redis',
  }));
  console.log("[LOG][TEST][IMAGE](pulled image) ")

  // Creating container
  const containerId = await nc.request(
    getEnv("CREATE_CONTAINER_TOPIC", "adapter:docker:create-container"),
    JSON.stringify({
      image_name: "redis",
      container_name: "myredis",
    }),
  );
  console.log("[LOG][TEST][CONTAINER](created container): " + containerId)

  // Should run container
  const runContainerMsg = await nc.request(
    getEnv("RUN_CONTAINER_TOPIC", "adapter:docker:run-container"),
    JSON.stringify({
      id: containerId,
    }),
  );
  console.log("[LOG][TEST][CONTAINER](run container): " + runContainerMsg)

  // Pulling Image
  // Should build image basing on Dockerfile
  const testDockerfilePath = '/tests/dockerfiles/redis';
  const testDockerfileName = 'Dockerfile';

  generateTestDockerfile(testDockerfilePath, testDockerfileName);
  await nc.publish(
    getEnv("BUILD_IMAGE_TOPIC", "adapter:docker:build-image"),
    JSON.stringify({
      dirname: __dirname,
      path: testDockerfilePath,
      options: { t: 'test-redis-dockerfile' },
    }),
  );

  // unlinkSync(`${testDockerfilePath}/${testDockerfileName}`);
  console.log("[LOG][TEST][CONTAINER](built image)")

  // Cleanup tests
  // Should stop container
  const stopContainerId = await nc.request(
    getEnv("STOP_CONTAINER_TOPIC", "adapter:docker:stop-container"),
    JSON.stringify({
      id: containerId,
    }),
  );
  console.log("[LOG][TEST][CONTAINER](stop container): " + stopContainerId)

  // Should remove container
  const removeContainerMsg = await nc.request(
    getEnv("REMOVE_CONTAINER_TOPIC", "adapter:docker:remove-container"),
    JSON.stringify({
      id: containerId,
    }),
  );
  console.log("[LOG][TEST][CONTAINER](remove container): " + removeContainerMsg)
  
  // Should remove image
  const removeImageMsg = await nc.request(
    getEnv("REMOVE_IMAGE_TOPIC", "adapter:docker:remove-image"),
    JSON.stringify({
      image_name: "redis",
    }),
  );
  console.log("[LOG][TEST][IMAGE](remove image): " + removeImageMsg)

  // Should list all images
  const imagesList = await nc.request(
    getEnv("LIST_IMAGES_TOPIC", "adapter:docker:list-images"),
    null,
  );
  console.log("[LOG][TEST][IMAGE](list of all images): " + imagesList)

  // Should list all containers
  const containersList = await nc.request(
    getEnv("LIST_CONTAINERS_TOPIC", "adapter:docker:list-containers"),
    null,
  );
  console.log("[LOG][TEST][CONTAINER](list of all containers): " + containersList)
}

main();

function generateTestDockerfile(path: string, filename: string = 'Dockerfile') {
  mkdirSync(`${__dirname}/${path}`, { recursive: true })
  writeFileSync(`${__dirname}/${path}/${filename}`, 'FROM redis', { flag: 'w+' })
}
