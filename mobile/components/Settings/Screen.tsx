import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation';
import {withLazy} from '../Common/Component';
import {addLanguageSource} from '../Common/I18Next';
import {
  NameEditorScreenName,
  Namespace,
  PeerIDScreenName,
  PeerPortEditorScreenName,
  QRCodeScreenName,
} from './ScreenRoute';

addLanguageSource('zh-CN', Namespace, () =>
  import(`../../locales/zh-CN/settings.json`).then(m => m.default),
);
addLanguageSource('en', Namespace, () =>
  import(`../../locales/en/settings.json`).then(m => m.default),
);

const SettingsAppScreen = () => (
  <Fragment>
    {QRCodeAppScreen()}
    {NameEditorAppScreen()}
    {PeerIDAppScreen()}
    {PeerPorEditorAppScreen()}
  </Fragment>
);

export default SettingsAppScreen;

const QRCodeScreen = lazy(() => import('./QRCodeScreen'));

const QRCodeAppScreen = () => (
  <Screen
    name={QRCodeScreenName}
    component={QRCodeScreen}
    options={{title: ''}}
  />
);

const NameEditorScreen = lazy(() => import('./NameEditorScreen'));
const NameEditorAppScreen = () => (
  <Screen
    name={NameEditorScreenName}
    component={NameEditorScreen}
    options={{title: ''}}
  />
);

const PeerIDScreen = lazy(() => import('./PeerIDScreen'));
const PeerIDHeaderActions = lazy(() =>
  import('./PeerIDScreen').then(m => ({default: m.PeerIDHeaderActions})),
);
const PeerIDAppScreen = () => (
  <Screen
    name={PeerIDScreenName}
    component={PeerIDScreen}
    options={{title: '', headerRight: PeerIDHeaderRight}}
  />
);

const PeerIDHeaderRight = withLazy(PeerIDHeaderActions);

const PeerPortEditorScreen = lazy(() => import('./PeerPortEditorScreen'));
const PeerPorEditorAppScreen = () => (
  <Screen
    name={PeerPortEditorScreenName}
    component={PeerPortEditorScreen}
    options={{title: ''}}
  />
);
