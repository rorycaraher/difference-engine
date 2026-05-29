import { Container } from "@cloudflare/containers";

export class Backend extends Container {
  defaultPort = 8000;
  sleepAfter = "2h";
}

export default {
  async fetch(request, env) {
    const id = env.BACKEND.idFromName("default");
    const stub = env.BACKEND.get(id);
    return stub.fetch(request);
  },
};
