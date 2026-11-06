import { describe, expect, it } from "vitest";
import { clampWorldIndex, worldLabel } from "./worlds";

describe("clampWorldIndex", () => {
  it("clamps into range", () => {
    expect(clampWorldIndex(5, 3)).toBe(2);
    expect(clampWorldIndex(-1, 3)).toBe(0);
    expect(clampWorldIndex(1, 3)).toBe(1);
  });
  it("is 0 when empty or invalid", () => {
    expect(clampWorldIndex(2, 0)).toBe(0);
    expect(clampWorldIndex(NaN, 3)).toBe(0);
  });
});

describe("worldLabel", () => {
  it("counts non-blank lines", () => {
    expect(worldLabel(0, "ephemeral", "a\n\nb\n")).toBe("World 1 (ephemeral, 2 relationships)");
    expect(worldLabel(1, "saved", "a")).toBe("World 2 (saved, 1 relationship)");
  });
  it("uses the name of a saved world file", () => {
    expect(worldLabel(2, "saved", "a\nb", "child-owner")).toBe("child-owner (saved, 2 relationships)");
  });
});
