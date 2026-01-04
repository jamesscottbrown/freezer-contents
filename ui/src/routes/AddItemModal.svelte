<script lang="ts">
    import { Dialog, DialogOverlay, DialogTitle } from "@rgossiaux/svelte-headlessui";
    import Select from "svelte-select";
    import {state} from "$lib/stores.ts";

    export let freezerName: string;
    export let isOpen: boolean;

    type ContainerType = {
        value: string;
        label: string;
    };

    const containerTypes: ContainerType[] = [
        { value: "containers", label: "Containers" },
        { value: "bowls", label: "Bowls" },
        { value: "jars", label: "Jars" },
        { value: "numContainers", label: "Other (unlabelled)" }
    ];

    let selectedContainerType: ContainerType = containerTypes[0];

    const generateContainerNames = () => {
        let containerNames = [];

        let numNamesToGenerate = +numContainers;

        let usedContainers: string[] = [];
        for (const freezer of $state.Freezers){
            for (const content of freezer.Contents){
                usedContainers = [...usedContainers, ...content.Containers];
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
        containers = "";
        bowls = "";
        jars = "";
        numContainers = "";
        selectedContainerType = containerTypes[0];
    }

    const addItem = () => {
        const url = "/add";

        let containerNames = [];

        if (selectedContainerType.value === "containers"){
            containerNames = containers.split(",").map(n => n.trim());
        } else if (selectedContainerType.value === "bowls"){
            containerNames = bowls.split(",").map(n => "bowl_" + n.trim());
        } else if (selectedContainerType.value === "jars"){
            containerNames = jars.split(",").map(n => "jar_" + n.trim());
        } else {
            containerNames = generateContainerNames();
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
            .then(d => $state = d);

        isOpen = false;
    }

    let itemName = "";
    let containers = "";
    let bowls = "";
    let jars = "";
    let numContainers = "";

</script>

<Dialog open={isOpen} on:close={() => (isOpen = false)}
        class="relative z-10 overflow-y-auto">
    <DialogOverlay class="fixed inset-0 bg-black bg-opacity-40"/>

    <!--Full-screen container to center the panel -->
    <div class="fixed inset-0 flex items-center justify-center p-4 pointer-events-none">
        <div
                class="inline-block w-full max-w-md my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl space-y-2 pb-4 pointer-events-auto"
        >
            <div class="bg-core-grey-600 text-white p-2 relative">
                <DialogTitle>Add Item to {freezerName}</DialogTitle>
                <button
                        on:click={() => (isOpen = false)}
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

            <div class="flex flex-col px-4 gap-4">

                <div class="flex flex-col pl-2">
                    <label for="name">Name:</label>
                    <input type="text" id="name" class="form-input" bind:value={itemName}/>
                </div>

                <div class="flex flex-col pl-2">
                    <label for="containerType">Container Type:</label>
                    <Select
                        items={containerTypes}
                        bind:value={selectedContainerType}
                        clearable={false}
                        searchable={false}
                    />
                </div>

                <div class="flex flex-col pl-2">
                    {#if selectedContainerType.value === "containers"}
                        <label for="containers">Numbers (comma-separated):</label>
                        <input type="text" id="containers" class="form-input" bind:value={containers}/>
                    {:else if selectedContainerType.value === "bowls"}
                        <label for="bowls">Numbers (comma-separated):</label>
                        <input type="text" id="bowls" class="form-input" bind:value={bowls}/>
                    {:else if selectedContainerType.value === "jars"}
                        <label for="jars">Numbers (comma-separated):</label>
                        <input type="text" id="jars" class="form-input" bind:value={jars}/>
                    {:else if selectedContainerType.value === "numContainers"}
                        <label for="numContainers">Number of unlabelled containers:</label>
                        <input type="text" id="numContainers" class="form-input" bind:value={numContainers}/>
                    {/if}
                </div>

                <button on:click={clearFields} class="px-2 border border-red-500 rounded">Clear</button>

                <button on:click={addItem} class="px-2 border border-green-500 rounded">Add</button>
            </div>

        </div>
    </div>
</Dialog>
