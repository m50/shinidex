/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './**/*.html',
    './**/*.templ',
    './frontend/*.ts',
    './frontend/**/*.ts',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms')
  ],
}

