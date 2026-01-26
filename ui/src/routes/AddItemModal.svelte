<script lang="ts">
    import Select from "svelte-select";
    import { appState } from "$lib/stores";
    import { buildContainerNames } from "$lib/containerNames";

    let { freezerName, isOpen = $bindable(false) }: { freezerName: string; isOpen: boolean } = $props();

    type SelectOption = {
        value: string;
        label: string;
        prefix?: string;
    };

    // Build container type options from state, plus always add "numContainers"
    const containerTypeOptions = $derived<SelectOption[]>([
        ...($appState?.ContainerTypes ?? []).map(ct => ({
            value: ct.Value,
            label: ct.Label,
            prefix: ct.Prefix
        })),
        { value: "numContainers", label: "Other (unlabelled)" }
    ]);

    let selectedContainerType: SelectOption = $state({ value: "", label: "" });

    // Set default selection when options become available
    $effect(() => {
        if (containerTypeOptions.length > 0 && !selectedContainerType.value) {
            selectedContainerType = containerTypeOptions[0];
        }
    });

    let itemName = $state("");
    let containerInput = $state("");
    let numContainers = $state("");
    let dialogEl: HTMLDialogElement;

    $effect(() => {
        if (isOpen) {
            dialogEl?.showModal();
        } else {
            dialogEl?.close();
        }
    });

    const closeDialog = () => {
        isOpen = false;
    };

    const generateContainerNames = () => {
        let containerNames: string[] = [];

        let numNamesToGenerate = +numContainers;

        let usedContainers: string[] = [];
        if ($appState) {
            for (const freezer of $appState.Freezers){
                for (const content of freezer.Contents){
                    usedContainers = [...usedContainers, ...content.Containers];
                }
            }
        }

        let i = 0;
        while (numNamesToGenerate > 0){
            const name = `u-${i}`;
            if (!usedContainers.includes(name)){
                containerNames.push(name);
                numNamesToGenerate--;
            }
            i++;
        }

        return containerNames;
    }

    const clearFields = () => {
        itemName = "";
        containerInput = "";
        numContainers = "";
        if (containerTypeOptions.length > 0) {
            selectedContainerType = containerTypeOptions[0];
        }
    }

    const addItem = () => {
        const url = "/add";

        let containerNames: string[] = [];

        if (selectedContainerType.value === "numContainers") {
            containerNames = generateContainerNames();
        } else {
            containerNames = buildContainerNames(containerInput, selectedContainerType.prefix ?? "");
        }

        fetch(url, {
            method: "POST",
            body: JSON.stringify({
                Name: itemName,
                Date: (new Date()).toISOString().slice(0,10),
                Freezer: freezerName,
                Containers: containerNames

            }),
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(res => res.json())
            .then(d => appState.set(d));

        closeDialog();
    }
</script>

<dialog
    bind:this={dialogEl}
    class="p-0 backdrop:bg-black backdrop:bg-opacity-40 max-w-md w-full"
    onclick={(e) => { if (e.target === dialogEl) closeDialog(); }}
    onclose={closeDialog}
    aria-labelledby="add-item-dialog-title"
>
    <div class="bg-core-grey-600 text-white p-2 relative">
        <h2 id="add-item-dialog-title" class="font-bold">Add Item to {freezerName}</h2>
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

    <div class="flex flex-col px-4 py-4 gap-4">
        <div class="flex flex-col pl-2">
            <label for="name">Name:</label>
            <input type="text" id="name" class="form-input" bind:value={itemName}/>
        </div>

        <div class="flex flex-col pl-2">
            <span id="containerType-label">Container Type:</span>
            <Select
                items={containerTypeOptions}
                value={selectedContainerType}
                on:change={(e) => selectedContainerType = e.detail}
                clearable={false}
                searchable={false}
                ariaAttributes={{ 'aria-labelledby': 'containerType-label' }}
            />
        </div>

        <div class="flex flex-col pl-2">
            {#if selectedContainerType.value === "numContainers"}
                <label for="numContainers">Number of unlabelled containers:</label>
                <input type="text" id="numContainers" class="form-input" bind:value={numContainers}/>
            {:else}
                <label for="containerInput">Numbers (comma-separated):</label>
                <input type="text" id="containerInput" class="form-input" bind:value={containerInput}/>
            {/if}
        </div>

        <button onclick={clearFields} class="px-2 border border-red-500 rounded">Clear</button>

        <button onclick={addItem} class="px-2 border border-green-500 rounded">Add</button>
    </div>
</dialog>
