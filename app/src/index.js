import { DurableObject } from "cloudflare:workers";

export class Backend extends DurableObject {
  async fetch(request) {
    return await this.container.fetch(request);
  }
}

export default {
  async fetch(request, env) {
    const id = env.BACKEND.idFromName("singleton");
    const stub = env.BACKEND.get(id);
    return stub.fetch(request);
  },
};
