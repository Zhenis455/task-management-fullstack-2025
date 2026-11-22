$ErrorActionPreference = 'Stop'
$baseUrl = 'http://localhost:8082'
try {
    $health = Invoke-RestMethod "$baseUrl/health"
    if (-not $health.ok) { throw 'Health check returned false' }
    Write-Host 'PASS: /health' -ForegroundColor Green
    $payload = @{ title='Smoke test task'; description='Temporary test'; status='todo'; priority='high' } | ConvertTo-Json
    $created = Invoke-RestMethod "$baseUrl/api/tasks" -Method Post -ContentType 'application/json' -Body $payload
    if (-not $created.id) { throw 'Create response has no ID' }
    Write-Host 'PASS: create task' -ForegroundColor Green
    $updatedPayload = @{ title='Smoke test updated'; description='Temporary test'; status='done'; priority='low' } | ConvertTo-Json
    $updated = Invoke-RestMethod "$baseUrl/api/tasks/$($created.id)" -Method Put -ContentType 'application/json' -Body $updatedPayload
    if (-not $updated.done) { throw 'Update did not set done' }
    Write-Host 'PASS: update task' -ForegroundColor Green
    $list = Invoke-RestMethod "$baseUrl/api/tasks"
    if (-not (@($list) | Where-Object { $_.id -eq $created.id })) { throw 'Task missing in list' }
    Write-Host 'PASS: list tasks' -ForegroundColor Green
    Invoke-RestMethod "$baseUrl/api/tasks/$($created.id)" -Method Delete | Out-Null
    Write-Host 'PASS: delete task' -ForegroundColor Green
    exit 0
} catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
