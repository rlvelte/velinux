#!/usr/bin/pwsh

# History
$MaximumHistoryCount = 50000

Set-PSReadLineOption -HistorySavePath "$PwshDotDir/.hist_pwsh"
Set-PSReadLineOption -HistoryNoDuplicates
Set-PSReadLineOption -PredictionSource History
Set-PSReadLineOption -PredictionViewStyle InlineView
