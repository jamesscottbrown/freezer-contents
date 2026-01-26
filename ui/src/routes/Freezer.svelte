<script lang="ts">
    import Item from "./Item.svelte";
    import AddItemModal from "./AddItemModal.svelte";
    import type { Freezer } from "$lib/stores";

    let { freezer }: { freezer: Freezer } = $props();
    let addModalIsOpen = $state(false);
</script>

<div class="flex flex-col gap-2">
    <h2 class="font-bold">{freezer.Name}</h2>

    <div>
        <ul class="list-decimal">
            {#each freezer.Contents as item}
                <li>
                    <Item {item} freezerName={freezer.Name}/>
                </li>
            {/each}
        </ul>
    </div>

    <button class="px-2 border border-green-500 rounded" onclick={() => addModalIsOpen = true}>Add Item</button>

    <AddItemModal bind:isOpen={addModalIsOpen} freezerName={freezer.Name} />
</div>