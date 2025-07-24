import {SafeAreaProvider} from 'react-native-safe-area-context';
import AppScreen from './components/App/Screen';
import {ReactQueryProvider} from './components/Common/ReactQuery';

const App = () => (
  <SafeAreaProvider>
    <ReactQueryProvider>
      <AppScreen />
    </ReactQueryProvider>
  </SafeAreaProvider>
);

export default App;
