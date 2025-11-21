
import { useEffect } from 'react';
import RootLayout from './components/RootLayout';
import { BROADCAST_CHANNEL_NAME } from './constants';
import { useBrowserSandboxSync } from './hooks/useBrowserStateSync';
import { usePlaygroundStore } from './stores/playgroundStore';

function App() {
    const initPlayground = usePlaygroundStore((state) => state.initPlayground);

    useBrowserSandboxSync({ channelName: BROADCAST_CHANNEL_NAME });

    useEffect(() => {
        void initPlayground();
    }, [initPlayground])

    return <RootLayout />;
}

export default App
