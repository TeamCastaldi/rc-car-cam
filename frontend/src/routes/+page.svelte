<script lang="ts">
	import { PUBLIC_API_BASE_URL, PUBLIC_STREAM_AUTH_TOKEN } from '$env/static/public';

	type ConnectionState = 'loading' | 'connected' | 'error';
	let connectionState: ConnectionState = $state('loading');

	// Bumped to force the <img> to open a fresh connection after recovering
	// from an error — reusing the exact same src string wouldn't reliably
	// trigger a new request in every browser.
	let reloadNonce = $state(0);

	const HEALTH_CHECK_INTERVAL_MS = 3000;

	$effect(() => {
		const interval = setInterval(async () => {
			try {
				await fetch(`${PUBLIC_API_BASE_URL}/healthz`, { mode: 'no-cors' });
				if (connectionState === 'error') {
					reloadNonce++;
					connectionState = 'loading';
				}
			} catch {
				connectionState = 'error';
			}
		}, HEALTH_CHECK_INTERVAL_MS);

		return () => clearInterval(interval);
	});
</script>

<h1>RC Car Cam</h1>

{#if connectionState === 'loading'}
	<p>Connecting to stream…</p>
{:else if connectionState === 'error'}
	<p>Unable to reach the camera stream.</p>
{/if}

<img
	src="{PUBLIC_API_BASE_URL}/stream?token={encodeURIComponent(
		PUBLIC_STREAM_AUTH_TOKEN
	)}&r={reloadNonce}"
	alt="Live RC car camera feed"
	onload={() => (connectionState = 'connected')}
	onerror={() => (connectionState = 'error')}
	style={connectionState === 'error' ? 'display: none' : ''}
/>
