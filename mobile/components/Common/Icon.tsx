import IonIcon from '@react-native-vector-icons/ionicons';
import {scale} from './SizeMatters';
import {Theme, withTheme} from './Theme';

import type {ComponentProps} from 'react';

export default IonIcon;

type IconProps = ComponentProps<typeof IonIcon>;
export const HeaderActionIcon = withTheme<IconProps, IconProps, Theme>(
  IonIcon,
  ({colors, sizes}, customProps) => ({
    color: colors.onBackground,
    size: scale(sizes.h4),
    ...customProps,
  }),
);
