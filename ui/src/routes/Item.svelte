<script lang="ts">
    import { appState } from "$lib/stores";
    import { sortContainerNames } from "$lib/containerNames";

    let { item, freezerName }: { item: any; freezerName: string } = $props();

    let selectedContainer: string | undefined = $state(undefined);
    let dialogEl: HTMLDialogElement;

    const openDialog = (container: string) => {
        selectedContainer = container;
        dialogEl?.showModal();
    };

    const closeDialog = () => {
        dialogEl?.close();
        selectedContainer = undefined;
    };

    const removeContainer = () => {
        const sure = confirm(`Are you sure you want to remove ${selectedContainer}?`);
        if (!sure) {
            return;
        }
        const url = "/remove";

        fetch(url, {
            method: "POST",
            body: JSON.stringify({Container: selectedContainer}),
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(res => res.json())
            .then(d => appState.set(d));

        closeDialog();
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
            .then(d => appState.set(d));

        closeDialog();
    };
</script>


<div class="flex gap-2">
    <span>{item.Name} ({item.Date})</span>
    {#each sortContainerNames(item.Containers) as container}
        <button onclick={() => openDialog(container)}
                class="px-2 border border-grey-500 rounded"
                aria-label="Edit container {container} for {item.Name}">{container}</button>
    {/each}
</div>


<dialog
    bind:this={dialogEl}
    class="p-0 backdrop:bg-black backdrop:bg-opacity-40 max-w-md w-full"
    onclick={(e) => { if (e.target === dialogEl) closeDialog(); }}
    aria-labelledby="edit-item-dialog-title"
>
    <div class="bg-core-grey-600 text-white p-2 relative">
        <h2 id="edit-item-dialog-title" class="font-bold">Edit {item.Name} ({item.Date})</h2>
        <button
            onclick={closeDialog}
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

    <div class="p-4 space-y-2">
        <p class="px-2">Selected Container: {selectedContainer}</p>

        <div class="flex flex-col pl-2">
            <div>
                <button class="px-2 border border-red-500 rounded" onclick={removeContainer}>Remove</button>

                {#each $appState?.Freezers ?? [] as freezer}
                    {#if freezer.Name !== freezerName}
                        <button class="px-2 border border-green-500 rounded" onclick={() => moveContainer(freezer.Name)}>Move to {freezer.Name}</button>
                    {/if}
                {/each}
            </div>
        </div>
    </div>
</dialog>
