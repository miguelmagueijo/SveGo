<script>
	import "../app.css";
	import Navbar from "../components/Navbar.svelte";
	import Footer from "../components/Footer.svelte";
	import { isAuthenticated } from "$lib/auth";

	export let data;

	isAuthenticated.set(data.isAuthenticated);
</script>

<div class="flex flex-col min-h-screen">
	<Navbar />
	<main class="mt-10 mb-auto">
		<slot />
	</main>
	<Footer />
</div>
