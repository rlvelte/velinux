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
│   └── keyd/        # Key remapping
│
├── etc/             → /etc/
│   └── greetd/      # Graphical greeter
│
├── vlx/             → /usr/local/bin/
│   ├── internal/
│   │   ├── app/
│   │   │   ├── pkg/      # Packages/Schemes
│   │   │   ├── stat/     # Statistics
│   │   │   └── themes/   # Themes
│   │   └── system/
│   ├── main.go
│   └── Makefile
│
└── README.md
```

---

Maintained with ❤️ (and maybe a little laziness)