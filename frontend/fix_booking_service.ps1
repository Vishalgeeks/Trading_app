$content = Get-Content "src\services\bookingService.js"
$new = @()
foreach ($line in $content) {
    if ($line -match "async getMyBookings") {
        $new += "  async createBooking(data) {"
        $new += "    const result = await api.post('/bookings', data);"
        $new += "    if (result.error) {"
        $new += "      return result;"
        $new += "    }"
        $new += "    return { error: false, data: mapBooking(result.data) };"
        $new += "  },"
        $new += ""
    }
    $new += $line
}
Set-Content "src\services\bookingService.js" $new
Write-Host "createBooking added"
