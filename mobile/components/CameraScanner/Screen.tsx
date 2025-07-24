import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation.tsx';

const CameraScannerAppScreen = () => (
  <Fragment>{QRScannerAppScreen()}</Fragment>
);

export default CameraScannerAppScreen;

const QRScannerScreen = lazy(() => import('./QRScannerScreen'));
export const QRScannerScreenName = 'cameraScanner.qr';

const QRScannerAppScreen = () => (
  <Screen
    name={QRScannerScreenName}
    component={QRScannerScreen}
    options={{
      title: 'QR Scanner',
    }}
  />
);
