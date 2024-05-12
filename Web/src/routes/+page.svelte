<script lang="ts">
	import Task from "../components/list/Task.svelte";
	import Icon from "@iconify/svelte";
	import { currPage } from "$lib/store";

	currPage.set("home");

	const FilterCompletionStatus = {
		ALL: "ALL",
		COMPLETED: "COMPLETED",
		NOT_COMPLETED: "NOT_COMPLETED",
	} as const;

	type FilterCompletionStatus = (typeof FilterCompletionStatus)[keyof typeof FilterCompletionStatus];

	let sampleTasks = [
		{
			id: 1,
			title: "One Simple Task",
			isCompleted: false,
		},
		{
			id: 2,
			title: "Another Simple Task",
			isCompleted: true,
		},
		{
			id: 3,
			title: "One Harder Task",
			isCompleted: false,
		},
		{
			id: 4,
			title: "Another Harder Task",
			isCompleted: true,
		},
		{
			id: 5,
			title: "I'm a task that is very big to fit inside the list and I hope I get ellipsis at the end",
			isCompleted: false,
		},
		{
			id: 6,
			title: "I'mATaskWithNoSpacesThatHopesToGetEllipsisAsWellAtTheEndOfTheTagBecauseItIsTooMuchToShow",
			isCompleted: false,
		},
	];
	let visibleTasks = sampleTasks;

	// Filter variables
	let fCompletion: FilterCompletionStatus = FilterCompletionStatus.ALL;

	// Form variables
	let newTaskTitle: string = "";

	// Functions
	function nextCompletionStatus(curr: FilterCompletionStatus): FilterCompletionStatus {
		let keys = Object.keys(FilterCompletionStatus);
		let nextValKey = keys.at((keys.indexOf(curr) + 1) % keys.length);

		return FilterCompletionStatus[nextValKey as keyof typeof FilterCompletionStatus];
	}

	function addNewItem() {
		sampleTasks.push({
			id: sampleTasks.length,
			title: newTaskTitle,
			isCompleted: false,
		});

		fCompletion = FilterCompletionStatus.ALL;
		sampleTasks = visibleTasks = sampleTasks;
	}

	function removeItem(rId: number) {
		sampleTasks = sampleTasks.filter((t) => t.id !== rId);
		visibleTasks = sampleTasks;
	}

	function filterCompleted() {
		fCompletion = nextCompletionStatus(fCompletion);

		switch (fCompletion) {
			case FilterCompletionStatus.ALL:
				visibleTasks = sampleTasks;
				break;
			case FilterCompletionStatus.COMPLETED:
				visibleTasks = sampleTasks.filter((t) => t.isCompleted);
				break;
			case FilterCompletionStatus.NOT_COMPLETED:
				visibleTasks = sampleTasks.filter((t) => !t.isCompleted);
				break;
			default:
				console.error(`Invalid completion status ${fCompletion}`);
		}
	}
</script>

<svelte:head>
	<title>SveGo - Your tasks</title>
</svelte:head>

<h1 class="text-3xl font-bold text-center">Tasks</h1>

<div class="w-[650px] mt-10 mx-auto rounded">
	<form on:submit={addNewItem}>
		<div class="flex gap-4">
			<input
				id="new_item_title"
				class="border-2 rounded p-2 flex-grow text-slate-950"
				type="text"
				placeholder="New task name"
				bind:value={newTaskTitle}
			/>
			<button type="submit" class="px-3 rounded bg-orange-600 duration-300 hover:bg-orange-800 hover:scale-110">
				<Icon icon="typcn:plus" class="mx-auto size-6" />
			</button>
		</div>
	</form>
	<div class="mt-6 mb-4 flex items-center justify-between border-b-2 border-white/30 pb-2">
		<span class="font-semibold text-lg">Filters</span>
		<button
			type="button"
			class="cursor-pointer bg-gray-700 px-4 py-1 text-center rounded text-sm font-semibold"
			on:click={filterCompleted}
		>
			{#if fCompletion === FilterCompletionStatus.ALL}
				All
			{:else if fCompletion === FilterCompletionStatus.COMPLETED}
				Completed
			{:else if fCompletion === FilterCompletionStatus.NOT_COMPLETED}
				To do
			{:else}
				ERROR
			{/if}
			({visibleTasks.length})
		</button>
	</div>
	<div class="space-y-6">
		{#each visibleTasks.toReversed() as t}
			<Task title={t.title} isCompleted={t.isCompleted} on:delete={() => removeItem(t.id)} />
		{/each}
	</div>
</div>
