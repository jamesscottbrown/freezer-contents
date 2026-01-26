import { describe, it, expect } from "vitest";
import {
  parseNumberInput,
  buildContainerNames,
  parseContainerName,
  sortContainerNames,
} from "./containerNames";

describe("parseNumberInput", () => {
  it("should parse a single number", () => {
    expect(parseNumberInput("1")).toEqual([1]);
    expect(parseNumberInput("42")).toEqual([42]);
  });

  it("should parse comma-separated numbers", () => {
    expect(parseNumberInput("1,2,3")).toEqual([1, 2, 3]);
    expect(parseNumberInput("1, 2, 3")).toEqual([1, 2, 3]);
  });

  it("should parse a range", () => {
    expect(parseNumberInput("1-5")).toEqual([1, 2, 3, 4, 5]);
    expect(parseNumberInput("3-5")).toEqual([3, 4, 5]);
  });

  it("should parse mixed numbers and ranges", () => {
    expect(parseNumberInput("1,3-5,7")).toEqual([1, 3, 4, 5, 7]);
    expect(parseNumberInput("1, 3-5, 7")).toEqual([1, 3, 4, 5, 7]);
  });

  it("should handle multiple ranges", () => {
    expect(parseNumberInput("1-3,5-7")).toEqual([1, 2, 3, 5, 6, 7]);
  });

  it("should handle whitespace", () => {
    expect(parseNumberInput(" 1 , 2 , 3 ")).toEqual([1, 2, 3]);
    expect(parseNumberInput("1 - 3")).toEqual([1, 2, 3]);
  });

  it("should handle a single-element range", () => {
    expect(parseNumberInput("5-5")).toEqual([5]);
  });

  it("should return empty array for empty input", () => {
    expect(parseNumberInput("")).toEqual([]);
    expect(parseNumberInput("   ")).toEqual([]);
  });

  it("should ignore invalid entries", () => {
    expect(parseNumberInput("1,abc,3")).toEqual([1, 3]);
    expect(parseNumberInput("5-3")).toEqual([]); // invalid range (start > end) is skipped
  });

  it("should handle complex mixed input", () => {
    expect(parseNumberInput("1,2,5-8,10,12-14")).toEqual([
      1, 2, 5, 6, 7, 8, 10, 12, 13, 14,
    ]);
  });
});

describe("buildContainerNames", () => {
  describe("with empty prefix", () => {
    it("should return string numbers", () => {
      expect(buildContainerNames("1,2,3", "")).toEqual(["1", "2", "3"]);
    });

    it("should handle ranges", () => {
      expect(buildContainerNames("1,3-5,7", "")).toEqual([
        "1",
        "3",
        "4",
        "5",
        "7",
      ]);
    });
  });

  describe("with bowl_ prefix", () => {
    it("should prefix with bowl_", () => {
      expect(buildContainerNames("1,2,3", "bowl_")).toEqual([
        "bowl_1",
        "bowl_2",
        "bowl_3",
      ]);
    });

    it("should handle ranges", () => {
      expect(buildContainerNames("1,3-5", "bowl_")).toEqual([
        "bowl_1",
        "bowl_3",
        "bowl_4",
        "bowl_5",
      ]);
    });
  });

  describe("with jar_ prefix", () => {
    it("should prefix with jar_", () => {
      expect(buildContainerNames("1,2,3", "jar_")).toEqual([
        "jar_1",
        "jar_2",
        "jar_3",
      ]);
    });

    it("should handle ranges", () => {
      expect(buildContainerNames("1,3-5", "jar_")).toEqual([
        "jar_1",
        "jar_3",
        "jar_4",
        "jar_5",
      ]);
    });
  });

  it("should return empty array for empty input", () => {
    expect(buildContainerNames("", "")).toEqual([]);
    expect(buildContainerNames("", "bowl_")).toEqual([]);
    expect(buildContainerNames("", "jar_")).toEqual([]);
  });
});

describe("parseContainerName", () => {
  it("should parse underscore-separated names", () => {
    expect(parseContainerName("bowl_5")).toEqual({ type: "bowl", number: 5 });
    expect(parseContainerName("jar_10")).toEqual({ type: "jar", number: 10 });
  });

  it("should parse hyphen-separated names", () => {
    expect(parseContainerName("u-3")).toEqual({ type: "u", number: 3 });
    expect(parseContainerName("box-42")).toEqual({ type: "box", number: 42 });
  });

  it("should parse plain numbers", () => {
    expect(parseContainerName("42")).toEqual({ type: "", number: 42 });
    expect(parseContainerName("1")).toEqual({ type: "", number: 1 });
  });

  it("should handle names without separator", () => {
    expect(parseContainerName("bowl5")).toEqual({ type: "bowl", number: 5 });
  });

  it("should fallback for unrecognized formats", () => {
    expect(parseContainerName("abc")).toEqual({ type: "abc", number: 0 });
  });
});

describe("sortContainerNames", () => {
  it("should sort numerically within same type", () => {
    expect(sortContainerNames(["bowl_10", "bowl_2", "bowl_1"])).toEqual([
      "bowl_1",
      "bowl_2",
      "bowl_10",
    ]);
  });

  it("should sort alphabetically by type first", () => {
    expect(sortContainerNames(["jar_1", "bowl_1", "u-1"])).toEqual([
      "bowl_1",
      "jar_1",
      "u-1",
    ]);
  });

  it("should handle mixed types and numbers", () => {
    expect(sortContainerNames(["jar_5", "bowl_3", "jar_2", "bowl_10"])).toEqual(
      ["bowl_3", "bowl_10", "jar_2", "jar_5"],
    );
  });

  it("should handle plain numbers (no type prefix)", () => {
    expect(sortContainerNames(["10", "2", "1"])).toEqual(["1", "2", "10"]);
  });

  it("should sort plain numbers before prefixed containers", () => {
    expect(sortContainerNames(["bowl_1", "5", "2"])).toEqual([
      "2",
      "5",
      "bowl_1",
    ]);
  });

  it("should handle empty array", () => {
    expect(sortContainerNames([])).toEqual([]);
  });

  it("should not mutate original array", () => {
    const original = ["bowl_10", "bowl_2", "bowl_1"];
    sortContainerNames(original);
    expect(original).toEqual(["bowl_10", "bowl_2", "bowl_1"]);
  });

  it("should handle u- prefix containers", () => {
    expect(sortContainerNames(["u-10", "u-2", "u-1"])).toEqual([
      "u-1",
      "u-2",
      "u-10",
    ]);
  });
});
