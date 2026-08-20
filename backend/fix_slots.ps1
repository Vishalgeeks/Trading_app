$content = Get-Content "internal\booking\service.go"
$new = @()
$i = 0
while ($i -lt $content.Length) {
    if ($content[$i] -eq "`tvar allSlots []Slot" -and $i + 1 -lt $content.Length -and $content[$i+1] -eq "" -and $i + 2 -lt $content.Length -and $content[$i+2] -match "for _, av := range availabilities") {
        $new += "`tdesign, err := s.designRepo.GetDesignByID(ctx, designID)"
        $new += "`tif err != nil {"
        $new += "`t`treturn nil, fmt.Errorf(`"failed to load design: %v`", err)"
        $new += "`t}"
        $new += "`tduration := time.Duration(design.DurationMinutes) * time.Minute"
        $new += "`tif duration <= 0 {"
        $new += "`t`tduration = 2 * time.Hour"
        $new += "`t}"
        $new += ""
        $new += "`tvar allSlots []Slot"
        $new += ""
        $new += "`tfor _, av := range availabilities {"
        $new += "`t`tslots := generateSlotsForAvailability(av, designID, bookingDateParsed, s.bookingRepo, duration)"
        $new += "`t`tallSlots = append(allSlots, slots...)"
        $new += "`t}"
        $new += ""
        $new += "`treturn allSlots, nil"
        $new += "}"
        $new += ""
        $new += "func generateSlotsForAvailability(av availability.Availability, designID string, bookingDate time.Time, bookingRepo BookingRepository, duration time.Duration) []Slot {"
        # skip old lines until closing } of old function
        while ($i -lt $content.Length -and $content[$i] -ne "}") { $i++ }
        $i++
        continue
    }
    if ($content[$i] -match "for current.Add\(2\*time.Hour\)") {
        $new += "for current.Add(duration).Before(endTime) || current.Add(duration).Equal(endTime) {"
        $i++
        continue
    }
    if ($content[$i] -match "slotEnd := current.Add\(2 \* time.Hour\)") {
        $new += "`t`tslotEnd := current.Add(duration)"
        $i++
        continue
    }
    $new += $content[$i]
    $i++
}
Set-Content "internal\booking\service.go" $new
Write-Host "Slots fixed"
