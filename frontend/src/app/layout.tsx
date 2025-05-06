import { Box, createTheme, MantineProvider } from "@mantine/core";
import "@mantine/core/styles.css";
import "@mantine/dates/styles.css";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import "./globals.css";

const theme = createTheme({
  primaryColor: 'dark',
  colors: {
    // override dark colors here to change them for all components
    dark: [
      '#d5d7e0',
      '#acaebf',
      '#8c8fa3',
      '#666980',
      '#4d4f66',
      '#34354a',
      '#2b2c3d',
      '#1d1e30',
      '#0c0d21',
      '#01010a',
    ],
  },
  fontFamily: "Instrument Sans, sans-serif",
  fontFamilyMonospace: "Azeret Mono, monospace",
  fontSizes: {
    xs: "10px",
    sm: "12px",
    md: "14px",
    lg: "16px",
    xl: "18px",
  },
});

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <MantineProvider
          theme={theme}
        >
          <Notifications />
          <Box bg="dark" h="100vh">
            {children}
          </Box>
        </MantineProvider>
      </body>
    </html>
  );
}
