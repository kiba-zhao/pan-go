import {Platform} from 'react-native';

/**
 * Default colors for light theme
 *
 *
 * @see {@link https://zenoo.github.io/mui-theme-creator}
 */
const DEFAULT_LIGHT_COLORS = {
  background: 'rgba(246, 247, 248, 0.51)',
  surface: '#fff',
  textPrimary: 'rgba(0, 0, 0, 0.87)',
  textSecondary: 'rgba(0, 0, 0, 0.6)',
  textDisabled: 'rgba(0, 0, 0, 0.38)',
  primary: '#3f51b5',
  primaryLight: 'rgb(101, 115, 195)',
  primaryDark: 'rgb(44, 56, 126)',
  primaryContrast: '#fff',
  secondary: '#f50057',
  secondaryLight: 'rgb(247, 51, 120)',
  secondaryDark: 'rgb(171, 0, 60)',
  secondaryContrast: '#fff',
  error: '#d32f2f',
  errorLight: '#ef5350',
  errorDark: '#c62828',
  errorContrast: '#fff',
  warning: '#ed6c02',
  warningLight: '#ff9800',
  warningDark: '#e65100',
  warningContrast: '#fff',
  info: '#0288d1',
  infoLight: '#03a9f4',
  infoDark: '#01579b',
  infoContrast: '#fff',
  success: '#2e7d32',
  successLight: '#4caf50',
  successDark: '#1b5e20',
  successContrast: '#fff',
  divider: 'rgba(0, 0, 0, 0.12)',
};
export const lightColors = Platform.select({
  android: {
    ...DEFAULT_LIGHT_COLORS,
  },
  ios: {
    ...DEFAULT_LIGHT_COLORS,
  },
  default: {
    ...DEFAULT_LIGHT_COLORS,
  },
});
export default lightColors;

const DEFAULT_DARK_COLORS = {
  ...DEFAULT_LIGHT_COLORS,
  background: '#121212',
  surface: 'rgba(255, 255, 255, 0.069)',
  textPrimary: '#fff',
  textSecondary: 'rgba(255, 255, 255, 0.7)',
  textDisabled: 'rgba(255, 255, 255, 0.5)',
  divider: 'rgba(255, 255, 255, 0.12)',
};
export const darkColors = Platform.select({
  ios: {
    ...DEFAULT_DARK_COLORS,
  },
  android: {
    ...DEFAULT_DARK_COLORS,
  },
  default: {
    ...DEFAULT_DARK_COLORS,
  },
});
