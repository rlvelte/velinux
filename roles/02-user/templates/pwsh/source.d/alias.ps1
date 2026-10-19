#!/usr/bin/pwsh

# Cat/Bat
function global:c { cat @args }
function global:b { bat @args }

# Ll
function global:ll { ls -lah @args }

# Navigation
function global:.. { Set-Location .. }

# Fastfetch
function global:ff { fastfetch @args }
function global:of { onefetch @args }

# Opencode
function global:oc { opencode @args }

# Git
function global:ga { git add . }
function global:gs { git status }
function global:gp { git push }
function global:gco { git commit @args }

# GnuPG
function global:gpg-reset {
  gpgconf --kill gpg-agent
  gpg-connect-agent /bye
  ssh-add -l
}

# Systemctl
function global:shutdown { systemctl poweroff }
function global:reboot { systemctl reboot }
function global:bios { sudo systemctl reboot --firmware-setup }
