# Cloudflare DNS Updater
This little webserver allows updating the Cloudflare DNS, by supplying ip addresses. It is protected by a simple password.

A docker image is available at [ghcr.io/jonathanbout/cloudflare-dns-updater](https://github.com/JonathanBout/cloudflare-dns-updater/pkgs/container/cloudflare-dns-updater).

# Usage
run the docker image: `docker run -p 5000:5000 ghcr.io/jonathanbout/cloudflare-dns-updater:latest` you can replace :latest with a version of your choice.

## Environment Variables
* = Required

| Name | Default Value | Description | Allowed Values |
| --- | --- | --- | --- |
| `CLOUDFLARE_API_TOKEN` * | None | The API token for Cloudflare | A Valid Cloudflare API Token. |
| `CLOUDFLARE_ZONE` * | None | The Cloudflare zone to update. | A domain name you manage through cloudflare, and is accessible with the provided CLOUDFLARE_API_TOKEN. |
| `CLOUDFLARE_RECORDS` | The value of CLOUDFLARE_ZONE | The records to update in Cloudflare. | A subdomain of the CLOUDFLARE_ZONE domain, including the domain itself.  |
| `DYNDNS_LISTEN_ADDR` | 5000 | What address the server should listen on. Not recommended in Docker, please use port mappings. | An available port number between 1 and 65535 |
| `DYNDNS_USERNAME` * | None | The username required to be able to call this services API | Any string |
| `DYNDNS_PASSWORD` * | None | The password required to be able to call this services API | Any string |


If you find any problems or defects, feel free to create an issue or even a pull request!
