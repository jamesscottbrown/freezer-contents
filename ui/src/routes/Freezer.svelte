<script>
    import Item from "./Item.svelte";
    import AddItemModal from "./AddItemModal.svelte";

    export let freezer;
    export let addModalIsOpen = false;

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

    <button class="px-2 border border-green-500 rounded" on:click={() => addModalIsOpen = true}>Add Item</button>

    <AddItemModal bind:isOpen={addModalIsOpen} freezerName={freezer.Name} />
</div>