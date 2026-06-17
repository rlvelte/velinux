# VeLinux - Opinionated openSUSE Configuration
My very own opinionated openSUSE dotfiles featuring a curated selection of packages, system configurations, and the `vlx` utility suite for theme management and some other tasks.

> This is not a plug-and-play setup! Some manual linking and configuration is still required.

## Repository Structure
The repository layout mirrors the target filesystem for clarity:

```
velinux/
├── config/          → ~/.config/
│   ├── hypr/        # Hyprland window manager
│   ├── sway/        # Sway window manager
│   ├── waybar/      # Waybar status bar
│   ├── rofi/        # Application launcher
│   ├── kitty/       # Terminal emulator
│   ├── mako/        # Notification daemon
│   ├── zsh/         # Shell configuration
│   ├── eza/         # Fancy ls
│   ├── git/         # Git configuration
│   ├── keyd/        # Key remapping
│   └── vlx/         # vlx utility config
│       ├── bundles/ # Bundle definitions
│       └── themes/  # Theme profiles and wallpapers
│
├── etc/             → /etc/
│   ├── greetd/      # Graphical greeter
│   └── systemd/
│       └── system/
│           ├── zypper-refresh.service  # Zypper refresh oneshot
│           └── zypper-refresh.timer    # Daily trigger
│
├── vlx/             → /usr/local/bin/
│   ├── internal/
│   │   ├── app/
│   │   │   ├── pkg/      # Package management
│   │   │   ├── bundle/   # Bundle management
│   │   │   └── themes/   # Theme management
│   │   └── core/
│   │       ├── fsys/     # Filesystem utilities
│   │       ├── guard/    # Precondition checks
│   │       ├── notify/   # Desktop notifications
│   │       ├── picker/   # Interactive selection
│   │       └── printer/  # Terminal output
│   ├── main.go
│   └── Makefile
│
└── README.md
```

---

Maintained with ❤️ (and maybe a little laziness)