import IonIcon from '@react-native-vector-icons/ionicons';
import {scale} from './SizeMatters';
import {Theme, useTheme} from './Theme';

import {useMemo, type ComponentProps} from 'react';

const IconSizeFactors = {
  small: 2,
  default: 2.5,
  medium: 3,
  large: 4,
};
type IconSizeFactorsKey = keyof typeof IconSizeFactors;
type ThemeColorsKey = keyof Theme['colors'];
type IonIconProps = ComponentProps<typeof IonIcon>;
type IconProps = Omit<IonIconProps, 'size' | 'color'> & {
  size?: IconSizeFactorsKey | IonIconProps['size'];
  color?: ThemeColorsKey | IonIconProps['color'];
};
const Icon = ({color = 'inherit', size = 'default', ...props}: IconProps) => {
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
  return <IonIcon color={color_} size={size_} {...props} />;
};

export default Icon;
