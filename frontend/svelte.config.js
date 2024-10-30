///////////////////////////// Build to run on server side (SSR) serve using a host with a server, vite server is for development only ////////////////////////////////
import adapter from '@sveltejs/adapter-auto';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		// adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://kit.svelte.dev/docs/adapters for more information about adapters.
		adapter: adapter()
	}
};

export default config;

//////////////////////////////// Build to run static, no need for a server, run local, all compiled, vite server is also included in the build////////////////////////////
// import adapter from '@sveltejs/adapter-static';

// const config = {
//   kit: {
//     adapter: adapter(),
//     prerender: {
//       entries: ['*']  // Prerender all routes by default
//     }
//   }
// };

// export default config;

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Static build  																				Server Side Build
//	 1. Extremely fast due to compiled static files and CDN delivery.								1. Potentially slower, depends on server response
//	 2. Pre rendered during build time no need for a server to render.								2. Dynamically rendered on each request
//	 3. Simple, just deploy static files															3. More complex, required server setup etc...
//	 4. Cheaper to host because no server needed													4. More expensive due to server cost
//	 5. Good for buildErrorMessage, landing pageXOffset, documentation.								5. Good for dashboards, e-commerce, dynamic apps