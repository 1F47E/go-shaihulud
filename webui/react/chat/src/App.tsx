import React from 'react';
import logo from './logo.svg';
import './App.css';
import ChatJoin from './Join';
import { useMantineColorScheme, Button, Group } from '@mantine/core';
import { IconSun, IconMoon } from '@tabler/icons-react';

function App() {

  return (
    <div className="App">
      <header className="App-header">
        <ChatJoin />
      </header>
    </div>
  );
}

export default App;
