<script lang="ts">
    import { appState, type FreezerItem, errorMessage, loadingState, clearErrorAfterDelay } from "$lib/stores";

    let { item, isOpen = $bindable(false) }: { item: FreezerItem; isOpen: boolean } = $props();

    let newName = $state(item.Name);
    let dialogEl: HTMLDialogElement;

    let isRenaming = $derived($loadingState.renameItem);

    $effect(() => {
        if (isOpen) {
            newName = item.Name;
            dialogEl?.showModal();
        } else {
            dialogEl?.close();
        }
    });

    const closeDialog = () => {
        isOpen = false;
    };

    const renameItem = async () => {
        if (newName.trim() === "" || newName === item.Name) {
            closeDialog();
            return;
        }

        const url = "/rename";

        loadingState.update(s => ({ ...s, renameItem: true }));
        errorMessage.set(null);

        try {
            const res = await fetch(url, {
                method: "POST",
                body: JSON.stringify({
                    OldName: item.Name,
                    OldDate: item.Date,
                    NewName: newName.trim()
                }),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!res.ok) {
                throw new Error(`Failed to rename item: ${res.status} ${res.statusText}`);
            }

            const d = await res.json();
            appState.set(d);
            closeDialog();
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Failed to rename item';
            errorMessage.set(message);
            clearErrorAfterDelay();
        } finally {
            loadingState.update(s => ({ ...s, renameItem: false }));
        }
    };

    const handleKeydown = (e: KeyboardEvent) => {
        if (e.key === 'Enter') {
            renameItem();
        }
    };
</script>

<dialog
    bind:this={dialogEl}
    class="p-0 backdrop:bg-black backdrop:bg-opacity-40 max-w-md w-full"
    onclick={(e) => { if (e.target === dialogEl) closeDialog(); }}
    onclose={closeDialog}
    aria-labelledby="rename-item-dialog-title"
>
    <div class="bg-core-grey-600 text-white p-2 relative">
        <h2 id="rename-item-dialog-title" class="font-bold">Rename Item</h2>
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

    <div class="flex flex-col px-4 py-4 gap-4">
        <div class="flex flex-col pl-2">
            <label for="newName">New Name:</label>
            <input
                type="text"
                id="newName"
                class="form-input"
                bind:value={newName}
                onkeydown={handleKeydown}
            />
        </div>

        <div class="flex gap-2 justify-end">
            <button
                onclick={closeDialog}
                disabled={isRenaming}
                class="px-4 py-2 border border-gray-400 text-gray-700 rounded-md hover:bg-gray-100 active:bg-gray-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
                Cancel
            </button>
            <button
                onclick={renameItem}
                disabled={isRenaming}
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 active:bg-blue-800 transition-colors shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
                {#if isRenaming}
                    Renaming...
                {:else}
                    Rename
                {/if}
            </button>
        </div>
    </div>
</dialog>
