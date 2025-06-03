/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './**/*.html',
    './**/*.templ',
    './frontend/*.ts',
    './frontend/**/*.ts',
    '!./node_modules/**/*',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms')
  ],
}

