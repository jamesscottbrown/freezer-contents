/** @type {import('tailwindcss').Config}*/
const defaultTheme = require('tailwindcss/defaultTheme');
const ldnColors = require("@ldn-viz/themes/colors.json");
const plugin = require('tailwindcss/plugin'); // eslint-disable-line
const svgToDataUri = require('mini-svg-data-uri'); // eslint-disable-line

const config = {
	content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/@ldn-viz/**/*.{html,js,svelte,ts}'],

	darkMode: 'class',
	theme: {
		container: {
			padding: {
				DEFAULT: '1rem',
				sm: '2rem',
				lg: '4rem',
				xl: '5rem',
				'2xl': '6rem'
			}
		},
		extend: {
			colors: ldnColors,
			fontFamily: {
				sans: ['"Roboto"', ...defaultTheme.fontFamily.sans]
			},
			typography: (theme) => ({
				DEFAULT: {
					css: {
						h1: { fontWeight: 'bold' },
						a: { color: theme('colors.core.blue.600') }
					}
				}
			})
		}
	},

  plugins: [
    require('@tailwindcss/typography'), // eslint-disable-line
    // provides form-element reset
    require('@tailwindcss/forms')({ strategy: 'class' }), // eslint-disable-line
    // LDN default form themeing check ui/lib/introduction.mdx for more
    plugin(function ({ addComponents }) {
      addComponents({
        '.form-input, .form-textarea, .form-select, .form-multiselect, .form-checkbox, .form-radio':
          {
            'border-color': ldnColors.core.grey[300],

            '&:disabled': {
              'background-color': ldnColors.core.grey[200]
            }
          },
        '.form-input, .form-textarea, .form-select, .form-multiselect': {
          '&:focus': {
            '--tw-ring-color': ldnColors.core.blue[600],
            'border-color': ldnColors.core.blue[600]
          }
        },
        '.form-checkbox, .form-radio': {
          color: ldnColors.core.blue[600],
          '&:focus': {
            '--tw-ring-color': ldnColors.core.blue[600]
          }
        },
        '.form-label': {
          color: ldnColors.core.grey[700]
        },
        '.form-select': {
          'background-image': `url("${svgToDataUri(
            `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="${ldnColors.core.grey[700]}" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 8l4 4 4-4"/></svg>`
          )}")`
        },
        '.dark': {
          '.form-input, .form-textarea, .form-select, .form-multiselect, .form-checkbox, .form-radio':
            {
              'border-color': ldnColors.core.grey[400],
              'background-color': ldnColors.core.grey[600],
              color: ldnColors.core.grey[50],
              '&:focus': {
                '--tw-ring-offset-color': ldnColors.core.grey[800]
              },
              '&:disabled, &:disabled::placeholder': {
                'background-color': ldnColors.core.grey[400],
                color: ldnColors.core.grey[600]
              }
            },
          '.form-input::placeholder, .form-textarea::placeholder': {
            color: ldnColors.core.grey[200]
          },
          '.form-checkbox, .form-radio': {
            color: ldnColors.core.blue[600]
          },
          '.form-checkbox:checked, .form-radio:checked': {
            'background-color': ldnColors.core.blue[600]
          },
          '.form-label': { color: ldnColors.core.grey[50] },
          '.form-select': {
            'background-image': `url("${svgToDataUri(
              `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="${ldnColors.core.grey[100]}" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 8l4 4 4-4"/></svg>`
            )}")`
          }
        }
      });
    })
  ]


};

module.exports = config;