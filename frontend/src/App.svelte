<script lang="ts">
    import { onMount } from "svelte";
    import AniToast from "./comps/aniToast.svelte";
<<<<<<< HEAD
    import { CopyToClipboard, GetGlyphs } from "../wailsjs/go/main/App";
=======
    import { CopyToClipboard, SearchGlyphs, GetAllGlyphs } from "../wailsjs/go/main/App";
>>>>>>> 6021e9c552a01e92c7674f270c3248a24c1ebf42
    import { WindowMinimise, Quit } from "../wailsjs/runtime";

    // Define the type for a single glyph object
    interface Glyph {
        name: string;
        glyph: string;
    }

    let searchTerm = "";
    let glyphs: Glyph[] = [];
    let toastVisible = false;
    let searchTimeout: NodeJS.Timeout;
    let isLoading = true;

    // Load all glyphs on initial render
    onMount(async () => {
<<<<<<< HEAD
        glyphs = await GetGlyphs("");
    });

    // Function to handle searching
    const handleSearch = async () => {
        glyphs = await GetGlyphs(searchTerm);
    };

    // Reactive statement to trigger search when searchTerm changes
    $: {
        const timeout = setTimeout(handleSearch, 300); // Debounce input
        () => clearTimeout(timeout);
=======
        try {
            // The Go backend automatically loads the embedded JSON on startup
            // So we just need to get all glyphs for the initial display
            filteredGlyphs = await GetAllGlyphs();
            isLoading = false;
        } catch (error) {
            console.error("Failed to load glyphs from Go backend:", error);
            isLoading = false;
        }
    });

    // Debounced search function for better performance
    const performSearch = async (query: string) => {
        try {
            filteredGlyphs = await SearchGlyphs(query);
        } catch (error) {
            console.error("Search failed:", error);
        }
    };

    // Reactive statement with debouncing for search
    $: {
        // Clear previous timeout
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        // Debounce search by 150ms to avoid excessive backend calls
        if (!isLoading) {
            searchTimeout = setTimeout(() => {
                performSearch(searchTerm);
            }, 150);
        }
>>>>>>> 6021e9c552a01e92c7674f270c3248a24c1ebf42
    }

    const handleGlyphClick = (glyph: string) => {
        CopyToClipboard(glyph);
        // Show "Copied!" toast message
        toastVisible = true;
        setTimeout(() => {
            toastVisible = false;
        }, 1900);
    };
</script>

<div id="app">
    <!-- Custom Title Bar -->
    <div class="title-bar draggable">
        <div class="title"></div>
        <div class="spacer"></div>
        <div class="window-controls">
            <button on:click={WindowMinimise}>-</button>
            <button on:click={Quit}>Ã—</button>
        </div>
    </div>

    <!-- Search Input -->
    <div class="search-container draggable">
<<<<<<< HEAD
        <input
            type="text"
            autocomplete="on"
            class="search-input"
            placeholder=""
            bind:value={searchTerm}
        />
    </div>

    <!-- Glyph Grid -->
    <div class="glyph-grid draggable">
        {#each glyphs as item, index (index)}
            <div
                class="glyph-card"
                title={`Click to copy "${item.glyph}"`}
                on:click={() => handleGlyphClick(item.glyph)}
                on:keydown={(e) =>
                    e.key === "Enter" && handleGlyphClick(item.glyph)}
                role="button"
                tabindex="0"
            >
                <span class="glyph-icon">{item.glyph}</span>
                <span class="glyph-name">{item.name}</span>
            </div>
        {/each}
    </div>
=======
        <input 
            type="text" 
            autocomplete="on" 
            class="search-input" 
            placeholder="" 
            bind:value={searchTerm}
            disabled={isLoading}
        />
    </div>

    <!-- Loading State -->
    {#if isLoading}
        <div class="loading">Loading glyphs...</div>
    {:else}
        <!-- Glyph Grid -->
        <div class="glyph-grid draggable">
            {#each filteredGlyphs as item (item.name)}
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
    {/if}
>>>>>>> 6021e9c552a01e92c7674f270c3248a24c1ebf42

    <!-- "Copied!" Toast Notification -->
    {#if toastVisible}
        <div class="toast">
<<<<<<< HEAD
            <AniToast
                width="12"
                height="12"
                fill="#45a847"
                class="animatedToast"
            />
=======
            <AniToast width="12" height="12" fill="#45a847" class="animatedToast" />
>>>>>>> 6021e9c552a01e92c7674f270c3248a24c1ebf42
        </div>
    {/if}
</div>

<style>
    .loading {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 200px;
        color: #888;
        font-size: 14px;
        background-color: #12121290;
    }
    
    /* Your existing CSS styles go here */
</style>