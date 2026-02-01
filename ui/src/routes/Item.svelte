<script lang="ts">
    import { appState, type FreezerItem, errorMessage, loadingState, clearErrorAfterDelay } from "$lib/stores";
    import { sortContainerNames } from "$lib/containerNames";
    import RenameItemModal from "./RenameItemModal.svelte";

    let { item, freezerName }: { item: FreezerItem; freezerName: string } = $props();

    let selectedContainer: string | undefined = $state(undefined);
    let dialogEl: HTMLDialogElement;
    let isRenameModalOpen = $state(false);

    let isRemoving = $derived($loadingState.removeContainer);
    let isMoving = $derived($loadingState.moveContainer);
    let isOperationPending = $derived(isRemoving || isMoving);

    const openDialog = (container: string) => {
        selectedContainer = container;
        dialogEl?.showModal();
    };

    const closeDialog = () => {
        dialogEl?.close();
        selectedContainer = undefined;
    };

    const removeContainer = async () => {
        const sure = confirm(`Are you sure you want to remove ${selectedContainer}?`);
        if (!sure) {
            return;
        }
        const url = "/remove";

        loadingState.update(s => ({ ...s, removeContainer: true }));
        errorMessage.set(null);

        try {
            const res = await fetch(url, {
                method: "POST",
                body: JSON.stringify({Container: selectedContainer}),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!res.ok) {
                throw new Error(`Failed to remove container: ${res.status} ${res.statusText}`);
            }

            const d = await res.json();
            appState.set(d);
            closeDialog();
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Failed to remove container';
            errorMessage.set(message);
            clearErrorAfterDelay();
        } finally {
            loadingState.update(s => ({ ...s, removeContainer: false }));
        }
    };

    const moveContainer = async (newFreezer: string) => {
        const url = "/move";

        loadingState.update(s => ({ ...s, moveContainer: true }));
        errorMessage.set(null);

        try {
            const res = await fetch(url, {
                method: "POST",
                body: JSON.stringify({Container: selectedContainer, newFreezer}),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!res.ok) {
                throw new Error(`Failed to move container: ${res.status} ${res.statusText}`);
            }

            const d = await res.json();
            appState.set(d);
            closeDialog();
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Failed to move container';
            errorMessage.set(message);
            clearErrorAfterDelay();
        } finally {
            loadingState.update(s => ({ ...s, moveContainer: false }));
        }
    };
</script>


<div class="flex gap-2 items-center">
    <span>{item.Name} ({item.Date})</span>
    <button
        onclick={() => isRenameModalOpen = true}
        class="p-1 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded transition-colors"
        aria-label="Rename {item.Name}"
    >
        <!-- Heroicons pencil-square -->
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
            <path stroke-linecap="round" stroke-linejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
        </svg>
    </button>
    {#each sortContainerNames(item.Containers) as container}
        <button onclick={() => openDialog(container)}
                class="px-3 py-1 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 active:bg-gray-300 transition-colors text-sm"
                aria-label="Edit container {container} for {item.Name}">{container}</button>
    {/each}
</div>

<RenameItemModal {item} bind:isOpen={isRenameModalOpen} />


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
            class="absolute top-2 right-2 p-1 rounded-md hover:bg-white/20 active:bg-white/30 transition-colors"
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

        <div class="flex flex-col pl-2 gap-2">
            <button
                class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 active:bg-red-800 transition-colors shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
                onclick={removeContainer}
                disabled={isOperationPending}
            >
                {#if isRemoving}
                    Removing...
                {:else}
                    Remove
                {/if}
            </button>

            Move to: 
            
            {#each $appState?.Freezers ?? [] as freezer}
                {#if freezer.Name !== freezerName}
                    <button
                        class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 active:bg-blue-800 transition-colors shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
                        onclick={() => moveContainer(freezer.Name)}
                        disabled={isOperationPending}
                    >
                        {#if isMoving}
                            Moving...
                        {:else}
                            {freezer.Name}
                        {/if}
                    </button>
                {/if}
            {/each}
        </div>
    </div>
</dialog>
