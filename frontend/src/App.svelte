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
    GetStats,
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
    cacheLoaded: false,
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

  // Handle glyph click (copy to clipboard)
  const handleGlyphClick = async (glyph: GlyphMatch) => {
    try {
      await CopyToClipboard(glyph.glyph);
      showToast();
    } catch (error) {
      console.error("Failed to copy to clipboard:", error);
    }
  };

  // Toggle favorite
  const handleToggleFavorite = async (glyphId: number, event: Event) => {
    event.stopPropagation();
    try {
      await ToggleFavorite(glyphId);

      filteredGlyphs = filteredGlyphs.map((g) =>
        g.id === glyphId ? { ...g, isFavorite: !g.isFavorite } : g
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
    const scrollPercentage =
      (target.scrollTop + target.clientHeight) / target.scrollHeight;

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
      viewingFavorites = true;
      isLoading = false;
    } catch (error) {
      console.error("Failed to load favorites:", error);
      isLoading = false;
    }
  };

  // Handle category change
  const handleCategoryChange = async (category: string) => {
    selectedCategory = category;
    showCategoryFilter = false;
    viewingFavorites = false;
    await loadGlyphs(true);
  };

  // Clear all filters
  const clearFilters = async () => {
    searchTerm = "";
    selectedCategory = "";
    viewingFavorites = false;
    await loadGlyphs(true);
  };

  // Reactive search with debouncing
  $: {
    if (searchTimeout) clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      if (!viewingFavorites) {
        loadGlyphs(true);
      }
    }, 300);
  }
</script>

<div id="app">
  <!-- Custom Title Bar -->
  <div class="title-bar draggable">
    <div class="title draggable">
      <img src={MonoNF} id="nf-icon" alt="Gylte" width="31" height="31" />
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
      placeholder=""
      bind:value={searchTerm}
      disabled={isLoading}
    />
  </div>

  <div class="toolbar">
    <div class="filters">
      <button
        class="filter-btn"
        on:click={() => (showCategoryFilter = !showCategoryFilter)}
        title="Filter by category"
      >
         {selectedCategory || "Categories"}
      </button>

      <button
        class="filter-btn favorites-btn"
        on:click={showFavorites}
        title="Show favorites"
      >
         {stats.totalFavorites}
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
        on:click={() => handleCategoryChange("")}
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
              title={item.isFavorite
                ? "Remove from favorites"
                : "Add to favorites"}
            >
              {item.isFavorite ? "★" : "☆"}
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
    background: var(--bg-color);
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
    background: #0c0e13;
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
    border-radius: 9px;
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

  /* Add remaining styles for search, toolbar, grid, etc. */
  .search-container {
    padding: 1rem;
    margin-top: 1empx;
    background: var(--bg-color);
  }

  .search-input {
    font-family: var(--font-sans);
    font-weight: 600;
    width: 40%;
    padding: 0.75rem;
    font-size: 0.7rem;
    border: none;
    border-radius: 11px;
    background: rgba(0, 0, 0, 0.2);
    color: #c5c8c6;
    backdrop-filter: blur(12px);
  }

  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 2rem 0.5rem;
    font-size: 12px;
    background: var(--bg-color);
  }

  .filters {
    font-size: 12px;
    display: flex;
    gap: 0.5rem;
  }

  .filter-btn {
    font-size: 101%;
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 6px;
    background: rgba(5, 5, 5, 0.5);
    color: #c5c8c6;
    cursor: pointer;
    transition: all 160ms ease-in;
  }

  .filter-btn:hover {
    background: rgba(8, 60, 73, 0.8);
  }

  .stats {
    display: flex;
    gap: 1rem;
    color: #8b8b8b;
    font-size: 0.875rem;
  }

  .glyph-grid-container {
    flex: 1;
    overflow-y: auto;
    padding: 0 1rem 1rem;
    background: var(--bg-color);
  }

  .glyph-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    margin-top: 0.5rem;
    gap: 1rem;
  }

  .glyph-card {
    position: relative;
    padding: 1rem;
    border-radius: 8px;
    background: rgba(5, 5, 5, 0.5);
    cursor: pointer;
    transition: all 210ms ease-in;
    text-align: center;
    max-height: 100px;
  }

  .glyph-card:hover {
    background: rgba(8, 60, 73, 0.6);
    transform: translateY(-2px);
  }

  .favorite-btn {
    position: absolute;
    top: 0.25rem;
    right: 0.25rem;
    background: none;
    border: none;
    font-size: 1rem;
    cursor: pointer;
    opacity: 0.5;
    transition: opacity 160ms ease;
  }

  .favorite-btn:hover,
  .favorite-btn.active {
    opacity: 1;
  }

  .glyph-icon {
    display: block;
    font-size: 2rem;
    margin-bottom: -1.5rem;
    transition: all 125ms ease-in;
  }

  .glyph-card:hover .glyph-icon {
    margin-bottom: 0.2rem;
  }

  .glyph-name {
    display: block;
    font-size: 0.6rem;
    color: #c0c0c000;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: normal;
    transition: all 230ms ease-in;
  }

  .glyph-card:hover .glyph-name {
    display: block;
    color: #c0c0c0;
  }

  .loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
  }

  .spinner {
    width: 48px;
    height: 48px;
    border: 3px solid rgba(255, 255, 255, 0.1);
    border-top-color: #083c49;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .toast {
    position: fixed;
    bottom: 2rem;
    /* right: 2rem; */
    padding: 1rem;
    background: rgba(5, 5, 5, 0.9);
    border-radius: 11px;
    backdrop-filter: blur(12px);
  }

  .no-results {
    text-align: center;
    padding: 2rem;
    color: #8b8b8b;
  }

  .category-dropdown {
    position: absolute;
    top: 160px;
    left: 1rem;
    background: rgba(5, 5, 5, 0.95);
    border-radius: 8px;
    padding: 0.5rem;
    max-height: 300px;
    overflow-y: auto;
    z-index: 100;
    backdrop-filter: blur(12px);
  }

  .category-item {
    display: block;
    width: 100%;
    padding: 0.5rem 1rem;
    text-align: left;
    background: none;
    border: none;
    color: #c5c8c6;
    cursor: pointer;
    border-radius: 4px;
    transition: all 160ms ease;
  }

  .category-item:hover,
  .category-item.active {
    background: rgba(8, 60, 73, 0.8);
  }

  .count {
    color: #8b8b8b;
    font-size: 0.875rem;
  }
</style>
