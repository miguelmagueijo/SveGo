<script>
	import Icon from "@iconify/svelte";
	import { createEventDispatcher } from "svelte";

	const dispatch = createEventDispatcher();

	export let title = "";
	export let isCompleted = false;
</script>

<div class="flex items-center gap-4 relative group">
	<button type="button" class="w-6 h-6 flex-shrink-0 border-2 rounded-full
																{isCompleted ? 'bg-sky-400 border-sky-800' : 'border-white/30'}"
					on:click={() => (isCompleted = !isCompleted)}></button>
	<p class="text-lg truncate {isCompleted ? 'line-through decoration-2' : null}">
		{ title }
	</p>
	<button type="button" class="p-1.5 bg-red-600 rounded absolute right-0 opacity-0 duration-500 group-hover:opacity-100"
			    on:click={() => dispatch("delete")}>
		<Icon icon="heroicons:trash-solid" class="text-white size-4" stroke="2" />
	</button>
</div>