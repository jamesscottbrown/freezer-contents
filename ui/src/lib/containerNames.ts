/**
 * Parses a string of comma-separated numbers and ranges into an array of numbers.
 * Supports individual numbers and ranges expressed with hyphens.
 * Example: "1,3-5,7" -> [1, 3, 4, 5, 7]
 */
export function parseNumberInput(input: string): number[] {
    const result: number[] = [];
    const parts = input.split(",").map(p => p.trim()).filter(p => p !== "");

    for (const part of parts) {
        if (part.includes("-")) {
            const [startStr, endStr] = part.split("-").map(s => s.trim());
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

export type ContainerType = "containers" | "bowls" | "jars";

/**
 * Builds an array of container names from a number input string.
 * @param input - Comma-separated numbers with optional ranges (e.g., "1,3-5,7")
 * @param containerType - The type of container ("containers", "bowls", or "jars")
 * @returns Array of container names with appropriate prefixes
 */
export function buildContainerNames(input: string, containerType: ContainerType): string[] {
    const numbers = parseNumberInput(input);

    switch (containerType) {
        case "containers":
            return numbers.map(n => String(n));
        case "bowls":
            return numbers.map(n => `bowl_${n}`);
        case "jars":
            return numbers.map(n => `jar_${n}`);
    }
}
