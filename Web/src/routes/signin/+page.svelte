<script lang="ts">
	import { currPage } from "$lib/store";
	import MetaTitle from "$lib/MetaTitle.svelte";
	import Icon from "@iconify/svelte";
	import { encodeAsFormData } from "$lib/functions";
	import { writable } from "svelte/store";

	currPage.set("signin");

	let isRequestPending = false;
	const formFeedback = {
		show: false,
		isError: false,
		message: "",
	};
	const formData = writable({
		username: "",
		password: "",
	});
	const formValidity = writable({
		username: false,
		password: false,
		__isDirty: false,
	});

	function validateInputFields() {
		formData.update((data) => {
			formValidity.set({
				username: new RegExp("^[A-Za-z0-9][A-Za-z0-9_]{2,14}[A-Za-z0-9]$").test(data.username),
				password: !(data.password.length < 8 || data.password.length > 60),
				__isDirty: true,
			});
			return data;
		});
	}

	async function submitLogin() {
		formFeedback.show = false;
		validateInputFields();

		const validity = $formValidity;

		if (!validity.username || !validity.password) {
			return;
		}

		try {
			isRequestPending = true;
			const res = await fetch("http://localhost:4555/v1/login", {
				method: "POST",
				body: encodeAsFormData($formData),
				credentials: "include",
				headers: {
					"Content-Type": "application/x-www-form-urlencoded",
				},
			});

			if (!res.ok) {
				if (res.status === 404) {
					formFeedback.show = true;
					formFeedback.isError = true;
					formFeedback.message = "Seems like the service is down. Please try again later.";
					isRequestPending = false;
					return;
				}

				formFeedback.show = true;
				formFeedback.isError = true;
				formFeedback.message = (await res.json())["error"];
				isRequestPending = false;
				return;
			}

			formFeedback.show = true;
			formFeedback.isError = false;
			formFeedback.message = "You've successfully signed in";

			window.location.href = "/";
		} catch (e) {
			console.error(e);
			formFeedback.show = true;
			formFeedback.isError = true;
			formFeedback.message = "Something wrong happened... try to refresh the page!";
			isRequestPending = false;
		}
	}
</script>

<MetaTitle title="Sign In into Svego List" />

<h1 class="text-3xl font-bold text-center">Sign In</h1>
<form class="w-[350px] mx-auto mt-10" on:submit|preventDefault={submitLogin}>
	<div>
		<label for="i_usr" class="block font-semibold">Username</label>
		<input
			id="i_usr"
			name="username"
			type="text"
			class="w-full text-black p-1 rounded border-2 {!$formValidity.username &&
			$formData.username.length > 0 &&
			$formValidity.__isDirty
				? 'border-red-500'
				: null}"
			required
			bind:value={$formData.username}
		/>
		{#if !$formValidity.username && $formData.username.length > 0 && $formValidity.__isDirty}
			<small class="text-xs text-red-400">Must be 4-16 characters, using letters, numbers, or underscores.</small>
		{/if}
	</div>
	<div class="mt-4">
		<label for="i_pw" class="block font-semibold">Password</label>
		<input
			id="i_pw"
			name="password"
			type="password"
			class="w-full text-black p-1 border-2 rounded {!$formValidity.password &&
			$formData.password.length > 0 &&
			$formValidity.__isDirty
				? 'border-red-500'
				: null}"
			required
			bind:value={$formData.password}
		/>
		{#if !$formValidity.password && $formData.password.length > 0 && $formValidity.__isDirty}
			<small class="text-xs text-red-400">Must be 8-60 characters.</small>
		{/if}
	</div>
	<button
		type="submit"
		class="group relative mt-8 py-1.5 rounded font-bold w-full flex gap-1 justify-center items-center
			   duration-300 {isRequestPending ? 'bg-gray-600' : 'bg-orange-700 hover:bg-amber-900'}"
		disabled={isRequestPending}
	>
		{#if !isRequestPending}
			<span class="block"> Sign In </span>
			<Icon
				icon="mingcute:arrow-right-fill"
				class="size-5 absolute duration-300 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100
						   right-10 group-hover:right-2 scale-y-[25%] group-hover:scale-y-100"
			/>
		{:else}
			<style>
				.bigger-loading path {
					stroke-width: 4;
				}
			</style>
			<span class="opacity-0">.</span>
			<Icon icon="line-md:loading-loop" class="size-5 bigger-loading" />
			<span class="opacity-0">.</span>
		{/if}
	</button>
</form>
{#if formFeedback.show}
	<div
		class="w-[350px] mx-auto p-2 border-2 mt-6 rounded text-center text-sm {formFeedback.isError
			? 'border-red-800 bg-red-300 text-red-950'
			: 'border-green-800 bg-green-300 text-green-950'}"
	>
		{formFeedback.message}
	</div>
{/if}
<hr class="w-[175px] mx-auto my-8 border rounded border-white/30" />
<p class="text-center text-sm">
	Don't have an account?
	<br />
	<a class="underline font-semibold text-sky-300" href="/signup">Sign up!</a>
</p>
