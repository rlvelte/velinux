# VeLinux – Opinionated openSUSE Configuration

My very own opinionated openSUSE dotfiles featuring a curated selection of packages,
system configurations, and the `vlx` utility suite for theme management and more.

> [!IMPORTANT]  
> This is not a plug-and-play setup! 
> Some manual linking and configuration is still required. Most configurations also work on other Linux distributions, but I couldn't be bothered to validate this.


## Repository Structure

The repository mirrors the target filesystem for clarity:

```
velinux/
├── config/              
│   ├── environment.d/   # Environment variables
│   ├── eza/             # Fancy ls
│   ├── git/             # Git configuration
│   ├── kitty/           # Terminal emulator
│   ├── mako/            # Notification daemon
│   ├── quickshell/      # Desktop shell
│   ├── sway/            # Sway window manager
│   ├── systemd/         # User scoped systemd services
│   ├── vlx/             # Configuration files for vlx
│   └── zsh/             # Shell configuration
│
├── etc/                 
│   ├── greetd/          # Graphical greeter
│   ├── keyd/            # Key remapping
│   └── systemd/		 # System scoped systemd services
│
├── vlx/                 
│   ├── ...
│   └── Makefile		 # Build and install for vlx
│
├── LICENSE.md
└── README.md
```

---

Maintained with ❤️ (and maybe a little laziness)