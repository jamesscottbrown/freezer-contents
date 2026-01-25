<script lang="ts">
    import { appState } from "$lib/stores";
    import { sortContainerNames } from "$lib/containerNames";

    interface ItemWithLocation {
        container: string;
        freezerName: string;
    }

    interface GroupedItem {
        name: string;
        containers: ItemWithLocation[];
    }

    const groupedByContents = $derived.by(() => {
        if (!$appState) return [];

        const itemMap = new Map<string, ItemWithLocation[]>();

        for (const freezer of $appState.Freezers) {
            for (const item of freezer.Contents) {
                const existing = itemMap.get(item.Name) || [];
                for (const container of item.Containers) {
                    existing.push({
                        container,
                        freezerName: freezer.Name
                    });
                }
                itemMap.set(item.Name, existing);
            }
        }

        const result: GroupedItem[] = [];
        for (const [name, containers] of itemMap) {
            const sortedContainers = containers.sort((a, b) => {
                const sorted = sortContainerNames([a.container, b.container]);
                return sorted[0] === a.container ? -1 : 1;
            });
            result.push({ name, containers: sortedContainers });
        }

        return result.sort((a, b) => a.name.localeCompare(b.name));
    });
</script>

<div class="flex flex-col gap-4">
    {#each groupedByContents as group}
        <div class="flex flex-col gap-1">
            <h2 class="font-bold">{group.name}</h2>
            <ul class="list-disc pl-6">
                {#each group.containers as item}
                    <li>{item.container} <span class="text-gray-500">({item.freezerName})</span></li>
                {/each}
            </ul>
        </div>
    {/each}
</div>
