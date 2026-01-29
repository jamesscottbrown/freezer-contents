<script lang="ts">
    import { appState, type FreezerItem } from "$lib/stores";
    import Item from "./Item.svelte";

    interface ItemWithLocation {
        item: FreezerItem;
        freezerName: string;
    }

    const sortedByDate = $derived.by(() => {
        if (!$appState) return [];

        const items: ItemWithLocation[] = [];

        for (const freezer of $appState.Freezers) {
            for (const item of freezer.Contents) {
                items.push({
                    item,
                    freezerName: freezer.Name
                });
            }
        }

        // Sort by date ascending (oldest first)
        return items.sort((a, b) => a.item.Date.localeCompare(b.item.Date));
    });
</script>

<div class="flex flex-col gap-4">
    {#each sortedByDate as entry}
        <div class="flex flex-col gap-1">
            <span class="text-gray-500">({entry.freezerName})</span>
            <Item item={entry.item} freezerName={entry.freezerName} />
        </div>
    {/each}
</div>
