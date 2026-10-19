#            _ _                                    _     
# __   _____| (_)_ __  _   ___  ___ ____      _____| |__  
# \ \ / / _ \ | | '_ \| | | \ \/ / '_ \ \ /\ / / __| '_ \ 
#  \ V /  __/ | | | | | |_| |>  <| |_) \ V  V /\__ \ | | |
#   \_/ \___|_|_|_| |_|\__,_/_/\_\ .__/ \_/\_/ |___/_| |_|
#                                |_|                      
#                                           pwsh by rvelte

# Interactivity
if ([Console]::IsInputRedirected -or [Console]::IsOutputRedirected) { return }

# Distrobox compatibility
if ((Test-Path '/run/.containerenv') -or $env:CONTAINER_ID_ANY) {
    $env:IN_DISTROBOX = '1'

    function global:prompt { "distrobox as $env:USER> " }
    function global:ll { ls -lah @args }
    function global:c { cat @args }

    return
}

# Pwsh
$PwshDotDir = "$HOME/.config/powershell"

# Source
Invoke-Expression (&starship init powershell)

# Configurations
Get-ChildItem "$PwshDotDir/source.d/*.ps1" | ForEach-Object { . $_ }

# Plugins
if (Get-Module -ListAvailable -Name PSFzf) {
    Set-PsFzfOption -PSReadlineChordProvider 'Ctrl+t' -PSReadlineChordReverseHistory 'Ctrl+r'
}
