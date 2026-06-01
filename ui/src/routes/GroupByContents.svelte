<script lang="ts">
    import { appState, type FreezerItem } from "$lib/stores";
    import Item from "./Item.svelte";

    interface ItemByFreezer {
        item: FreezerItem;
        freezerName: string;
    }

    interface GroupedByName {
        name: string;
        itemsByFreezer: ItemByFreezer[];
    }

    const groupedByContents = $derived.by(() => {
        if (!$appState) return [];

        // Group by item name, collecting all items (including those with different dates)
        const nameMap = new Map<string, ItemByFreezer[]>();

        for (const freezer of $appState.Freezers) {
            for (const item of freezer.Contents) {
                if (!nameMap.has(item.Name)) {
                    nameMap.set(item.Name, []);
                }
                nameMap.get(item.Name)!.push({ item, freezerName: freezer.Name });
            }
        }

        const result: GroupedByName[] = [];
        for (const [name, itemsByFreezer] of nameMap) {
            // Sort by freezer name, then by date for consistent ordering
            itemsByFreezer.sort((a, b) => {
                const freezerCmp = a.freezerName.localeCompare(b.freezerName);
                if (freezerCmp !== 0) return freezerCmp;
                return a.item.Date.localeCompare(b.item.Date);
            });
            result.push({ name, itemsByFreezer });
        }

        return result.sort((a, b) => a.name.localeCompare(b.name));
    });
</script>

<div class="flex flex-col gap-4">
    {#each groupedByContents as group}
        <div class="flex flex-col gap-1">
            <h2 class="font-bold">{group.name}</h2>
            {#each group.itemsByFreezer as entry}
                <div class="pl-4 flex items-center gap-2">
                    <span class="text-gray-500">({entry.freezerName})</span>
                    <Item item={entry.item} freezerName={entry.freezerName} />
                </div>
            {/each}
        </div>
    {/each}
</div>
