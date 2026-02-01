<script lang="ts">
    import { appState, type FreezerItem, errorMessage, loadingState, clearErrorAfterDelay } from "$lib/stores";

    let { item, isOpen = $bindable(false) }: { item: FreezerItem; isOpen: boolean } = $props();

    let newName = $state(item.Name);
    let newDate = $state(item.Date);
    let dialogEl: HTMLDialogElement;

    let isEditing = $derived($loadingState.editItem);

    $effect(() => {
        if (isOpen) {
            newName = item.Name;
            newDate = item.Date;
            dialogEl?.showModal();
        } else {
            dialogEl?.close();
        }
    });

    const closeDialog = () => {
        isOpen = false;
    };

    const editItem = async () => {
        const trimmedName = newName.trim();
        const trimmedDate = newDate.trim();
        
        if (trimmedName === "" || trimmedDate === "") {
            return;
        }
        
        if (trimmedName === item.Name && trimmedDate === item.Date) {
            closeDialog();
            return;
        }

        const url = "/edit";

        loadingState.update(s => ({ ...s, editItem: true }));
        errorMessage.set(null);

        try {
            const res = await fetch(url, {
                method: "POST",
                body: JSON.stringify({
                    OldName: item.Name,
                    OldDate: item.Date,
                    NewName: trimmedName,
                    NewDate: trimmedDate
                }),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!res.ok) {
                throw new Error(`Failed to edit item: ${res.status} ${res.statusText}`);
            }

            const d = await res.json();
            appState.set(d);
            closeDialog();
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Failed to edit item';
            errorMessage.set(message);
            clearErrorAfterDelay();
        } finally {
            loadingState.update(s => ({ ...s, editItem: false }));
        }
    };

    const handleKeydown = (e: KeyboardEvent) => {
        if (e.key === 'Enter') {
            editItem();
        }
    };
</script>

<dialog
    bind:this={dialogEl}
    class="p-0 backdrop:bg-black backdrop:bg-opacity-40 max-w-md w-full"
    onclick={(e) => { if (e.target === dialogEl) closeDialog(); }}
    onclose={closeDialog}
    aria-labelledby="edit-item-dialog-title"
>
    <div class="bg-core-grey-600 text-white p-2 relative">
        <h2 id="edit-item-dialog-title" class="font-bold">Edit Item</h2>
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
            <label for="newName">Name:</label>
            <input
                type="text"
                id="newName"
                class="form-input"
                bind:value={newName}
                onkeydown={handleKeydown}
            />
        </div>

        <div class="flex flex-col pl-2">
            <label for="newDate">Date:</label>
            <input
                type="date"
                id="newDate"
                class="form-input"
                bind:value={newDate}
                onkeydown={handleKeydown}
            />
        </div>

        <div class="flex gap-2 justify-end">
            <button
                onclick={closeDialog}
                disabled={isEditing}
                class="px-4 py-2 border border-gray-400 text-gray-700 rounded-md hover:bg-gray-100 active:bg-gray-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
                Cancel
            </button>
            <button
                onclick={editItem}
                disabled={isEditing}
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 active:bg-blue-800 transition-colors shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
                {#if isEditing}
                    Saving...
                {:else}
                    Save
                {/if}
            </button>
        </div>
    </div>
</dialog>
