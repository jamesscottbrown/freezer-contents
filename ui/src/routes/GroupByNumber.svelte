<script lang="ts">
    import { appState } from "$lib/stores";
    import { parseContainerName } from "$lib/containerNames";

    interface ContainerInfo {
        name: string;
        number: number;
        itemName: string;
        freezerName: string;
    }

    interface ContainerTypeGroup {
        type: string;
        displayName: string;
        containers: ContainerInfo[];
        missingNumbers: number[];
    }

    function getDisplayName(type: string): string {
        if (type === "") return "Containers";
        if (type === "bowl") return "Bowls";
        if (type === "jar") return "Jars";
        return type.charAt(0).toUpperCase() + type.slice(1) + "s";
    }

    const groupedByNumber = $derived.by(() => {
        if (!$appState) return [];

        const typeMap = new Map<string, ContainerInfo[]>();

        for (const freezer of $appState.Freezers) {
            for (const item of freezer.Contents) {
                for (const container of item.Containers) {
                    const parsed = parseContainerName(container);
                    const existing = typeMap.get(parsed.type) || [];
                    existing.push({
                        name: container,
                        number: parsed.number,
                        itemName: item.Name,
                        freezerName: freezer.Name
                    });
                    typeMap.set(parsed.type, existing);
                }
            }
        }

        const result: ContainerTypeGroup[] = [];
        for (const [type, containers] of typeMap) {
            // Sort containers by number
            const sortedContainers = containers.sort((a, b) => a.number - b.number);

            // Find missing numbers
            const numbers = sortedContainers.map(c => c.number);
            const minNum = Math.min(...numbers);
            const maxNum = Math.max(...numbers);
            const missingNumbers: number[] = [];
            for (let i = minNum; i <= maxNum; i++) {
                if (!numbers.includes(i)) {
                    missingNumbers.push(i);
                }
            }

            result.push({
                type,
                displayName: getDisplayName(type),
                containers: sortedContainers,
                missingNumbers
            });
        }

        // Sort groups: plain containers first, then alphabetically
        return result.sort((a, b) => {
            if (a.type === "" && b.type !== "") return -1;
            if (a.type !== "" && b.type === "") return 1;
            return a.type.localeCompare(b.type);
        });
    });
</script>

<div class="flex flex-col gap-6">
    {#each groupedByNumber as group}
        <div class="flex flex-col gap-2">
            <h2 class="font-bold text-lg">{group.displayName}</h2>
            <ol class="list-decimal pl-6">
                {#each group.containers as container}
                    <li value={container.number}>
                        {container.name}: {container.itemName} <span class="text-gray-500">({container.freezerName})</span>
                    </li>
                {/each}
            </ol>
            {#if group.missingNumbers.length > 0}
                <p class="text-gray-500 text-sm pl-2">
                    Missing: {group.missingNumbers.join(", ")}
                </p>
            {/if}
        </div>
    {/each}
</div>
