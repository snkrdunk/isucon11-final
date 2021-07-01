module.exports = {
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eef5f2',
          100: '#cef0ed',
          200: '#95e8d5',
          300: '#59d0a9',
          400: '#1eb379',
          500: '#19935c',
          600: '#13853b',
          700: '#136630',
          800: '#0f4626',
          900: '#0b2b1e',
        },
        secondary: {
          50: '#fdfcfb',
          100: '#fbf0f0',
          200: '#f6cfe0',
          300: '#eba2c0',
          400: '#e8739a',
          500: '#da4f7b',
          600: '#c2355a',
          700: '#9a2842',
          800: '#6f1c2b',
          900: '#431117',
        },
      },
      gridTemplateColumns: {
        calendar: 'minmax(0, 2.5rem) repeat(5, minmax(0, 1fr));',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
}
