<script lang="ts">
    import { onMount } from "svelte";
    import AniToast from "./comps/aniToast.svelte";
    import MonoNF from "./comps/mono-nf.webp";
    import { 
        CopyToClipboard, 
        GetGlyphs, 
        ToggleFavorite,
        GetFavorites,
        GetCategories,
        GetStats
    } from "../wailsjs/go/main/App";
    import { WindowMinimise, Quit } from "../wailsjs/runtime";
    import type { main } from "../wailsjs/go/models";

    // Type definitions
    interface Glyph {
        id: number;
        name: string;
        glyph: string;
        category?: string;
    }

    interface GlyphMatch extends Glyph {
        score: number;
        isFavorite: boolean;
    }

    interface SearchResult {
        glyphs: GlyphMatch[];
        total: number;
        searchTime: number;
        hasMore: boolean;
        categories?: string[];
    }

    // State
    let searchTerm = "";
    let selectedCategory = "";
    let filteredGlyphs: GlyphMatch[] = [];
    let toastVisible = false;
    let searchTimeout: NodeJS.Timeout;
    let isLoading = true;
    let currentOffset = 0;
    let hasMore = false;
    let total = 0;
    let searchTime = 0;
    let viewingFavorites = false;
    
    // Categories
    let categories: Record<string, number> = {};
    let showCategoryFilter = false;

    // Stats
    let stats = {
        totalGlyphs: 0,
        totalFavorites: 0,
        totalCategories: 0,
        cacheLoaded: false
    };

    const LIMIT = 100;

    // Load initial data
    onMount(async () => {
        try {
            stats = await GetStats();
            categories = await GetCategories();
            await loadGlyphs(true);
            isLoading = false;
        } catch (error) {
            console.error("Failed to initialize app:", error);
            isLoading = false;
        }
    });

    // Load glyphs with pagination
    const loadGlyphs = async (reset: boolean = false) => {
        try {
            if (reset) {
                currentOffset = 0;
                filteredGlyphs = [];
            }

            const result: SearchResult = await GetGlyphs(
                searchTerm,
                selectedCategory,
                LIMIT,
                currentOffset
            );

            if (reset) {
                filteredGlyphs = result.glyphs;
            } else {
                filteredGlyphs = [...filteredGlyphs, ...result.glyphs];
            }

            total = result.total;
            hasMore = result.hasMore;
            searchTime = result.searchTime;
            currentOffset += result.glyphs.length;

        } catch (error) {
            console.error("Failed to load glyphs:", error);
        }
    };

    // Debounced search
    const performSearch = async () => {
        await loadGlyphs(true);
    };

    // Reactive search with debouncing
    $: {
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        if (!isLoading) {
            searchTimeout = setTimeout(() => {
                performSearch();
            }, 200);
        }
    }

    // Category filter change
    const handleCategoryChange = async (category: string) => {
        selectedCategory = category;
        showCategoryFilter = false;
        await loadGlyphs(true);
    };

    // Handle glyph click (copy to clipboard)
    const handleGlyphClick = async (glyph: GlyphMatch) => {
        await CopyToClipboard(glyph.glyph);
        showToast();
    };

    // Toggle favorite
    const handleToggleFavorite = async (glyphId: number, event: Event) => {
        event.stopPropagation();
        try {
            await ToggleFavorite(glyphId);
            
            filteredGlyphs = filteredGlyphs.map(g => 
                g.id === glyphId 
                    ? { ...g, isFavorite: !g.isFavorite }
                    : g
            );

            stats = await GetStats();
        } catch (error) {
            console.error("Failed to toggle favorite:", error);
        }
    };

    // Load more glyphs (infinite scroll)
    const loadMore = async () => {
        if (hasMore && !isLoading) {
            await loadGlyphs(false);
        }
    };

    // Show toast notification
    const showToast = () => {
        toastVisible = true;
        setTimeout(() => {
            toastVisible = false;
        }, 1900);
    };

    // Scroll handler for infinite scroll
    const handleScroll = (event: Event) => {
        const target = event.target as HTMLElement;
        const scrollPercentage = (target.scrollTop + target.clientHeight) / target.scrollHeight;
        
        if (scrollPercentage > 0.8 && hasMore && !isLoading) {
            loadMore();
        }
    };

    // Show only favorites
    const showFavorites = async () => {
        try {
            isLoading = true;
            const favs = await GetFavorites();
            filteredGlyphs = favs;
            total = favs.length;
            hasMore = false;
            searchTerm = "";
            selectedCategory = "";
            isLoading = false;
        } catch (error) {
            console.error("Failed to load favorites:", error);
            isLoading = false;
        }
    };

    // Clear all filters
    const clearFilters = async () => {
        searchTerm = "";
        selectedCategory = "";
        await loadGlyphs(true);
    };
</script>

