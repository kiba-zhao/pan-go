import {HeaderButton} from '@react-navigation/elements';
import type {ViewProps, ViewStyle} from 'react-native';
import {View} from 'react-native';
import {scale} from 'react-native-size-matters';
import type {Theme} from '../Common/Theme';

import {withTheme} from '../Common/StyleSheet';

export const HeaderAction = HeaderButton;

export const HeaderActions = withTheme<ViewStyle, ViewProps, Theme>(
  View,
  ({sizes}) => ({
    flexDirection: 'row',
    alignItems: 'flex-end',
    gap: scale(sizes.base),
    paddingRight: scale(sizes.base),
  }),
);
