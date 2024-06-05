<script lang="ts">
	import { currPage } from "$lib/store";
	import { isAuthenticated, requestLogout } from "$lib/auth";
</script>

<nav class="p-4 flex items-center justify-between border-b-2 border-white/10">
	<a href="/">
		<h1 class="font-bold text-2xl select-none">
			<span style="color: #f73c00">Sve</span><span style="color: #08afd8">Go</span> List
		</h1>
	</a>
	<div class="flex gap-4">
		{#if $isAuthenticated}
			<a class="nav-link-base nav-link-normal {$currPage === 'home' ? 'active' : null}" href="/"> Home </a>
			<button type="button" class="nav-link-base" on:click={requestLogout}> Logout </button>
		{:else}
			<a class="nav-link-base nav-link-signin {$currPage === 'signin' ? 'active' : null}" href="/signin">
				Sign In
			</a>
		{/if}
	</div>
</nav>

<style>
	.nav-link-base {
		@apply px-4 py-1 rounded font-semibold border-2 border-transparent;
		transition-duration: 300ms;
	}

	.nav-link-signin {
		@apply bg-amber-600;
	}

	.nav-link-signin:not(.active):hover {
		@apply bg-amber-800;
	}

	.nav-link-normal:not(.active):hover {
		@apply bg-sky-800;
	}

	.nav-link-base.active {
		@apply border-b-sky-600;
	}

	.nav-link-signin.active {
		@apply border-b-amber-600;
		background: none;
	}
</style>
