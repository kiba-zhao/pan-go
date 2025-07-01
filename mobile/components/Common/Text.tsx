import type {TextProps, TextStyle} from 'react-native';
import {Text as NativeText} from 'react-native';
import {scale} from './SizeMatters';
import {withTheme} from './StyleSheet';
import {Theme} from './Theme';

export const Text = withTheme<TextStyle, TextProps, Theme>(
  NativeText,
  ({colors, fonts, sizes}) => ({
    color: colors.onBackground,
    fontSize: scale(sizes.text),
    fontFamily: fonts.regular.fontFamily,
    fontWeight: fonts.regular.fontWeight as TextStyle['fontWeight'],
  }),
);
