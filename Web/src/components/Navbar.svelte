<script lang="ts">
    import { currPage } from "$lib/store";
</script>

<nav class="p-4 flex items-center justify-between border-b-2 border-black/20">
    <h1 class="font-bold text-2xl">
        <span style="color: #f73c00">Sve</span><span style="color: #08afd8">Go</span> List
    </h1>
    <div class="flex gap-4">
        <a class="nav-link {$currPage === 'home' ? 'active' : null}" href="/">
            Home
        </a>
        <a class="nav-link {$currPage === 'login' ? 'active' : null}" href="/login">
            Login
        </a>
    </div>
</nav>

<style>
    .nav-link {
        @apply px-4 py-1 rounded font-semibold;
        transition-duration: 300ms;
    }

    .nav-link:hover {
        @apply bg-sky-800;
    }

    .nav-link.active {
        @apply bg-sky-600;
    }
</style>