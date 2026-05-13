// Styles
import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";
import colors from "vuetify/util/colors";

// Composable
import { createVuetify } from "vuetify";

export default createVuetify({
  defaults: {
    VRow: {
      gap: '0.5em'
    },
    VSheet: {
      elevation: 2,
    },
    VSwitch: {
      color: "primary",
    },
    VToolbar: {
      color: "primary",
    },
    VFooter: {
      color: "primary",
    },
    VTextarea: {
      color: "primary",
    },
    VAppBar: {
      VAppBarNavIcon: {
        elevation: 0,
      },
      VBtn: {
        elevation: 2,
        color: "white",
      },
    },
    VBtn: {
      elevation: 2,
      color: "primary",
    },
  },
  theme: {
    defaultTheme: "system",
    themes: {
      light: {
        dark: false,
        colors: {
          primary: colors.green.base,
          clickable: colors.yellow.base,
          noclick: colors.green.darken2,
          red: "E30000",
          yellow: "D9E400",
          purple: "660078",
          black: "070D0A",
        },
      },
      dark: {
        dark: true,
        colors: {
          primary: colors.green.darken2,
          clickable: colors.yellow.darken1,
          noclick: colors.green.darken4,
          red: "E30000",
          yellow: "D9E400",
          purple: "660078",
          black: "070D0A",
        },
      },
    },
  },
})
