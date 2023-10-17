<script lang="ts">
    import {Dialog, DialogDescription, DialogOverlay, DialogTitle} from "@rgossiaux/svelte-headlessui";
    import {state} from "$lib/stores.ts";

    export let freezerName: string;
    export let isOpen: boolean;

    const addItem = () => {
        const url = "/add";

        fetch(url, {
            method: "POST",
            body: JSON.stringify({
                Name: itemName,
                Date: (new Date()).toISOString().slice(0,10),
                Freezer: freezerName,
                Containers: containers.split(",").map(n => n.trim())

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


            <div class="flex flex-col pl-2">


                <label for="name">Name:</label>
                <input type="text" id="name" class="form-input" bind:value={itemName}/>

                <label for="containers">Containers:</label>
                <input type="text" id="containers" class="form-input" bind:value={containers}/>

            </div>

            <button on:click={addItem} class="px-2 border border-green-500 rounded">Add</button>

        </div>
    </div>
</Dialog>
