/** @type {import('tailwindcss').Config} */
const tailwindConfig = {
  content: [
    './src/**/*.{html,svelte,js,ts}', // Adjust based on your file structure
  ],
  darkMode: 'class', // Enable class-based dark mode
  theme: {
    extend: {},
  },
  plugins: [],
};

export default tailwindConfig;
