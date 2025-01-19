import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

//读取上级目录中的config.yaml
import { get } from "es-toolkit/compat";
import fs from "fs";
import yaml from "js-yaml";
import path from "path";

const config: any = yaml.load(
  fs.readFileSync(path.resolve(__dirname, "../config.yaml"), "utf-8")
);

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: get(config, "frontend.port", 21422),
    host: get(config, "frontend.host", "0.0.0.0"),
  },
});
