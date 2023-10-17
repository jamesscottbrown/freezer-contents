<script lang="ts">
    import {Dialog, DialogDescription, DialogOverlay, DialogTitle} from "@rgossiaux/svelte-headlessui";
    import {state} from "$lib/stores.ts";

    export let item;
    export let freezerName;

    let selectedContainer;

    const removeContainer = () => {
        const sure = confirm(`Are you sure you want to remove ${selectedContainer}?`);
        if (!sure) {
            return;
        }
        // {Container: selectedContainer}
        const url = "/remove";

        fetch(url, {
            method: "POST",
            body: JSON.stringify({Container: selectedContainer}),
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(res => res.json())
            .then(d => $state = d);

        selectedContainer = undefined;
    };

    const moveContainer = (newFreezer: string) => {
        const url = "/move";

        fetch(url, {
            method: "POST",
            body: JSON.stringify({Container: selectedContainer, newFreezer}),
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(res => res.json())
            .then(d => $state = d);

        selectedContainer = undefined;
    };

</script>


<div class="flex gap-2">
    <span>{item.Name} ({item.Date})</span>
    {#each item.Containers.sort() as container}
        <button on:click={() => selectedContainer = container}
                class="px-2 border border-grey-500 rounded">{container}</button>
    {/each}
</div>


<Dialog open={!!selectedContainer} on:close={() => (selectedContainer = undefined)}
        class="relative z-10 overflow-y-auto">
    <DialogOverlay class="fixed inset-0 bg-black bg-opacity-40"/>

    <!--Full-screen container to center the panel -->
    <div class="fixed inset-0 flex items-center justify-center p-4 pointer-events-none">
        <div
                class="inline-block w-full max-w-md my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl space-y-2 pb-4 pointer-events-auto"
        >
            <div class="bg-core-grey-600 text-white p-2 relative">
                <DialogTitle>Edit {item.Name}  ({item.Date})</DialogTitle>
                <button
                        on:click={() => (selectedContainer = undefined)}
                        class="bg-core-grey-500 absolute top-2 right-2 hover:bg-core-grey-800"
                >
                    <span class="sr-only">Close</span>
                    <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke-width="1.5"
                            stroke="currentColor"
                            class="w-6 h-6"
                    >
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                </button>
            </div>

            <DialogDescription class="px-2">Selected Container: {selectedContainer}</DialogDescription>

            <div class="flex flex-col pl-2">

                <div>
                    <button class="px-2 border border-red-500 rounded" on:click={removeContainer}>Remove</button>

                    {#each $state.Freezers as freezer}
                        {#if freezer.Name !== freezerName}
                            <button class="px-2 border border-green-500 rounded" on:click={() => moveContainer(freezer.Name)}>Move to {freezer.Name}</button>
                        {/if}
                    {/each}
                </div>

            </div>

        </div>
    </div>
</Dialog>
