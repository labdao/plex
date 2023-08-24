/* Components */
import { Providers } from '@/lib/providers'
import { TopNav } from './components/TopNav/TopNav'
import { FootBar } from './components/Footbar/FootBar'
import { UserLoader } from './components/UserLoader/UserLoader'

import { Container, Grid, Box } from '@mui/material'



/* Instruments */
import './styles/globals.css'

export default function RootLayout(props: React.PropsWithChildren) {
  return (
    <Providers>
      <html lang="en">
      <body>
        <Box display="flex" flexDirection="column" height="100vh"> {/* Fill entire view height and set up for flex */}
          <TopNav />

          <Container maxWidth="lg">
            <Box mt={5} mb={5} flexGrow={1} display="flex" alignItems="center" justifyContent="center"> {/* Center content */}
              <Grid container direction="column" spacing={3}>
                <Grid item xs={12}>
                  <UserLoader>
                    <main>{props.children}</main>
                  </UserLoader>
                </Grid>
              </Grid>
            </Box>
          </Container>

          <FootBar />
        </Box>
      </body>
      </html>
    </Providers>
  )
}
