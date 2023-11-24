import { useDisclosure } from '@mantine/hooks';
import { AppShell, Burger, Group, Skeleton } from '@mantine/core';
// import { MantineLogo } from '@mantine/ds';
import Navbar from './Navbar';
import Chat from './Chat';

export default function AppRoot() {
  const [opened, { toggle }] = useDisclosure();

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 300, breakpoint: 'sm', collapsed: { mobile: !opened } }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          logo
          {/* <MantineLogo size={30} /> */}
        </Group>
      </AppShell.Header>
      <AppShell.Navbar p="md">
      <Navbar />
        {/* Navbar
        {Array(15)
          .fill(0)
          .map((_, index) => (
            <Skeleton key={index} h={28} mt="sm" animate={false} />
          ))} */}
      </AppShell.Navbar>
      <AppShell.Main
              style={{ 
                  // border: '1px solid #0f0',
                  // padding: '5rem', 
                  // width: '100vw', 
                  // height: '100vh' 
                  height: '100%' 
                  }}
      >
      <Chat />
      </AppShell.Main>
    </AppShell>
  );
}