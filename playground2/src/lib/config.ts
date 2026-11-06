import type { GenerateConfig } from "./types";

export const DEFAULT_CONFIG: GenerateConfig = { count: 1, scope: 3, actors: 3 };
export const MAX_CONFIG: GenerateConfig = { count: 50, scope: 6, actors: 10 };

export function clampConfig(c: GenerateConfig): GenerateConfig {
  const f = (v: number, d: number, max: number) =>
    Number.isFinite(v) ? Math.min(max, Math.max(1, Math.round(v))) : d;
  return {
    count: f(c.count, DEFAULT_CONFIG.count, MAX_CONFIG.count),
    scope: f(c.scope, DEFAULT_CONFIG.scope, MAX_CONFIG.scope),
    actors: f(c.actors, DEFAULT_CONFIG.actors, MAX_CONFIG.actors),
  };
}
