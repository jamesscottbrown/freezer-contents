/**
 * Parses a string of comma-separated numbers and ranges into an array of numbers.
 * Supports individual numbers and ranges expressed with hyphens.
 * Example: "1,3-5,7" -> [1, 3, 4, 5, 7]
 */
export function parseNumberInput(input: string): number[] {
  const result: number[] = [];
  const parts = input
    .split(",")
    .map((p) => p.trim())
    .filter((p) => p !== "");

  for (const part of parts) {
    if (part.includes("-")) {
      const [startStr, endStr] = part.split("-").map((s) => s.trim());
      const start = parseInt(startStr, 10);
      const end = parseInt(endStr, 10);

      if (!isNaN(start) && !isNaN(end) && start <= end) {
        for (let i = start; i <= end; i++) {
          result.push(i);
        }
      }
    } else {
      const num = parseInt(part, 10);
      if (!isNaN(num)) {
        result.push(num);
      }
    }
  }

  return result;
}

/**
 * Builds an array of container names from a number input string.
 * @param input - Comma-separated numbers with optional ranges (e.g., "1,3-5,7")
 * @param prefix - The prefix to prepend to each container number (e.g., "", "bowl_", "jar_")
 * @returns Array of container names with appropriate prefixes
 */
export function buildContainerNames(input: string, prefix: string): string[] {
  const numbers = parseNumberInput(input);
  return numbers.map((n) => `${prefix}${n}`);
}

/**
 * Parses a container name into its type and number parts.
 * Container names can be separated by hyphen or underscore.
 * Examples: "bowl_5" -> { type: "bowl", number: 5 }
 *           "u-3" -> { type: "u", number: 3 }
 *           "42" -> { type: "", number: 42 }
 */
export function parseContainerName(name: string): {
  type: string;
  number: number;
} {
  const match = name.match(/^([a-zA-Z]*)[-_]?(\d+)$/);
  if (match) {
    return {
      type: match[1] || "",
      number: parseInt(match[2], 10),
    };
  }
  // Fallback: treat entire name as type with no number
  return { type: name, number: 0 };
}

/**
 * Sorts container names intelligently:
 * - First alphabetically by container type (e.g., "bowl", "jar", "u")
 * - Then numerically by container number
 */
export function sortContainerNames(containers: string[]): string[] {
  return [...containers].sort((a, b) => {
    const parsedA = parseContainerName(a);
    const parsedB = parseContainerName(b);

    // First compare by type alphabetically
    const typeCompare = parsedA.type.localeCompare(parsedB.type);
    if (typeCompare !== 0) {
      return typeCompare;
    }

    // Then compare by number numerically
    return parsedA.number - parsedB.number;
  });
}
