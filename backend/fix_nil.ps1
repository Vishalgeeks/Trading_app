$content = Get-Content "internal\booking\service.go"
$new = @()
$i = 0
while ($i -lt $content.Length) {
    if ($content[$i] -eq "`tdesign, err := s.designRepo.GetDesignByID(ctx, req.DesignID)") {
        $new += "`tduration := 2 * time.Hour"
        $new += "`tif s.designRepo != nil {"
        $new += "`t`tdesign, err := s.designRepo.GetDesignByID(ctx, req.DesignID)"
        $new += "`t`tif err == nil {"
        $new += "`t`t`tduration = time.Duration(design.DurationMinutes) * time.Minute"
        $new += "`t`t`tif duration <= 0 {"
        $new += "`t`t`t`tduration = 2 * time.Hour"
        $new += "`t`t`t}"
        $new += "`t`t}"
        $new += "`t}"
        $i += 9
        continue
    }
    if ($content[$i] -eq "`tdesign, err := s.designRepo.GetDesignByID(ctx, designID)") {
        $new += "`tduration := 2 * time.Hour"
        $new += "`tif s.designRepo != nil {"
        $new += "`t`tdesign, err := s.designRepo.GetDesignByID(ctx, designID)"
        $new += "`t`tif err == nil {"
        $new += "`t`t`tduration = time.Duration(design.DurationMinutes) * time.Minute"
        $new += "`t`t`tif duration <= 0 {"
        $new += "`t`t`t`tduration = 2 * time.Hour"
        $new += "`t`t`t}"
        $new += "`t`t}"
        $new += "`t}"
        $i += 9
        continue
    }
    $new += $content[$i]
    $i++
}
Set-Content "internal\booking\service.go" $new
Write-Host "Nil-safe design lookup added"
