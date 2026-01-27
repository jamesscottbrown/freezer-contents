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

        // Group by item name, then by freezer
        const nameMap = new Map<string, Map<string, FreezerItem>>();

        for (const freezer of $appState.Freezers) {
            for (const item of freezer.Contents) {
                if (!nameMap.has(item.Name)) {
                    nameMap.set(item.Name, new Map());
                }
                const freezerMap = nameMap.get(item.Name)!;
                freezerMap.set(freezer.Name, item);
            }
        }

        const result: GroupedByName[] = [];
        for (const [name, freezerMap] of nameMap) {
            const itemsByFreezer: ItemByFreezer[] = [];
            for (const [freezerName, item] of freezerMap) {
                itemsByFreezer.push({ item, freezerName });
            }
            // Sort by freezer name for consistent ordering
            itemsByFreezer.sort((a, b) => a.freezerName.localeCompare(b.freezerName));
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
