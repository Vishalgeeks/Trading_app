$content = Get-Content "internal\booking\service.go"
$new = @()
$i = 0
while ($i -lt $content.Length) {
    if ($content[$i] -eq "}" -and $i + 1 -lt $content.Length -and $content[$i+1] -eq "" -and $i + 2 -lt $content.Length -and $content[$i+2] -eq "type Service struct {") {
        $new += "}"
        $new += ""
        $new += "type DesignGetter interface {"
        $new += "`tGetDesignByID(ctx context.Context, id string) (design.Design, error)"
        $new += "}"
        $new += ""
        $i += 2
    } else {
        $new += $content[$i]
    }
    $i++
}
Set-Content "internal\booking\service.go" $new
Write-Host "DesignGetter added"
