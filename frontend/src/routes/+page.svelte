<script lang="ts">
	import { PUBLIC_API_BASE_URL } from '$env/static/public';

	type ConnectionState = 'loading' | 'connected' | 'error';
	let state: ConnectionState = $state('loading');

	const HEALTH_CHECK_INTERVAL_MS = 3000;

	$effect(() => {
		const interval = setInterval(async () => {
			if (state === 'error') return;
			try {
				await fetch(`${PUBLIC_API_BASE_URL}/healthz`, { mode: 'no-cors' });
			} catch {
				state = 'error';
			}
		}, HEALTH_CHECK_INTERVAL_MS);

		return () => clearInterval(interval);
	});
</script>

<h1>RC Car Cam</h1>

{#if state === 'loading'}
	<p>Connecting to stream…</p>
{:else if state === 'error'}
	<p>Unable to reach the camera stream.</p>
{/if}

<img
	src="{PUBLIC_API_BASE_URL}/stream"
	alt="Live RC car camera feed"
	onload={() => (state = 'connected')}
	onerror={() => (state = 'error')}
	style={state === 'error' ? 'display: none' : ''}
/>
