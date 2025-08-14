import type {PropsWithChildren} from 'react';
import {Provider as I18NextProviderBase} from '../Common/I18Next';
import {Namespace} from './ScreenRoute';

export const I18NextProvider = ({children}: PropsWithChildren<{}>) =>
  I18NextProviderBase({children, defaultNS: Namespace});
