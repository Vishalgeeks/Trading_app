$content = Get-Content "src\services\bookingService.js"
$new = @()
$i = 0
while ($i -lt $content.Length) {
    if ($content[$i] -match "async createBooking\(data\) \{") {
        $new += "  async createBooking(data) {"
        $new += "    const payload = {"
        $new += "      design_id: data.designId || data.design_id,"
        $new += "      booking_date: data.bookingDate || data.booking_date,"
        $new += "      start_time: data.startTime || data.start_time,"
        $new += "      notes: data.notes || '',"
        $new += "    };"
        $new += "    const result = await api.post('/bookings', payload);"
        $new += "    if (result.error) {"
        $new += "      return result;"
        $new += "    }"
        $new += "    return { error: false, data: mapBooking(result.data) };"
        $new += "  },"
        $i += 2
        continue
    }
    $new += $content[$i]
    $i++
}
Set-Content "src\services\bookingService.js" $new
Write-Host "createBooking fixed with snake_case"
