$content = Get-Content "internal\booking\service.go"
$new = @()
foreach ($line in $content) {
    if ($line -match '"mehndi-booking-backend/internal/notification"') {
        $new += '"mehndi-booking-backend/internal/design"'
        $new += "`t`"mehndi-booking-backend/internal/notification`""
    } else {
        $new += $line
    }
}
Set-Content "internal\booking\service.go" $new
Write-Host "Import added"
