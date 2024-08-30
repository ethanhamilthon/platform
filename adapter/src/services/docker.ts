import { z } from "zod";
import DockerSDK from "dockerode";

// DockerAdapter provides api to docker sdk
export class DockerAdapter {
  private client: DockerSDK;
  constructor() {
    this.client = new DockerSDK();
  }
  
  // creates new container with existing image, be sure that image is pulled. It returns id of container
  createContainer = async (data: string) => {
    const info = CreateContainerSchema.parse(JSON.parse(data));
    const container = await this.client.createContainer({
      Image: info.image_name,
      name: info.container_name,
    });
    return container.id;
  };
  
  // Runs container with id
  runContainer = async (data: string) => {
    const info = RunContainerSchema.parse(JSON.parse(data));
    const container = this.client.getContainer(info.id);
    container.start();
  };
}
export const RunContainerSchema = z.object({
  id: z.string(),
});
export const CreateContainerSchema = z.object({
  image_name: z.string(),
  container_name: z.string(),
});
