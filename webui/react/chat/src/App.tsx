import React from 'react';
import logo from './logo.svg';
import './App.css';
import ChatJoin from './Join';
import Chat from './Chat';
import { useMantineColorScheme, Button, Group } from '@mantine/core';
import { IconSun, IconMoon } from '@tabler/icons-react';

function App() {

  return (
    <div className="App">
      {/* <ChatJoin /> */}
      <Chat />
    </div>
  );
}

export default App;
