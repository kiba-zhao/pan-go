import type {DeviceQRCodeValue} from '@pango/data';
import {MarshalDeviceQRCode} from '@pango/data';
import type {ComponentProps, PropsWithChildren} from 'react';
import {useMemo} from 'react';
import QRCode from 'react-native-qrcode-svg';
import {Provider as I18NextProviderBase} from '../Common/I18Next';
import {scale} from '../Common/SizeMatters';
import {useTheme} from '../Common/Theme';
import {Namespace} from './ScreenRoute';

const QRCodeSizes = {
  default: scale(120),
  medium: scale(180),
  large: scale(280),
};
type QRCodeSizeKey = keyof typeof QRCodeSizes;

type QRCodeProps = ComponentProps<typeof QRCode>;
type DeviceQRCodeProps = DeviceQRCodeValue &
  Omit<QRCodeProps, 'value' | 'size'> & {
    disabled?: boolean;
    size?: QRCodeSizeKey | number;
  };
export const DeviceQRCode = ({
  name,
  peerId,
  disabled,
  size = 'default',
  ...props
}: DeviceQRCodeProps) => {
  const {colors} = useTheme();
  const value = useMemo(
    () =>
      name.length > 0 && peerId.length > 0 && !disabled
        ? MarshalDeviceQRCode({name, peerId})
        : void 0,
    [disabled, name, peerId],
  );

  const size_ = useMemo(() => {
    if (typeof size === 'string') {
      return QRCodeSizes[size as QRCodeSizeKey];
    }
    if (typeof size === 'number') {
      return scale(size);
    }
    return size;
  }, [size]);

  const color = useMemo(() => {
    if (disabled) return colors.textDisabled;
    return value === void 0 ? colors.error : colors.textPrimary;
  }, [disabled, value, colors]);

  return <QRCode size={size_} color={color} value={value} {...props} />;
};

export const I18NextProvider = ({children}: PropsWithChildren<{}>) =>
  I18NextProviderBase({children, defaultNS: Namespace});
