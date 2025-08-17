import FontAwesome6 from '@react-native-vector-icons/fontawesome6';
import IonIcon from '@react-native-vector-icons/ionicons';
import MaterialDesignIcon from '@react-native-vector-icons/material-design-icons';
import {scale} from './SizeMatters';
import {Theme, useTheme} from './Theme';

import {useMemo, type ComponentProps} from 'react';

const IconSizeFactors = {
  xsmall: 1.5,
  small: 2,
  default: 2.5,
  medium: 3,
  large: 4,
};
type IconSizeFactorsKey = keyof typeof IconSizeFactors;
type ThemeColorsKey = keyof Theme['colors'];
type IonIconProps = ComponentProps<typeof IonIcon>;
export type IconProps = Omit<IonIconProps, 'size' | 'color'> & {
  size?: IconSizeFactorsKey | IonIconProps['size'];
  color?: ThemeColorsKey | IonIconProps['color'];
};
const Icon = ({
  color = 'inherit',
  size = 'default',
  children,
  ...props
}: IconProps) => {
  const {sizes, colors} = useTheme();

  const color_ = useMemo(
    () =>
      typeof color === 'string' && colors[color as ThemeColorsKey] !== void 0
        ? colors[color as ThemeColorsKey]
        : color,
    [color, colors],
  );

  const size_ = useMemo(() => {
    if (typeof size === 'number') {
      return scale(size);
    }
    const factor = IconSizeFactors[size];
    return scale(sizes.base) * factor;
  }, [size, sizes]);
  return (
    <IonIcon color={color_} size={size_} {...props}>
      {children}
    </IonIcon>
  );
};

export default Icon;

type FontAwesome6Props = ComponentProps<typeof FontAwesome6>;
export type FAIconProps = Omit<FontAwesome6Props, 'size' | 'color'> & {
  size?: IconSizeFactorsKey | FontAwesome6Props['size'];
  color?: ThemeColorsKey | FontAwesome6Props['color'];
};

export const FAIcon = ({
  color = 'inherit',
  size = 'default',
  children,
  ...props
}: FAIconProps) => {
  const {sizes, colors} = useTheme();
  const color_ = useMemo(
    () =>
      typeof color === 'string' && colors[color as ThemeColorsKey] !== void 0
        ? colors[color as ThemeColorsKey]
        : color,
    [color, colors],
  );
  const size_ = useMemo(() => {
    if (typeof size === 'number') {
      return scale(size);
    }
    const factor = IconSizeFactors[size];
    return scale(sizes.base) * factor;
  }, [size, sizes]);
  return (
    <FontAwesome6 color={color_} size={size_} {...(props as FontAwesome6Props)}>
      {children}
    </FontAwesome6>
  );
};

type MaterialDesignIconProps = ComponentProps<typeof MaterialDesignIcon>;
export type MaterialIconProps = Omit<
  MaterialDesignIconProps,
  'size' | 'color'
> & {
  size?: IconSizeFactorsKey | MaterialDesignIconProps['size'];
  color?: ThemeColorsKey | MaterialDesignIconProps['color'];
};

export const MaterialIcon = ({
  color = 'inherit',
  size = 'default',
  children,
  ...props
}: MaterialIconProps) => {
  const {sizes, colors} = useTheme();
  const color_ = useMemo(
    () =>
      typeof color === 'string' && colors[color as ThemeColorsKey] !== void 0
        ? colors[color as ThemeColorsKey]
        : color,
    [color, colors],
  );
  const size_ = useMemo(() => {
    if (typeof size === 'number') {
      return scale(size);
    }
    const factor = IconSizeFactors[size];
    return scale(sizes.base) * factor;
  }, [size, sizes]);
  return (
    <MaterialDesignIcon
      color={color_}
      size={size_}
      {...(props as MaterialDesignIconProps)}>
      {children}
    </MaterialDesignIcon>
  );
};
