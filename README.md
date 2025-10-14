# VictoriaTraces MCP Server

[![Latest Release](https://img.shields.io/github/v/release/VictoriaMetrics-Community/mcp-victoriatraces?sort=semver&label=&filter=!*-victoriatraces&logo=github&labelColor=gray&color=gray&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics-Community%2Fmcp-victoriatraces%2Freleases%2Flatest)](https://github.com/VictoriaMetrics-Community/mcp-victoriatraces/releases)
[![Trust Score](https://archestra.ai/mcp-catalog/api/badge/quality/VictoriaMetrics-Community/mcp-victoriatraces)](https://archestra.ai/mcp-catalog/victoriametrics-community__mcp-victoriatraces)
![License](https://img.shields.io/github/license/VictoriaMetrics-Community/mcp-victoriatraces?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics-Community%2Fmcp-victoriatraces%2Fblob%2Fmain%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

The implementation of [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for [VictoriaTraces](https://docs.victoriametrics.com/victoriatraces/).

This provides access to your VictoriaTraces instance and seamless integration with [VictoriaTraces APIs](https://docs.victoriametrics.com/victoriatraces/querying/#http-api) and [documentation](https://docs.victoriametrics.com/victoriatraces/).
It can give you a comprehensive interface for traces, observability, and debugging tasks related to your VictoriaTraces instances, enable advanced automation and interaction capabilities for engineers and tools.

## Features

This MCP server allows you to use almost all read-only APIs of VictoriaTraces:

- Get services and operations (span names)
- Query traces, explore and analyze traces data
 
In addition, the MCP server contains embedded up-to-date documentation and is able to search it without online access.

More details about the exact available tools and prompts can be found in the [Usage](#usage) section.

You can combine functionality of tools, docs search in your prompts and invent great usage scenarios for your VictoriaTraces instance.
And please note the fact that the quality of the MCP Server and its responses depends very much on the capabilities of your client and the quality of the model you are using.

You can also combine the MCP server with other observability or doc search MCP Servers and get even more powerful results.

## Requirements

- [VictoriaTraces](https://docs.victoriametrics.com/victoriatraces/) instance: ([single-node](https://docs.victoriametrics.com/victoriatraces/) or [cluster](https://docs.victoriametrics.com/victoriatraces/cluster/))
- Go 1.25 or higher (if you want to build from source)

## Installation

### Go

```bash
go install github.com/VictoriaMetrics-Community/mcp-victoriatraces/cmd/mcp-victoriatraces@latest
```

### Binaries

Just download the latest release from [Releases](https://github.com/VictoriaMetrics-Community/mcp-victoriatraces/releases) page and put it to your PATH.

Example for Linux x86_64 (note that other architectures and platforms are also available):

```bash
latest=$(curl -s https://api.github.com/repos/VictoriaMetrics-Community/mcp-victoriatraces/releases/latest | grep 'tag_name' | cut -d\" -f4)
wget https://github.com/VictoriaMetrics-Community/mcp-victoriatraces/releases/download/$latest/mcp-victoriatraces_Linux_x86_64.tar.gz
tar axvf mcp-victoriatraces_Linux_x86_64.tar.gz
```

### Docker

You can run VictoriaTraces MCP Server using Docker.

This is the easiest way to get started without needing to install Go or build from source.

```bash
docker run -d --name mcp-victoriatraces \
  -e VT_INSTANCE_ENTRYPOINT=https://localhost:10428 \
  -e MCP_SERVER_MODE=http \
  -e MCP_LISTEN_ADDR=:8081 \
  -p 8081:8081 \
  ghcr.io/victoriametrics-community/mcp-victoriatraces
```

You should replace environment variables with your own parameters.

Note that the `MCP_SERVER_MODE=http` flag is used to enable Streamable HTTP mode. 
More details about server modes can be found in the [Configuration](#configuration) section.

See available docker images in [github registry](https://github.com/orgs/VictoriaMetrics-Community/packages/container/package/mcp-victoriatraces).

Also see [Using Docker instead of binary](#using-docker-instead-of-binary) section for more details about using Docker with MCP server with clients in stdio mode.


### Source Code

For building binary from source code you can use the following approach:

- Clone repo:

  ```bash
  git clone https://github.com/VictoriaMetrics-Community/mcp-victoriatraces.git
  cd mcp-victoriatraces
  ```
- Build binary from cloned source code:

  ```bash
  make build
  # after that you can find binary mcp-victoriatraces and copy this file to your PATH or run inplace
  ```
- Build image from cloned source code:

  ```bash
  docker build -t mcp-victoriatraces .
  # after that you can use docker image mcp-victoriatraces for running or pushing
  ```

## Configuration

MCP Server for VictoriaTraces is configured via environment variables:

| Variable                   | Description                                             | Required | Default          | Allowed values         |
|----------------------------|---------------------------------------------------------|----------|------------------|------------------------|
| `VT_INSTANCE_ENTRYPOINT`   | URL to VictoriaTraces instance                            | Yes      | -                | -                      |
| `VT_INSTANCE_BEARER_TOKEN` | Authentication token for VictoriaTraces API               | No       | -                | -                      |
| `VT_INSTANCE_HEADERS`      | Custom HTTP headers to send with requests (comma-separated key=value pairs) | No       | -                | -                      |
| `MCP_SERVER_MODE`          | Server operation mode. See [Modes](#modes) for details. | No       | `stdio`          | `stdio`, `sse`, `http` |
| `MCP_LISTEN_ADDR`          | Address for SSE or HTTP server to listen on             | No       | `localhost:8081` | -                      |
| `MCP_DISABLED_TOOLS`       | Comma-separated list of tools to disable                | No       | -                | -                      |
| `MCP_HEARTBEAT_INTERVAL`   | Defines the heartbeat interval for the streamable-http protocol. <br /> It means the MCP server will send a heartbeat to the client through the GET connection, <br /> to keep the connection alive from being closed by the network infrastructure (e.g. gateways) | No | `30s`  | -                      |

### Modes

MCP Server supports the following modes of operation (transports):

- `stdio` - Standard input/output mode, where the server reads commands from standard input and writes responses to standard output. This is the default mode and is suitable for local servers.
- `sse` - Server-Sent Events. Server will expose the `/sse` and `/message` endpoints for SSE connections.
- `http` - Streamable HTTP. Server will expose the `/mcp` endpoint for HTTP connections.

More info about traqnsports you can find in MCP docs:

- [Core concepts -> Transports](https://modelcontextprotocol.io/docs/concepts/transports)
- [Specifications -> Transports](https://modelcontextprotocol.io/specification/2025-03-26/basic/transports)

### Сonfiguration examples

```bash
export VT_INSTANCE_ENTRYPOINT="https://localhost:10428"

# Custom headers for authentication (e.g., behind a reverse proxy)
# Expected syntax is key=value separated by commas
export VT_INSTANCE_HEADERS="<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"

# Server mode
export MCP_SERVER_MODE="sse"
export MCP_SSE_ADDR="0.0.0.0:8082"
```

## Endpoints

In SSE and HTTP modes the MCP server provides the following endpoints:

| Endpoint            | Description                                                                                      |
|---------------------|--------------------------------------------------------------------------------------------------|
| `/sse` + `/message` | Endpoints for messages in SSE mode (for MCP clients that support SSE)                            |
| `/mcp`              | HTTP endpoint for streaming messages in HTTP mode (for MCP clients that support Streamable HTTP) |
| `/metrics`          | Metrics in Prometheus format for monitoring the MCP server                                       |
| `/health/liveness`  | Liveness check endpoint to ensure the server is running                                          |
| `/health/readiness` | Readiness check endpoint to ensure the server is ready to accept requests                        |

## Setup in clients

### Cursor

Go to: `Settings` -> `Cursor Settings` -> `MCP` -> `Add new global MCP server` and paste the following configuration into your Cursor `~/.cursor/mcp.json` file:

```json
{
  "mcpServers": {
    "victoriatraces": {
      "command": "/path/to/mcp-victoriatraces",
      "env": {
        "VT_INSTANCE_ENTRYPOINT": "<YOUR_VT_INSTANCE>",
        "VT_INSTANCE_BEARER_TOKEN": "<YOUR_VT_BEARER_TOKEN>",
        "VT_INSTANCE_HEADERS": "<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
      }
    }
  }
}
```

See [Cursor MCP docs](https://docs.cursor.com/context/model-context-protocol) for more info.

### Claude Desktop

Add this to your Claude Desktop `claude_desktop_config.json` file (you can find it if open `Settings` -> `Developer` -> `Edit config`):

```json
{
  "mcpServers": {
    "victoriatraces": {
      "command": "/path/to/mcp-victoriatraces",
      "env": {
        "VT_INSTANCE_ENTRYPOINT": "<YOUR_VT_INSTANCE>",
        "VT_INSTANCE_BEARER_TOKEN": "<YOUR_VT_BEARER_TOKEN>",
        "VT_INSTANCE_HEADERS": "<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
      }
    }
  }
}
```

See [Claude Desktop MCP docs](https://modelcontextprotocol.io/quickstart/user) for more info.

### Claude Code

Run the command:

```sh
claude mcp add victoriatraces -- /path/to/mcp-victoriatraces \
  -e VT_INSTANCE_ENTRYPOINT=<YOUR_VT_INSTANCE> \
  -e VT_INSTANCE_BEARER_TOKEN=<YOUR_VT_BEARER_TOKEN> \
  -e VT_INSTANCE_HEADERS="<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
```

See [Claude Code MCP docs](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/tutorials#set-up-model-context-protocol-mcp) for more info.

### Visual Studio Code

Add this to your VS Code MCP config file:

```json
{
  "servers": {
    "victoriatraces": {
      "type": "stdio",
      "command": "/path/to/mcp-victoriatraces",
      "env": {
        "VT_INSTANCE_ENTRYPOINT": "<YOUR_VT_INSTANCE>",
        "VT_INSTANCE_BEARER_TOKEN": "<YOUR_VT_BEARER_TOKEN>",
        "VT_INSTANCE_HEADERS": "<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
      }
    }
  }
}
```

See [VS Code MCP docs](https://code.visualstudio.com/docs/copilot/chat/mcp-servers) for more info.

### Zed

Add the following to your Zed config file:

```json
  "context_servers": {
    "victoriatraces": {
      "command": {
        "path": "/path/to/mcp-victoriatraces",
        "args": [],
        "env": {
          "VT_INSTANCE_ENTRYPOINT": "<YOUR_VT_INSTANCE>",
          "VT_INSTANCE_BEARER_TOKEN": "<YOUR_VT_BEARER_TOKEN>",
          "VT_INSTANCE_HEADERS": "<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
        }
      },
      "settings": {}
    }
  }
}
```

See [Zed MCP docs](https://zed.dev/docs/ai/mcp) for more info.

### JetBrains IDEs

- Open `Settings` -> `Tools` -> `AI Assistant` -> `Model Context Protocol (MCP)`.
- Click `Add (+)`
- Select `As JSON`
- Put the following to the input field:

```json
{
  "mcpServers": {
    "victoriatraces": {
      "command": "/path/to/mcp-victoriatraces",
      "env": {
        "VT_INSTANCE_ENTRYPOINT": "<YOUR_VT_INSTANCE>",
        "VT_INSTANCE_BEARER_TOKEN": "<YOUR_VT_BEARER_TOKEN>",
        "VT_INSTANCE_HEADERS": "<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
      }
    }
  }
}
```

### Windsurf

Add the following to your Windsurf MCP config file.

```json
{
  "mcpServers": {
    "victoriatraces": {
      "command": "/path/to/mcp-victoriatraces",
      "env": {
        "VT_INSTANCE_ENTRYPOINT": "<YOUR_VT_INSTANCE>",
        "VT_INSTANCE_BEARER_TOKEN": "<YOUR_VT_BEARER_TOKEN>",
        "VT_INSTANCE_HEADERS": "<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
      }
    }
  }
}
```

See [Windsurf MCP docs](https://docs.windsurf.com/windsurf/mcp) for more info.

### Using Docker instead of binary

You can run VictoriaTraces MCP Server using Docker instead of local binary.

You should replace run command in configuration examples above in the following way:

```
{
  "mcpServers": {
    "victoriametrics": {
      "command": "docker",
        "args": [
          "run",
          "-i", "--rm",
          "-e", "VT_INSTANCE_ENTRYPOINT",
          "-e", "VT_INSTANCE_BEARER_TOKEN",
          "-e", "VT_INSTANCE_HEADERS",
          "ghcr.io/victoriametrics-community/mcp-victoriatraces",
        ],
      "env": {
        "VT_INSTANCE_ENTRYPOINT": "<YOUR_VT_INSTANCE>",
        "VT_INSTANCE_BEARER_TOKEN": "<YOUR_VT_BEARER_TOKEN>",
        "VT_INSTANCE_HEADERS": "<HEADER>=<HEADER_VALUE>,<HEADER>=<HEADER_VALUE>"
      }
    }
  }
}
```

## Usage

After [installing](#installation) and [configuring](#setup-in-clients) the MCP server, you can start using it with your favorite MCP client.

You can start dialog with AI assistant from the phrase:

```
Use MCP VictoriaTraces in the following answers
```

But it's not required, you can just start asking questions and the assistant will automatically use the tools and documentation to provide you with the best answers.

### Toolset

MCP VictoriaTraces provides numerous tools for interacting with your VictoriaTraces instance.

Here's a list of available tools:

| Tool            | Description                                      |
|-----------------|--------------------------------------------------|
| `documentation` | Search in embedded VictoriaTraces documentation  |
| `services`      | List of all traced services                      |
| `service_names` | Get all the span names (operations) of a service |
| `traces`        | Query traces                                     |
| `trace`         | Get trace info by trace ID                       |
| `dependencies`  | Query the service dependency graph               |

### Prompts

The server includes pre-defined prompts for common tasks.

These are just examples at the moment, the prompt library will be added to in the future:

| Prompt          | Description                                             |
|-----------------|---------------------------------------------------------|
| `documentation` | Search VictoriaTraces documentation for specific topics |

## Roadmap
 
- [ ] Implement multitenant version of MCP (that will support several deployments)
- [x] Add service graph tool after release of [this feature](https://github.com/VictoriaMetrics/VictoriaTraces/pull/52) (see [the PR](#7))

## Disclaimer

AI services and agents along with MCP servers like this cannot guarantee the accuracy, completeness and reliability of results.
You should double check the results obtained with AI.

The quality of the MCP Server and its responses depend very much on the capabilities of your client and the quality of the model you are using.

## Contributing

Contributions to the MCP VictoriaTraces project are welcome! 

Please feel free to submit issues, feature requests, or pull requests.
