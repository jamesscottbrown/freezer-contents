<script>
    import {onMount} from "svelte";
    import {state} from "$lib/stores.ts";

    import Freezer from "./Freezer.svelte";


    const url = "/state";


    const getData = () => {

        fetch(url)
            .then(res => res.json())
            .then(d => $state = d);

    }

    $: console.log($state);

    onMount(getData)

</script>


<div class="container flex flex-col gap-6">
    <h1 class="font-bold text-2xl">Freezer contents</h1>

    {#if !$state}
        <p>Loading...</p>
    {:else}
        {#each $state.Freezers as freezer}
            <Freezer {freezer}/>
        {/each}
    {/if}
</div>