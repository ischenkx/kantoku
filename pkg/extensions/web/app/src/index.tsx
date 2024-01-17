import React from "react";
import { createRoot } from "react-dom/client";

import App from "./App";
import * as qs from 'qs'

const qsParse = qs.parse

qs.parse = (str, opts) => {
    opts ||= {}
    opts.arrayLimit = 10000000
    return qsParse(str, opts)
}

const container = document.getElementById("root") as HTMLElement;
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
