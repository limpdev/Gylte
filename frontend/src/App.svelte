<script lang="ts">
    import { onMount } from "svelte";
    import glyphData from "./glyphs.json";
    import AniToast from "./comps/aniToast.svelte";
    import { CopyToClipboard } from "../wailsjs/go/main/App";
    import { WindowMinimise, Quit } from "../wailsjs/runtime";

    // Define the type for a single glyph object
    interface Glyph {
        name: string;
        glyph: string;
    }

    let searchTerm = "";
    let filteredGlyphs: Glyph[] = [];
    let toastVisible = false;

    // Load all glyphs on initial render and remove duplicates
    onMount(() => {
        // Remove duplicates by name, keeping the first occurrence
        const uniqueGlyphs = glyphData.filter((item, index, self) => index === self.findIndex((g) => g.name === item.name));
        filteredGlyphs = uniqueGlyphs;
    });

    // Reactive statement to filter glyphs when search term changes
    $: {
        const sourceData = glyphData.filter((item, index, self) => index === self.findIndex((g) => g.name === item.name));

        if (searchTerm.trim() === "") {
            filteredGlyphs = sourceData;
        } else {
            const lowercasedFilter = searchTerm.toLowerCase();
            filteredGlyphs = sourceData.filter((item) => item.name.toLowerCase().includes(lowercasedFilter));
        }
    }

    const handleGlyphClick = (glyph: string) => {
        CopyToClipboard(glyph);
        // Show "Copied!" toast message
        toastVisible = true;
        setTimeout(() => {
            toastVisible = false;
        }, 1900); // Hide after the animation is almost done
    };
</script>

<div id="app">
    <!-- Custom Title Bar -->
    <div class="title-bar draggable">
        <div class="title"></div>
        <div class="spacer"></div>
        <div class="window-controls">
            <button on:click={WindowMinimise}>-</button>
            <button on:click={Quit}>×</button>
        </div>
    </div>

    <!-- Search Input -->
    <div class="search-container draggable">
        <input type="text" autocomplete="on" class="search-input" placeholder="" bind:value={searchTerm} />
    </div>

    <!-- Glyph Grid -->
    <div class="glyph-grid draggable">
        {#each filteredGlyphs as item, index (index)}
            <div
                class="glyph-card"
                title={`Click to copy "${item.glyph}"`}
                on:click={() => handleGlyphClick(item.glyph)}
                on:keydown={(e) => e.key === "Enter" && handleGlyphClick(item.glyph)}
                role="button"
                tabindex="0"
            >
                <span class="glyph-icon">{item.glyph}</span>
                <span class="glyph-name">{item.name}</span>
            </div>
        {/each}
    </div>

    <!-- "Copied!" Toast Notification -->
    {#if toastVisible}
        <div class="toast">
            <!-- IMPORT THE SVG FOR THE TOAST HERE AS A COMPONENT -->
            <AniToast width="12" height="12" fill="#45a847" class="animatedToast" />
        </div>
    {/if}
</div>

<style>
    /* You'll need to move your existing CSS here or import it */
    /* This is where your current style.css content should go */
</style>
