<script lang="ts">
    import { onMount } from "svelte";
    import { appState } from "$lib/stores";

    import Freezer from "./Freezer.svelte";

    const url = "/state";

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

    {#if !$appState}
        <p>Loading...</p>
    {:else}
        {#each $appState.Freezers as freezer}
            <Freezer {freezer}/>
        {/each}
    {/if}
</div>