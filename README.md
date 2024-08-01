# Simple Container Runtime in Go

This project implements a basic container runtime in Go, demonstrating core container concepts such as namespaces, chroot, and cgroups. It's designed for educational purposes to help understand how containers work at a low level.

## Features

- **Namespace Isolation**: Utilizes UTS, PID, and Mount namespaces for process isolation.
- **Filesystem Isolation**: Implements chroot to provide a separate root filesystem for the container.
- **Resource Limitation**: Uses cgroups to limit container resources (currently limits the number of processes).
- **Simple CLI**: Provides a straightforward command-line interface to run commands in containers.

## Prerequisites

- Go 1.16 or higher
- Linux operating system (tested on Ubuntu 20.04)
- Root privileges (required for creating namespaces and cgroups)

## Building the Project

1. Clone the repository:

   ```
   git clone https://github.com/senbo1/go-container
   cd go-container
   ```

2. Build the project:
   ```
   go build -o container main.go
   ```

## Usage

To run a command in a container:

```
sudo ./container run <command> [args...]
```

For example:

```
sudo ./container run /bin/bash
```

This will start a new bash session inside a container.

## Project Structure

- `main.go`: The main source file containing the container runtime implementation.
- `README.md`: This file, providing project information and usage instructions.

## How It Works

1. The `run` command creates a new process with isolated namespaces.
2. The child process sets up the container environment:
   - Sets a new hostname
   - Changes the root directory (chroot)
   - Sets up a new proc filesystem
   - Configures cgroups for resource limitation
3. The specified command is executed within this containerized environment.

## Limitations

This is a simplified container runtime for educational purposes. It lacks many features of production container runtimes:

- No image management
- Limited network isolation
- Basic resource constraints (only limits number of processes)
- No container lifecycle management

## Contributing

Contributions to improve the project or add new features are welcome! Please feel free to submit a pull request or open an issue for discussion.

## License

This project is open source and available under the [MIT License](LICENSE).
