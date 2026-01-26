<script lang="ts">
    import { onMount } from "svelte";
    import { appState } from "$lib/stores";

    import Freezer from "./Freezer.svelte";
    import GroupByContents from "./GroupByContents.svelte";
    import GroupByNumber from "./GroupByNumber.svelte";
    import GroupByDate from "./GroupByDate.svelte";

    const url = "/state";

    type TabId = "location" | "contents" | "number" | "date";
    let activeTab: TabId = $state("location");

    const tabs: { id: TabId; label: string }[] = [
        { id: "location", label: "Location" },
        { id: "contents", label: "Contents" },
        { id: "number", label: "Label" },
        { id: "date", label: "Date" },
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

    <div class="flex items-center gap-2" role="tablist" aria-label="View options">
        <span class="text-gray-600 py-2">Group by</span>
        {#each tabs as tab}
            <button
                role="tab"
                aria-selected={activeTab === tab.id}
                aria-controls="tab-panel"
                id="tab-{tab.id}"
                class="px-4 py-2 -mb-px {activeTab === tab.id ? 'border-b-2 border-blue-500 font-bold' : 'border-b border-gray-300 text-gray-600 hover:text-gray-800'}"
                onclick={() => activeTab = tab.id}
            >
                {tab.label}
            </button>
        {/each}
    </div>

    <div id="tab-panel" role="tabpanel" aria-labelledby="tab-{activeTab}">
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
      {:else if activeTab === "date"}
          <GroupByDate />
      {/if}
    </div>
</div>
