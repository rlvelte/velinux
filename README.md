# VeLinux - Opinionated openSUSE Configuration
My very own opinionated openSUSE dotfiles featuring a curated selection of packages, system configurations, and the `vlx` utility suite for theme management and some other tasks.

> This is not a plug-and-play setup! Some manual linking and configuration is still required.

## Repository Structure
The repository layout mirrors the target filesystem for clarity:

```
velinux/
├── config/          
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
│   ├── systemd/
│   │   └── user/
│   │       ├── vlx-bl-tracker.service  # Bundesliga poll oneshot
│   │       └── vlx-bl-tracker.timer    # 2-minute trigger
│   └── vlx/         
│       ├── bundles/ 		# Bundle definitions
│       ├── themes/  		# Theme profiles and wallpapers
│       └── bundesliga/  	# Bundesliga config
│
├── etc/             
│   ├── greetd/      # Graphical greeter
│   └── systemd/
│       └── system/
│           ├── zypper-refresh.service  # Zypper refresh oneshot
│           └── zypper-refresh.timer    # Daily trigger
│
├── vlx/             
│   ├── internal/
│   │   ├── app/
│   │   │   ├── bundesliga/  # Bundesliga match tracker
│   │   │   ├── bundle/      # Bundle management
│   │   │   ├── pkg/         # Package management
│   │   │   └── themes/      # Theme management
│   │   └── core/
│   │       ├── fsys/     # Filesystem utilities
│   │       ├── guard/    # Precondition checks
│   │       ├── http/     # HTTP client
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