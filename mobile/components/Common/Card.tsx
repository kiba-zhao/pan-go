import type {ViewProps, ViewStyle} from 'react-native';
import {StyleSheet, View} from 'react-native';
import {OpacityPressable} from './Pressable';
import {scale, verticalScale} from './SizeMatters';
import {withTheme} from './StyleSheet';
import {Theme} from './Theme';

const styles = StyleSheet.create({
  card: {
    boxShadow: [
      {
        color: 'rgba(0,0,0, .2)',
        offsetX: 0,
        offsetY: 2,
        blurRadius: 1,
        spreadDistance: -1,
      },
      {
        color: 'rgba(0,0,0, .14)',
        offsetX: 0,
        offsetY: 1,
        blurRadius: 1,
        spreadDistance: 0,
      },
      {
        color: 'rgba(0,0,0, .12)',
        offsetX: 0,
        offsetY: 1,
        blurRadius: 3,
        spreadDistance: 0,
      },
    ],
  },
});

export const Card = withTheme<ViewStyle, ViewProps, Theme>(
  View,
  ({sizes, colors}) => ({
    ...styles.card,
    paddingHorizontal: scale(sizes.base),
    paddingVertical: verticalScale(sizes.base),
    borderRadius: scale(sizes.base),
    backgroundColor: colors.surface,
  }),
);

export const PressableCard = OpacityPressable;
