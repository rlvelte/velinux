#!/usr/bin/pwsh

# Git
$global:OnefetchDone = @{}
$global:OnefetchLastPwd = ''
$global:StarshipPrompt = ${function:prompt}

function global:prompt {
  if ($PWD.Path -ne $global:OnefetchLastPwd) {
    $global:OnefetchLastPwd = $PWD.Path

    $gitRoot = git rev-parse --show-toplevel 2>$null
    if ($gitRoot -and -not $global:OnefetchDone[$gitRoot]) {
      onefetch 2>$null
      $global:OnefetchDone[$gitRoot] = $true
    }
  }

  & $global:StarshipPrompt
}
