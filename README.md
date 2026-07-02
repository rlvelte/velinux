# VeLinux - Opinionated openSUSE Configuration

My very own opinionated openSUSE dotfiles featuring a curated selection of packages,
system configurations, and the `vlx` utility suite for theme management and more.

> This is not a plug-and-play setup! Some manual linking and configuration is still required.

## Repository Structure

The repository mirrors the target filesystem for clarity:

```
velinux/
├── config/              
│   ├── environment.d/   # Environment variables
│   ├── eza/             # Fancy ls
│   ├── git/             # Git configuration
│   ├── hypr/            # Hyprland window manager
│   ├── kitty/           # Terminal emulator
│   ├── mako/            # Notification daemon
│   ├── mango/           # Mangowm window manager
│   ├── quickshell/      # Desktop shell
│   ├── sway/            # Sway window manager
│   ├── systemd/		 # User scoped systemd services
│   ├── vlx/
│   │   ├── bundles/     # Bundle definitions
│   │   ├── themes/      # Theme profiles and wallpapers
│   │   └── bundesliga/  # Bundesliga config
│   └── zsh/             # Shell configuration
│
├── etc/                 
│   ├── greetd/          # Graphical greeter
│   ├── keyd/            # Key remapping
│   └── systemd/		 # System scoped systemd services
│
├── vlx/                 
│   ├── internal/
│   │   ├── app/
│   │   │   ├── bundesliga/ # Bundesliga match tracker
│   │   │   ├── bundle/     # Bundle installer
│   │   │   ├── package/    # Package wrapper
│   │   │   └── themes/     # Theme management
│   │   └── core/			# Utilities
│   ├── main.go
│   └── Makefile
│
├── LICENSE.md
└── README.md
```

---

Maintained with ❤️ (and maybe a little laziness)