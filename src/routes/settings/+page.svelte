<script lang="ts">
    import { currentTheme } from '$lib/stores/theme';
    
    const themeFiles = import.meta.glob('/src/themes/*.css');
    
    const themes = Object.keys(themeFiles).map(path => {
        return path.split('/').pop()?.split('.')[0] || '';
    }).sort((a, b) => {
        if (a.toLowerCase() === 'default') return -1;
        if (b.toLowerCase() === 'default') return 1;
        return a.localeCompare(b);
    });
    
    function handleThemeChange(event: Event) {
        const select = event.target as HTMLSelectElement;
        const newTheme = select.value;
        localStorage.setItem('selectedTheme', newTheme);
        currentTheme.set(newTheme);
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const systemTheme = prefersDark ? 'dark' : 'light';
    }
</script>

<div class="p-6">
    <h1 class="text-2xl font-bold mb-6">Settings</h1>

    <div class=" shadow p-6" style="background-color: var(--background-1); border-radius: var(--radius-lg);">
        <div class="flex items-center justify-between">
            <div>
                <h2 class="text-lg font-medium">Theme Selector</h2>
                <p class="text-sm" style="color: var(--foreground-2);">
                    Select your theme
                </p>
            </div>
            
            <select
                value={$currentTheme}
                on:change={handleThemeChange}
                class="mt-1 block w-64 border py-2 pl-3 pr-10 text-base focus:outline-none focus:ring-2 focus:ring-primary sm:text-sm"
                style="background-color: var(--background-2); color: var(--foreground-0); border-color: var(--foreground-1); border-radius: var(--radius-md);"
            >
                {#each themes as theme}
                    <option value={theme}>{theme}</option>
                {/each}
            </select>
        </div>
    </div>
</div>