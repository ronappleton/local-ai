# LocalAI UI

This folder contains a minimal [Vue 3](https://vuejs.org/) project powered by [Vite](https://vitejs.dev/).

## Available Scripts

- `npm install` – install dependencies
- `npm run dev` – start the development server
- `npm run build` – build for production
- `npm run preview` – preview the production build

The dev server proxies requests starting with `/api` to `http://localhost:8080` so it can be served alongside the Go backend.

Styling is provided via [Tailwind CSS](https://tailwindcss.com/) which is loaded from the CDN in `index.html`. Feel free to use Tailwind utility classes when extending the UI.
