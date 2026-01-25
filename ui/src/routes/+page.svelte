<script lang="ts">
    import { onMount } from "svelte";
    import { appState } from "$lib/stores";

    import Freezer from "./Freezer.svelte";
    import GroupByContents from "./GroupByContents.svelte";
    import GroupByNumber from "./GroupByNumber.svelte";

    const url = "/state";

    type TabId = "location" | "contents" | "number";
    let activeTab: TabId = $state("location");

    const tabs: { id: TabId; label: string }[] = [
        { id: "location", label: "Group by Location" },
        { id: "contents", label: "Group by Contents" },
        { id: "number", label: "Group by Number" },
    ];

    const getData = () => {
        fetch(url)
            .then(res => res.json())
            .then(d => appState.set(d));
    }

    $effect(() => {
        console.log($appState);
    });

    onMount(getData);
</script>


<div class="container flex flex-col gap-6">
    <h1 class="font-bold text-2xl">Freezer contents</h1>

    <div class="flex gap-2 border-b border-gray-300">
        {#each tabs as tab}
            <button
                class="px-4 py-2 -mb-px {activeTab === tab.id ? 'border-b-2 border-blue-500 font-bold' : 'text-gray-600 hover:text-gray-800'}"
                onclick={() => activeTab = tab.id}
            >
                {tab.label}
            </button>
        {/each}
    </div>

    {#if !$appState}
        <p>Loading...</p>
    {:else if activeTab === "location"}
        {#each $appState.Freezers as freezer}
            <Freezer {freezer}/>
        {/each}
    {:else if activeTab === "contents"}
        <GroupByContents />
    {:else if activeTab === "number"}
        <GroupByNumber />
    {/if}
</div>