# Simple HTTPS proxy for local dev
Docker image for running Traefik as an HTTPS proxy for up to three other containers by just providing a few environment variables for configuration

## What & Why
Sometimes you need HTTPS for local development, for example when implementing FIDO U2F 2-Step Verification / 
Two-Factor Authentication. The FIDO U2F spec strictly requires HTTPS in the browser. 

This container runs Traefik using a local configuration file. The entrypoint script updates the configuration based
on environment variables to keep it as simple as possible to use. 

## DNS Requirements
Let's Encrypt can either verify your SSL certificate request by making an HTTP call to your server or verifying a DNS record. Since we're talking about local development the HTTP challenge will not work, but DNS can so long as your
DNS is managed by a compatible provider. A list of compatible providers is available at https://docs.traefik.io/https/acme/#providers.

Create A records in your DNS  provider for the local development domains you want, something like dev1.domain.com or app1.domain.com and have them point to `127.0.0.1` so that when requested in your browser you'll just be directed to your own machine. When Traefik calls Let's Encrypt to provision certificates it will receive a challenge and create
a TXT record on your DNS provider for Let's Encrypt to verify before it receives your certificate. 

## Usage
In your `compose.yaml`, include a service like:

```yaml
  proxy:
    image: ghcr.io/sil-org/traefik-https-proxy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      #- ./traefik.toml:/etc/traefik/traefik.toml
      - ./cert/:/cert/
    env_file:
      - ./local.env
```

Copy the `local.env.example` file to `local.env` and update it with appropriate values. Using a separate file
for environment configuration helps prevent committing it to version control with credentials in it. 

Required env vars:
- `DNS_PROVIDER` - A valid value from https://docs.traefik.io/https/acme/#providers. Each provider will also required additional env vars for authentication. For example `cloudflare` requires either a `CLOUDFLARE_EMAIL` and `CLOUDFLARE_API_KEY` or just a `CLOUDFLARE_DNS_API_TOKEN`.
- `LETS_ENCRYPT_EMAIL` - An email address to use with Lets Encrypt, does not need to be previously "registered"
- `LETS_ENCRYPT_CA` - Either `staging` or `production`. Traefik does not appear to respect the staging caServer at the moment though.
- `TLD` - Used as the main domain on Lets Encrypt certificate, something like `domain.com`
- `SANS` - Comma separated list of domains to include on cert, something like `app1.domain.com,app2.domain.com`
- `BACKEND1_URL` - Url to backend #1, usually the name of the docker service in url form, example: `http://app1:80`
- `FRONTEND1_DOMAIN` - The domain name that should be routed to `BACKEND1_URL`, example: `app1.domain.com`

Optional env vars:
- `BACKEND2_URL` - If you need to route a second domain to a different container, define backend url here, example: `http://app2:80`
- `FRONTEND2_DOMAIN` - The domain name that should be routed to `BACKEND2_URL`, example: `app2.domain.com`
- `BACKEND3_URL` - If you need to route a third domain to a different container, define backend url here, example: `http://app3:80`
- `FRONTEND3_DOMAIN` - The domain name that should be routed to `BACKEND3_URL`, example: `app3.domain.com`

## Overriding `traefik.toml`
You'll notice in the `compose.yaml` example above a commented out volume for `traefik.toml`. If you 
don't want to use the simplified template that comes with this container and want to customize it, just provide 
your own config file and volume it in. The entrypoint script looks for specific placeholders and should not 
modify your own provided config. 

## License - MIT
MIT License

Copyright (c) 2021 SIL Global

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
