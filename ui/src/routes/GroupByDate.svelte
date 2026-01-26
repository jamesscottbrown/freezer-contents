<script lang="ts">
    import { appState } from "$lib/stores";
    import { sortContainerNames } from "$lib/containerNames";

    interface ItemWithLocation {
        name: string;
        date: string;
        containers: string[];
        freezerName: string;
    }

    const sortedByDate = $derived.by(() => {
        if (!$appState) return [];

        const items: ItemWithLocation[] = [];

        for (const freezer of $appState.Freezers) {
            for (const item of freezer.Contents) {
                items.push({
                    name: item.Name,
                    date: item.Date,
                    containers: sortContainerNames(item.Containers),
                    freezerName: freezer.Name
                });
            }
        }

        // Sort by date ascending (oldest first)
        return items.sort((a, b) => a.date.localeCompare(b.date));
    });
</script>

<div class="flex flex-col gap-4">
    {#each sortedByDate as item}
        <div class="flex flex-col gap-1">
            <h2 class="font-bold">{item.date}: {item.name}</h2>
            <ul class="list-disc pl-6">
                {#each item.containers as container}
                    <li>{container} <span class="text-gray-500">({item.freezerName})</span></li>
                {/each}
            </ul>
        </div>
    {/each}
</div>
