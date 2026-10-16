// Code for backend proxy in Cloudflare Workers
// Doesn't automatically get updated, do changes manually in the Cloudflare Workers dashboard

const CLOUD_RUN_URL = "https://booklet-backend-258487485284.australia-southeast1.run.app";

export default {
  async fetch(request) {
    const url = new URL(request.url);
    url.hostname = new URL(CLOUD_RUN_URL).hostname;

    const proxiedRequest = new Request(url, {
      method: request.method,
      headers: request.headers,
      body: request.body,
      redirect: "follow",
    });

    return fetch(proxiedRequest);
  },
};
