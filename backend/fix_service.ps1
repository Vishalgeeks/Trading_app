$content = Get-Content "internal\booking\service.go"
$new = @()
$i = 0
while ($i -lt $content.Length) {
    if ($content[$i] -match "type Service struct \{") {
        $new += "type Service struct {"
        $new += "`tbookingRepo      BookingRepository"
        $new += "`tavailabilityRepo AvailabilityRepository"
        $new += "`tnotificationSvc  notification.NotificationRepository"
        $new += "`tuserRepo         notification.UserRepository"
        $new += "`tdesignRepo       DesignGetter"
        $new += "}"
        # skip old lines until closing }
        while ($i -lt $content.Length -and $content[$i] -ne "}") { $i++ }
        $i++
        continue
    }
    if ($content[$i] -match "func NewService\(bookingRepo BookingRepository, availabilityRepo AvailabilityRepository, notificationSvc notification.NotificationRepository, userRepo notification.UserRepository\) \*Service \{") {
        $new += "func NewService(bookingRepo BookingRepository, availabilityRepo AvailabilityRepository, notificationSvc notification.NotificationRepository, userRepo notification.UserRepository, designRepo DesignGetter) *Service {"
        $new += "`treturn &Service{"
        $new += "`t`tbookingRepo:      bookingRepo,"
        $new += "`t`tavailabilityRepo: availabilityRepo,"
        $new += "`t`tnotificationSvc:  notificationSvc,"
        $new += "`t`tuserRepo:         userRepo,"
        $new += "`t`tdesignRepo:       designRepo,"
        $new += "`t}"
        $new += "}"
        # skip old lines until closing }
        while ($i -lt $content.Length -and $content[$i] -ne "}") { $i++ }
        $i++
        continue
    }
    $new += $content[$i]
    $i++
}
Set-Content "internal\booking\service.go" $new
Write-Host "Service struct and NewService updated"
