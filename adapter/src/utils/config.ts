export function getEnv(env: string, defaultValue: string) {
  return process.env[env] ?? defaultValue;
}