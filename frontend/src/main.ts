import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";

import "./assets/styles/main.css";
import { globalComponents } from "./components/global/index";

const pinia = createPinia();
const app = createApp(App);

// Register global components
Object.keys(globalComponents).forEach((name) => {
  app.component(name, globalComponents[name]);
});

app.use(pinia);
app.use(router);

app.mount("#app");
