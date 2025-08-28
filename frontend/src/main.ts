import "./style.css"; // Keep your existing CSS file name
import App from "./App.svelte";

const app = new App({
    target: document.getElementById("app")!,
});

export default app;
