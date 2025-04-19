# DayZ Server Status Landing

A simple static page for displaying DayZ server statuses via the `/info`
(`A2S_INFO`) endpoint from [dayz-exporter]

> 🔖 **Disclaimer**  
> Publishing JSON with `A2S_INFO` is not the primary purpose of
> [dayz-exporter], and status updates only occur when metrics are refreshed.
> For example, if the `scrape_interval` in Prometheus is set to `15s`, the
> `/info` data will update every 15 seconds. On one hand, this prevents
> excessive requests; on the other, it’s just an additional feature.
>
> If you don’t use metric collection and only need JSON server info, this
> solution is not for you—it would be overkill. Consider using a more
> specialized tool instead.

This document describes how to create a simple landing page with detailed
information about your servers.

![example]

## Exporter settings

* `/info` is disabled by default. To enable it, set the environment variable
  `DAYZ_EXPORTER_LISTEN_EXPOSE_INFO` or the YAML config parameter
  `listen.expose_info` to `true`.
* By default, if metric authentication is enabled, `/info` ignores it. To
  also password-protect `/info`, enable `DAYZ_EXPORTER_LISTEN_INFO_AUTH` or
  set `listen.info_auth` to `true`.
* If `/info` is served from a different domain or for local development, you
  can configure CORS via `DAYZ_EXPORTER_LISTEN_CORS_DOMAINS` (or
  `listen.cors_domains`). For example, `*` works, but in production,
  explicit domains are recommended.

## Status Page

The [index.html] file is straightforward. Make basic edits like the server
name and other details, but the key part is defining your servers in the
`SERVERS` constant:

```html
<script type="text/javascript">
  const SERVERS = [
    {
      name: "My Server 1",
      apiUrl: "http://127.0.0.1:8098/info", // direct
      ip: "public.server1.ip"
    },
    {
      name: "My Server 2",
      apiUrl: "/info/2", // via reverse proxy
      ip: "public.server2.ip"
    },
    {
      name: "My Server 3",
      apiUrl: "/info/3",
      ip: "public.server3.ip"
    }
  ];
</script>
```

Where:

* `name`: The display name if the server is unreachable.
* `apiUrl`: The `/info` endpoint URL for this server’s exporter.
* `ip`: The public IP shown in the description (port is taken from `/info`).

## Testing and Development

A disabled mock script is included for debugging. If you’re modifying the
main script or styles, uncomment it to test without connecting to a live
server—it generates random data for 6 servers:

```html
<!-- For testing only -->
<script src="servers-mock-test.js"></script>
```

## Discord Invites

To simplify Discord invite link management, use [discord-invite]. It
generates unique, personalized server invites and avoids expired links.

## Deployment

> ℹ️ **Info**
> For public hosting, use a reverse proxy to serve static files and
> consolidate all exporter instances under a single endpoint while
> restricting access to only `/info`.

A detailed production example with **Nginx**, caching, and rate-limiting:

```nginx
http {
  # Variables for /info routing
  map $1 $dayz_exporter_port {
    default 8091;  # /info/ (default)
    1       8091;  # /info/1 → 1st server (port 8091)
    2       8092;  # /info/2 → 2nd server (port 8092)
    3       8093;  # /info/3 → 3rd server (port 8093)
  }

  # Cache settings
  proxy_cache_path /tmp/cache levels=1:2 keys_zone=dayz_exporter_cache:1m;

  # Rate limits
  limit_req_zone $request_uri zone=dayz_exporter_limit:10m rate=5r/s;

  server {
    listen 80;
    server_name your_domain.com;

    # Static files
    location / {
      alias /path/to/your/html/;

      try_files $uri $uri/ /index.html;
      location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
        expires 1M;
        access_log off;
        add_header Cache-Control "public";
      }
    }

    # Consolidated /info endpoints
    location ~ ^/info/([0-9]+)$ {
      proxy_pass http://127.0.0.1:$dayz_exporter_port/info;

      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

      # Cache rules
      proxy_buffering on;
      proxy_ignore_headers Expires Cache-Control X-Accel-Expires;
      proxy_ignore_headers Set-Cookie;
      proxy_cache_methods GET;
      proxy_cache dayz_exporter_cache;
      proxy_cache_valid 1m;

      # Rate limiting
      limit_req zone=dayz_exporter_limit burst=6 nodelay;
    }

    # Discord invite generator (optional)
    location ~ ^/(d|dsc|disc|discord|invite)$ {
      proxy_pass http://127.0.0.1:8080; # discord-invite port

      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

      limit_req zone=dayz_exporter_limit;
    }
  }
}
```

[dayz-exporter]: https://github.com/WoozyMasta/dayz-exporter
[discord-invite]: https://github.com/WoozyMasta/discord-invite
[example]: example.jpg
[index.html]: index.html
