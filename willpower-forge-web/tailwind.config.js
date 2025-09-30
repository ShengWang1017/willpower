export default {
  content: ['index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        display: ['"Plus Jakarta Sans"', 'Inter', 'system-ui', 'sans-serif'],
        body: ['Inter', 'system-ui', 'sans-serif']
      },
      colors: {
        midnight: {
          50: '#f5f7fb',
          100: '#e8edf6',
          200: '#d4d9ec',
          300: '#b2bde0',
          400: '#8b9acb',
          500: '#6777b5',
          600: '#535f97',
          700: '#424d79',
          800: '#333d5f',
          900: '#282f4c'
        },
        moss: {
          50: '#f4faf7',
          100: '#e4f3ea',
          200: '#c7e6d4',
          300: '#a2d5b6',
          400: '#5cb986',
          500: '#379d64',
          600: '#2a7c4d',
          700: '#236243',
          800: '#1f5038',
          900: '#1a4331'
        }
      },
      boxShadow: {
        'soft-xl': '0 40px 80px -40px rgba(31, 41, 55, 0.25)',
        'inner-card': 'inset 0 1px 0 rgba(255, 255, 255, 0.35)'
      },
      backgroundImage: {
        'surface-gradient': 'linear-gradient(145deg, rgba(255,255,255,0.94), rgba(248,250,252,0.8))'
      }
    }
  },
  plugins: []
};
