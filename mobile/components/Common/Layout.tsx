import {
  ComponentProps,
  PropsWithChildren,
  useMemo,
  type ComponentType,
} from 'react';
import type {ColorValue, FlexStyle, ViewProps} from 'react-native';
import {ScrollView, StyleSheet, View} from 'react-native';
import {scale, verticalScale} from './SizeMatters';
import type {Theme} from './Theme';
import {useTheme} from './Theme';

type ThemeColorsKey = keyof Theme['colors'];
type LayoutProps<BaseProps extends Pick<ViewProps, 'style'>> =
  PropsWithChildren<{
    padding?: number | [number, number] | [number, number, number, number];
    margin?: number | [number, number] | [number, number, number, number];
    color?: ThemeColorsKey | ColorValue;
    radius?: number;
    direction?: FlexStyle['flexDirection'];
    gap?: number;
  }> &
    BaseProps;

export function withLayout<BaseProps extends Pick<ViewProps, 'style'>>(
  BaseComponent: ComponentType<BaseProps>,
  defaultProps: Pick<
    LayoutProps<BaseProps>,
    'padding' | 'margin' | 'radius' | 'color' | 'direction' | 'gap' | 'style'
  >,
): ComponentType<LayoutProps<BaseProps>> {
  return ({
    padding = defaultProps.padding || 0,
    margin = defaultProps.margin || 0,
    radius = defaultProps.radius || 0,
    color = defaultProps.color,
    direction = defaultProps.direction,
    gap = defaultProps.gap,
    style,
    children,
    ...props
  }: LayoutProps<BaseProps>) => {
    const {sizes, colors} = useTheme();

    const color_ = useMemo(
      () =>
        typeof color === 'string' && colors[color as ThemeColorsKey] !== void 0
          ? colors[color as ThemeColorsKey]
          : color,
      [color, colors],
    );

    const [top, right, bottom, left] = useMemo(() => {
      if (typeof padding === 'number') {
        if (padding === 0) return [0, 0, 0, 0];
        return [padding, padding, padding, padding].map(_ => _ * sizes.base);
      }
      if (Array.isArray(padding)) {
        if (padding.length === 4) return padding.map(_ => _ * sizes.base);
        if (padding.length === 2)
          return [
            padding[0] * sizes.base,
            padding[1] * sizes.base,
            padding[0] * sizes.base,
            padding[1] * sizes.base,
          ];
      }
      return [0, 0, 0, 0];
    }, [padding, sizes]);

    const [mrTop, mrRight, mrBottom, mrLeft] = useMemo(() => {
      if (typeof margin === 'number') {
        if (margin === 0) return [0, 0, 0, 0];
        return [margin, margin, margin, margin].map(_ => _ * sizes.base);
      }
      if (Array.isArray(margin)) {
        if (margin.length === 4) return margin.map(_ => _ * sizes.base);
        if (margin.length === 2)
          return [
            margin[0] * sizes.base,
            margin[1] * sizes.base,
            margin[0] * sizes.base,
            margin[1] * sizes.base,
          ];
      }
      return [0, 0, 0, 0];
    }, [margin, sizes]);

    const borderRadius = useMemo(() => radius * sizes.base, [radius, sizes]);

    const gap_ = useMemo(() => {
      if (direction === void 0 || gap === void 0) return void 0;
      if (gap === 0) return 0;
      const value = gap * sizes.base;
      return direction.startsWith('row') ? scale(value) : verticalScale(value);
    }, [direction, gap, sizes]);

    return (
      <BaseComponent
        style={StyleSheet.compose(
          [
            defaultProps.style,
            {
              backgroundColor: color_,
              paddingTop: verticalScale(top),
              paddingRight: scale(right),
              paddingBottom: verticalScale(bottom),
              paddingLeft: scale(left),
              marginTop: verticalScale(mrTop),
              marginRight: scale(mrRight),
              marginBottom: verticalScale(mrBottom),
              marginLeft: scale(mrLeft),
              borderRadius: scale(borderRadius),
              flexDirection: direction,
              gap: gap_,
            },
          ],
          style,
        )}
        {...(props as BaseProps)}>
        {children}
      </BaseComponent>
    );
  };
}

type ScreenLayoutProps = ComponentProps<typeof ScrollView>;
export const ScreenLayout = ({
  children,
  style,
  ...props
}: ScreenLayoutProps) => {
  const {sizes} = useTheme();
  return (
    <ScrollView
      {...props}
      style={StyleSheet.compose({paddingHorizontal: scale(sizes.base)}, style)}>
      <View style={{height: verticalScale(sizes.base * 1.5)}} />
      {children}
    </ScrollView>
  );
};

export const RowLayout = withLayout(View, {
  style: {flexWrap: 'wrap'},
  direction: 'row',
  gap: 1,
});
