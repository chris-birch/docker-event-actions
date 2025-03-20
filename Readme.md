# Docker Event Actions

Using the excellent [yubiuser/docker-event-monitor](https://github.com/yubiuser/docker-event-monitor) project as a base, remove the notification feature, and add actions for use in my homelab.


## Features

- Auto update Sync DNS with records pointing to containers
- Auto update Cloudflare tunnels to allow internet access to containers

## Technical

The application uses Docker's API to connect to the [event stream](https://docs.docker.com/engine/config/v1.43/#tag/System/operation/SystemEvents). Each new event is procesed, logged and can be reported.

### Configuration

Configuration is loaded from a config file, by defaut `config.yml`. The path can be adjusted by `docker-event-monitor --config [path]`.
Currently the following options can be set via `config.yml`

```yaml
---
options:
  filter_strings: ["type=container"]
  exclude_strings: ["Action=exec_start", "Action=exec_die", "Action=exec_create"]
  log_level: debug
  server_tag: My Server

```

### Filter and exclude events

Docker Event Monitor offers two options that sound alike, but aren't: `Filter` and `Exclude`.
By default, the docker system event stream will contain **all** events. The `filter` option  is a docker built-in function that allows filtering certain events from the stream. It's a **positive** filter only, meaning it defines which events will pass the filter. The possible filters and syntax are described [here](https://docs.docker.com/engine/reference/commandline/events/#filter).

However, docker has no native support for **negative** filter (let all events pass, except those defined) - so I added it. To distingush it from postive filters, this option is named `exclude`.
Based on how it is implemented, **exclusion happens after filtering**. Together you can create configurations like filtering events of type container, but exclude reporting for a specific container or certain actions.

The syntax for exclusion is also `key=value`.  But as the exclusion happens on the data contained in the reported event, the `key`s are different from those used for `filtering`. E.g. instead of `event`, `Action` is used. To figure out which keys to use, it's best to enable debug logging and carefully inspect the event's data structure. A typical container event looks like

```
{Status:"start", ID:"b4a2a54c4487ddc0bbae006e48ae970d4b2fa4b9fd2bef390d8875cb6158c888", From:"squidfunk/mkdocs-material", Type:"container", Action:"start", Actor:events.Actor{ID:"b4a2a54c4487ddc0bbae006e48ae970d4b2fa4b9fd2bef390d8875cb6158c888", Attributes:map[string]string{"com.docker.compose.config-hash":"cd464ac038ddc9ee7a53599aaa9db6a85a01683a9a08a749582d0c0b8c0a595d", "com.docker.compose.container-number":"1", "com.docker.compose.depends_on":"", "com.docker.compose.image":"sha256:feb8ba83cb7272046551c69a58ec03ecda2306410a07844d22c166e810034aa6", "com.docker.compose.oneoff":"False", "com.docker.compose.project":"mkdocs-material", "com.docker.compose.project.config_files":"/home/pi/docker/mkdocs-material/docker-compose.yml", "com.docker.compose.project.working_dir":"/home/pi/docker/mkdocs-material", "com.docker.compose.service":"mkdocs-material", "com.docker.compose.version":"2.24.5", "image":"squidfunk/mkdocs-material", "name":"mkdocs-material", "org.opencontainers.image.created":"2024-02-10T06:18:18.743Z", "org.opencontainers.image.description":"Documentation that simply works", "org.opencontainers.image.licenses":"MIT", "org.opencontainers.image.revision":"a6286ef3ac3407e8b6c985cf0571fc0e2caa6f5b", "org.opencontainers.image.source":"https://github.com/squidfunk/mkdocs-material", "org.opencontainers.image.title":"mkdocs-material", "org.opencontainers.image.url":"https://github.com/squidfunk/mkdocs-material", "org.opencontainers.image.version":"9.5.9"}}, Scope:"local", Time:1708201602, TimeNano:1708201602805856956}
```

Keys of nested elements are joind by dots. E.g. `Actor.Attributes.com.docker.compose.project` or `Actor.Attributes.image`.
