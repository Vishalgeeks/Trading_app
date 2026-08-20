$content = Get-Content "internal\booking\service.go"
$new = @()
$i = 0
while ($i -lt $content.Length) {
    if ($content[$i] -eq "`tendTime := startTime.Add(2 * time.Hour)") {
        $new += "`tdesign, err := s.designRepo.GetDesignByID(ctx, req.DesignID)"
        $new += "`tif err != nil {"
        $new += "`t`treturn Booking{}, fmt.Errorf(`"failed to load design: %v`", err)"
        $new += "`t}"
        $new += ""
        $new += "`tduration := time.Duration(design.DurationMinutes) * time.Minute"
        $new += "`tif duration <= 0 {"
        $new += "`t`tduration = 2 * time.Hour"
        $new += "`t}"
        $new += ""
        $new += "`tendTime := startTime.Add(duration)"
        $i++
        continue
    }
    $new += $content[$i]
    $i++
}
Set-Content "internal\booking\service.go" $new
Write-Host "CreateBooking duration fixed"
