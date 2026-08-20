$content = Get-Content "src\pages\Booking.jsx"
$new = @()
$i = 0
while ($i -lt $content.Length) {
    if ($content[$i] -match "if \(bookingResult\?\.type === 'success'\) \{" ) {
        $new += "  if (bookingResult?.type === 'success' && bookingResult.data) {"
        $i++
        continue
    }
    if ($content[$i] -match "const booking = bookingResult.data;" ) {
        $new += "  const booking = bookingResult.data;"
        $i++
        continue
    }
    $new += $content[$i]
    $i++
}
Set-Content "src\pages\Booking.jsx" $new
Write-Host "Booking.jsx success guard added"
