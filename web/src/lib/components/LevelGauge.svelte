<script lang="ts">
	let { pctRemaining }: { pctRemaining: number } = $props();

	let fill = $derived(Math.max(0, Math.min(100, pctRemaining)));
	let colorClass = $derived(
		fill > 50 ? 'bg-green-500' : fill > 20 ? 'bg-amber-400' : 'bg-red-500'
	);
</script>

<div class="relative h-64 w-12 overflow-hidden rounded-xl bg-gray-700">
	<div
		class="absolute bottom-0 left-0 right-0 transition-all duration-1000 ease-out {colorClass}"
		style="height: {fill}%"
	></div>
	<div class="absolute inset-0 flex items-end justify-center pb-2">
		<span class="text-xs font-bold text-white drop-shadow">{fill.toFixed(0)}%</span>
	</div>
</div>
