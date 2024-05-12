<script lang="ts">
	import { currPage } from "$lib/store";
	import MetaTitle from "$lib/MetaTitle.svelte";
	import Icon from "@iconify/svelte";

	currPage.set("login");
</script>

<MetaTitle title="Login" />

<h1 class="text-3xl font-bold text-center">Login</h1>
<form class="w-[350px] mx-auto mt-10">
	<div>
		<label for="i_usr" class="block font-semibold">Username</label>
		<input id="i_usr" name="username" type="text" class="w-full text-black p-1 rounded" />
	</div>
	<div class="mt-4">
		<label for="i_pw" class="block font-semibold">Username</label>
		<input id="i_pw" name="password" type="password" class="w-full text-black p-1 rounded" />
	</div>
	<button
		type="submit"
		class="group relative mt-8 bg-orange-700 py-1.5 rounded font-bold w-full flex gap-1 justify-center items-center
			   hover:bg-amber-900 duration-300"
	>
		<span class="block"> Login </span>
		<Icon
			icon="mingcute:arrow-right-fill"
			class="size-5 absolute duration-300 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100
					   right-10 group-hover:right-2 scale-x-[25%] group-hover:scale-x-100"
		/>
	</button>
</form>
