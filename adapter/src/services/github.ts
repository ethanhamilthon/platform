import simpleGit, { SimpleGit } from 'simple-git';
import { z } from 'zod';
import { Octokit } from '@octokit/rest';
import { encodeResponse } from '../utils/format';

export class GitHubAdapter {
  private client: SimpleGit;

  constructor() {
    this.client = simpleGit();
  }

  // Fetches all repositories (public and private) for the authenticated user
  getAllRepositories = async (data: string) => {
    const info = GetAllRepositoriesRepositorySchema.parse(JSON.parse(data));
    const octokit = new Octokit({ auth: info.token });
    const response = await octokit.repos.listForAuthenticatedUser();
    return encodeResponse({ message: "success", data: response.data });
  };

  // Clones a repository from GitHub, returns success message
  cloneRepository = async (data: string) => {
    const info = CloneRepositorySchema.parse(JSON.parse(data));
    const cloneUrl = this.getCloneUrl(info.url, info.token);

    await this.client.clone(cloneUrl, info.destination);
    return encodeResponse({ message: "success" });
  };

  // Constructs the clone URL using the token for authentication
  private getCloneUrl(url: string, token: string) {
    const urlParts = new URL(url);
    urlParts.username = token;

    return urlParts.toString();
  }
}

export const GetAllRepositoriesRepositorySchema = z.object({
  token: z.string().min(1),
});

export const CloneRepositorySchema = z.object({
  url: z.string().url(),
  destination: z.string(),
  token: z.string().min(1),
});