<div id="app">
    <!-- Custom Title Bar -->
    <div class="title-bar">
        <div class="title draggable">
            <img 
                src={MonoNF}
                id="nf-icon" 
                alt="Gylte"
                width="31"
                height="31"
            />
        </div>
        <div class="spacer draggable"></div>
        <div class="window-controls">
            <button on:click={WindowMinimise} title="Minimize">−</button>
            <button on:click={Quit} title="Close">×</button>
        </div>
    </div>

    <!-- Search & Filters -->
    <div class="search-container draggable">
        <input 
            type="text" 
            class="search-input" 
            placeholder="🔍" 
            bind:value={searchTerm}
            disabled={isLoading}
        />
    </div>

    <div class="toolbar">
        <div class="filters">
            <button 
                class="filter-btn"
                on:click={() => showCategoryFilter = !showCategoryFilter}
                title="Filter by category"
            >
                 {selectedCategory || 'Categories'}
            </button>
            
            <button 
                class="filter-btn favorites-btn"
                on:click={showFavorites}
                title="Show favorites"
            >
                 {stats.totalFavorites}
            </button>

            {#if searchTerm || selectedCategory}
                <button 
                    class="filter-btn clear-btn"
                    on:click={clearFilters}
                    title="Clear filters"
                >
                    ✕
                </button>
            {/if}
        </div>

        <div class="stats">
            <span class="stat-item">{filteredGlyphs.length} / {total}</span>
            {#if searchTime > 0}
                <span class="stat-item">{(searchTime * 1000).toFixed(0)}ms</span>
            {/if}
        </div>
    </div>

    <!-- Category Dropdown -->
    {#if showCategoryFilter}
        <div class="category-dropdown">
            <button 
                class="category-item {selectedCategory === '' ? 'active' : ''}"
                on:click={() => handleCategoryChange('')}
            >
                All Categories
            </button>
            {#each Object.entries(categories).sort((a, b) => b[1] - a[1]) as [cat, count]}
                <button 
                    class="category-item {selectedCategory === cat ? 'active' : ''}"
                    on:click={() => handleCategoryChange(cat)}
                >
                    {cat} <span class="count">({count})</span>
                </button>
            {/each}
        </div>
    {/if}

    <!-- Loading State -->
    {#if isLoading && filteredGlyphs.length === 0}
        <div class="loading">
            <div class="spinner"></div>
        </div>
    {:else}
        <!-- Glyph Grid -->
        <div class="glyph-grid-container" on:scroll={handleScroll}>
            <div class="glyph-grid">
                {#each filteredGlyphs as item (item.id)}
                    <div
                        class="glyph-card"
                        title={item.name}
                        on:click={() => handleGlyphClick(item)}
                        on:keydown={(e) => e.key === "Enter" && handleGlyphClick(item)}
                        role="button"
                        tabindex="0"
                    >
                        <button 
                            class="favorite-btn {item.isFavorite ? 'active' : ''}"
                            on:click={(e) => handleToggleFavorite(item.id, e)}
                            title={item.isFavorite ? 'Remove from favorites' : 'Add to favorites'}
                        >
                            {item.isFavorite ? '⭐' : '☆'}
                        </button>
                        <span class="glyph-icon">{item.glyph}</span>
                        <span class="glyph-name">{item.name}</span>
                    </div>
                {/each}
            </div>

            <!-- Load More -->
            {#if hasMore}
                <div class="load-more">
                    <button on:click={loadMore}>Load More</button>
                </div>
            {/if}

            <!-- No Results -->
            {#if filteredGlyphs.length === 0 && !isLoading}
                <div class="no-results">
                    <p>No glyphs found</p>
                    {#if searchTerm || selectedCategory}
                        <button on:click={clearFilters}>Clear filters</button>
                    {/if}
                </div>
            {/if}
        </div>
    {/if}

    <!-- Toast Notification -->
    {#if toastVisible}
        <div class="toast">
            <AniToast width="12" height="12" fill="#45a847" />
        </div>
    {/if}
</div>

<style>
    /* ===== GLOBAL RESET ===== */
    * {
        box-sizing: border-box;
    }

    #app {
        display: flex;
        flex-direction: column;
        height: 100vh;
        overflow: hidden;
        background: transparent;
    }

    /* ===== TITLE BAR ===== */
    .title-bar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 8px;
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        z-index: 1000;
        -webkit-app-region: drag;
    }

    .title {
        display: flex;
        align-items: center;
        background: rgba(5, 5, 5, 0.55);
        backdrop-filter: blur(12px);
        -webkit-backdrop-filter: blur(12px);
        border-radius: 6px;
        padding: 0.4em;
        margin-top: 1.2em;
        margin-left: 0.5em;
        transition: all 210ms ease;
    }

    .title:hover {
        cursor: grab;
        opacity: 0.9;
    }

    .title:active {
        cursor: grabbing;
    }

    .spacer {
        flex: 1;
    }

    .window-controls {
        display: flex;
        gap: 8px;
        background: rgba(5, 5, 5, 0.65);
        backdrop-filter: blur(12px);
        -webkit-backdrop-filter: blur(12px);
        border-radius: 6px;
        margin-top: 1.2em;
        -webkit-app-region: no-drag;
    }

    .window-controls button {
        width: 25px;
        background: none;
        border: none;
        color: #8b8b8b;
        font-size: 19px;
        cursor: pointer;
        padding: 0.3em;
        border-radius: 6px;
        line-height: 1;
        transition: all 160ms ease;
    }

    .window-controls button:hover {
        transform: scale(1.12);
        color: #c5c8c6;
        background-color: #083c49;
    }
</style>