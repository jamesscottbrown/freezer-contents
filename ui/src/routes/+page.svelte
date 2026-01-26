<script lang="ts">
    import { onMount } from "svelte";
    import { appState, errorMessage, loadingState, clearErrorAfterDelay } from "$lib/stores";

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

    const getData = async () => {
        loadingState.update(s => ({ ...s, initialLoad: true }));
        errorMessage.set(null);

        try {
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error(`Failed to load data: ${res.status} ${res.statusText}`);
            }
            const d = await res.json();
            appState.set(d);
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Failed to load data';
            errorMessage.set(message);
            clearErrorAfterDelay();
        } finally {
            loadingState.update(s => ({ ...s, initialLoad: false }));
        }
    }

    $effect(() => {
        console.log($appState);
    });

    onMount(getData);
</script>


<div class="container flex flex-col gap-6 pt-4">
    <h1 class="font-bold text-2xl">Freezer contents</h1>

    {#if $errorMessage}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <span class="block sm:inline">{$errorMessage}</span>
            <button
                class="absolute top-0 bottom-0 right-0 px-4 py-3"
                onclick={() => errorMessage.set(null)}
                aria-label="Dismiss error"
            >
                <svg class="fill-current h-6 w-6 text-red-500" role="button" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                    <title>Close</title>
                    <path d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z"/>
                </svg>
            </button>
        </div>
    {/if}

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
      {#if $loadingState.initialLoad}
          <p>Loading...</p>
      {:else if !$appState}
          <p class="text-gray-500">No data available. {$errorMessage ? '' : 'Try refreshing the page.'}</p>
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
