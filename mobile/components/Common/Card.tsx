import type {TextProps, TextStyle, ViewProps, ViewStyle} from 'react-native';
import {StyleSheet, Text, View} from 'react-native';
import {OpacityPressable} from './Pressable';
import {scale, verticalScale} from './SizeMatters';
import {withTheme} from './StyleSheet';
import {Theme} from './Theme';

const styles = StyleSheet.create({
  cardPanel: {},
});

export const CardPanel = withTheme<ViewStyle, ViewProps, Theme>(
  View,
  ({colors, sizes}) => ({
    ...styles.cardPanel,
    paddingHorizontal: scale(sizes.base),
    paddingVertical: verticalScale(sizes.base),
    borderRadius: scale(sizes.base),
    backgroundColor: colors.surfaceVariant,
  }),
);

export const CardTitleText = withTheme<TextStyle, TextProps, Theme>(
  Text,
  ({colors, sizes, fonts}) => ({
    color: colors.onSurfaceVariant,
    fontFamily: fonts.medium.fontFamily,
    fontWeight: fonts.medium.fontWeight as TextStyle['fontWeight'],
    fontSize: scale(sizes.text),
  }),
);

export const CardText = withTheme<TextStyle, TextProps, Theme>(
  Text,
  ({colors, sizes, fonts}) => ({
    color: colors.onSurfaceVariant,
    fontFamily: fonts.regular.fontFamily,
    fontWeight: fonts.regular.fontWeight as TextStyle['fontWeight'],
    fontSize: scale(sizes.text * 0.8),
  }),
);

export const PressableCard = OpacityPressable;
