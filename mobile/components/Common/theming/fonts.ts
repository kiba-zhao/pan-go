import {Platform} from 'react-native';

/**
 * Default fonts
 * @see {@link https://github.com/react-native-elements/react-native-elements}
 * @see {@link https://github.com/react-native-elements/react-native-elements/blob/next/packages/base/src/helpers/fonts.tsx}
 */
const fonts = Platform.select({
  default: {
    regular: {
      fontFamily: 'sans-serif',
      fontWeight: 'normal',
    },
    medium: {
      fontFamily: 'sans-serif-medium',
      fontWeight: 'normal',
    },
    bold: {
      fontFamily: 'sans-serif',
      fontWeight: 'bold',
    },
    heavy: {
      fontFamily: 'sans-serif',
      fontWeight: 'heavy',
    },
  },
  ios: {
    regular: {
      fontFamily: 'System',
      fontWeight: '400',
    },
    medium: {
      fontFamily: 'System',
      fontWeight: '500',
    },
    bold: {
      fontFamily: 'System',
      fontWeight: '600',
    },
    heavy: {
      fontFamily: 'System',
      fontWeight: '700',
    },
  },
});

export default fonts;
