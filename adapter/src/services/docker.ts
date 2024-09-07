import { z } from "zod";
import DockerSDK from "dockerode";
import { encodeResponse } from '../utils/format';
import NatsMessage from '../message';
import { getEnv } from '../utils/config';

// DockerAdapter provides api to docker sdk
export class DockerAdapter {
  private client: DockerSDK;
  private nc: NatsMessage;

  constructor(nc: NatsMessage) {
    this.nc = nc;
    this.client = new DockerSDK();
  }

  // creates new container with existing image, be sure that image is pulled. It returns id of container [sync]
  createContainer = async (data: string) => {
    const info = CreateContainerSchema.parse(JSON.parse(data));
    const container = await this.client.createContainer
    ({
      Image: info.image_name,
      name: info.container_name,
    });
    return container.id;
  };
  
  // should be async? because of long start up time
  // Runs container with id [sync]
  runContainer = async (data: string) => {
    const info = BaseContainerSchema.parse(JSON.parse(data));
    const container = this.client.getContainer(info.id);
    await container.start();
    return encodeResponse({ message: "success" });
  };

  // should be async? because of long start up time(10 secs)
  // Stops container with id [sync]
  stopContainer = async (data: string) => {
    const info = BaseContainerSchema.parse(JSON.parse(data));
    const container = this.client.getContainer(info.id);
    await container.stop();
    return encodeResponse({ message: "success" });
  };

  removeContainer = async (data: string) => {
    const info = BaseContainerSchema.parse(JSON.parse(data));
    const container = this.client.getContainer(info.id);
    await container.remove();
    return encodeResponse({ message: "success" });
  };

  pullImage = async (data: string) => {
    const info = BaseImageSchema.parse(JSON.parse(data));
    const stream = await this.client.pull(info.image_name);

    stream.on('data', (data) => {
      this.nc.publish(getEnv("PULL_IMAGE_PROGRESS_TOPIC", "adapter:docker:pull-image-progress"), data.toString());
    });

    stream.on('end', () => {
      this.nc.publish(getEnv("PULL_IMAGE_END_TOPIC", "adapter:docker:pull-image-end"), "success");
    })
  };

  buildImage = async (data: string) => {
    const info = BuildImageSchema.parse(JSON.parse(data));
    this.client.buildImage({
      context: `${info.dirname}${info.path}`,
      src: ['Dockerfile']
    }, info.options, (err, stream) => {
      if (err) {
        console.log("ERR");
        console.log(err);
        throw Error(err);
      }

      if (!stream) {
        throw Error("Stream was not provided");
      }
      
      stream.on('data', (data) => {
        this.nc.publish(getEnv("BUILD_IMAGE_PROGRESS_TOPIC", "adapter:docker:build-image-progress"), data.toString());
      });

      stream.on('end', () => {
        this.nc.publish(getEnv("BUILD_IMAGE_END_TOPIC", "adapter:docker:build-image-end"), "success");
      });
    });
  }

  removeImage = async (data: string) => {
    const info = BaseImageSchema.parse(JSON.parse(data));
    await this.client.getImage(info.image_name).remove();
    return encodeResponse({ message: "success" });
  };

  listImages = async () => {
    const images = await this.client.listImages();
    return encodeResponse({ message: "success", data: images });
  };

  listContainers = async () => {
    const containers = await this.client.listContainers();
    return encodeResponse({ message: "success", data: containers });
  };
}

export const BaseContainerSchema = z.object({
  id: z.string(),
});

export const BaseImageSchema = z.object({
  image_name: z.string(),
})

export const BuildImageOptionsSchema = z.object({
  t: z.string().optional(),
});

export const BuildImageSchema = z.object({
  path: z.string(),
  dirname: z.string(),
  options: BuildImageOptionsSchema,
});

export const CreateContainerSchema = z.object({
  image_name: z.string(),
  container_name: z.string(),
});
