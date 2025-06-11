<script lang="ts">
    import '../app.css';
    import { page } from '$app/stores';
    import { currentTheme } from '$lib/stores/theme';
    import { onMount } from 'svelte';

    // Theme management
    const themeFiles = import.meta.glob('/src/themes/*.css', {eager: true, query: '?raw'});
    const themes = Object.keys(themeFiles).map(path => {
        return path.split('/').pop()?.split('.')[0] || '';
    });

    let isThemeLoaded = $state(false);
    let currentThemeLink: HTMLLinkElement | null = null;

    async function applyTheme(themeName: string) {
        if (currentThemeLink) {
            currentThemeLink.remove();
        }

        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = `/src/themes/${themeName}.css`;

        const loadPromise = new Promise((resolve) => {
            link.onload = resolve;
        });

        document.head.appendChild(link);
        currentThemeLink = link;

        await loadPromise;
        isThemeLoaded = true;
    }


    // Subscribe to theme changes
    currentTheme.subscribe((theme) => {
        if (typeof window !== 'undefined') {
            isThemeLoaded = false;
            applyTheme(theme);
        }
    });

    onMount(() => {
        const storedTheme = localStorage.getItem('selectedTheme');
        if (storedTheme) {
            currentTheme.set(storedTheme);
        } else {
            const defaultTheme = themes.find(theme => theme.toLowerCase() === 'default') || themes[0] || '';
            currentTheme.set(defaultTheme);
        }
    });

</script>

<div class="min-h-screen">
    <!-- Sidebar -->
    <div class="fixed inset-y-0 left-0 w-64 shadow-lg" style="background-color: var(--background-0);">
        <div class="flex flex-col h-full">
            <!-- Logo -->
            <div class="p-4 border-b" style="border-color: var(--accent-color);">
                <h1 class="text-xl font-bold" style="color: var(--foreground-0);">HyveKeeper</h1>
            </div>

            <!-- Navigation -->
            <nav class="flex-1 p-4">
                <ul class="space-y-2">
                    <li>
                        <a
                                href="/"
                                class="flex items-center px-4 py-2"
                                style="border-radius: var(--radius-md);"
                                class:active={$page.url.pathname === '/'}
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6l9-4 9 4v12l-9 4-9-4V6z" />
                            </svg>
                            Host
                        </a>
                    </li>
                    <li>
                        <a
                                href="/settings"
                                class="flex items-center px-4 py-2"
                                style="border-radius: var(--radius-md);"
                                class:active={$page.url.pathname === '/settings'}
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            </svg>
                            Settings
                        </a>
                    </li>
                </ul>
            </nav>
        </div>
    </div>

    <!-- Main Content -->
    <div class="ml-64">
        <slot />
    </div>
</div>

<style>
    :global(body) {
        background-color: var(--background-0);
        color: var(--foreground-0);
    }

    .active {
        background-color: var(--primary-color);
        color: white;
    }

    a:not(.active) {
        color: var(--foreground-0);
    }

    a:not(.active):hover {
        background-color: var(--secondary-color);
        opacity: 0.7;
    }
</style>
