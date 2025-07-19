# Teste de Cenários de Falha - Rinha de Backend 2025
# Simula instabilidades nos processadores para testar fallback

$BaseUrl = "http://localhost:9999"
$DefaultProcessor = "http://localhost:8001"
$FallbackProcessor = "http://localhost:8002"
$Headers = @{"Content-Type" = "application/json"; "X-Rinha-Token" = "123"}

Write-Host "=== TESTE DE CENÁRIOS DE FALHA ===" -ForegroundColor Cyan
Write-Host ""

function Test-Payment {
    param([string]$TestName)

    $correlationId = [System.Guid]::NewGuid().ToString()
    $paymentData = @{
        correlationId = $correlationId
        amount = 100.00
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Uri "$BaseUrl/payments" -Method POST -Body $paymentData -Headers @{"Content-Type" = "application/json"}
        Write-Host "✅ $TestName - Pagamento processado: $correlationId" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "❌ $TestName - Falhou: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

function Get-Summary {
    try {
        $summary = Invoke-RestMethod -Uri "$BaseUrl/payments-summary" -Method GET
        Write-Host "   Default: $($summary.default.totalRequests) req, R$ $($summary.default.totalAmount)" -ForegroundColor Gray
        Write-Host "   Fallback: $($summary.fallback.totalRequests) req, R$ $($summary.fallback.totalAmount)" -ForegroundColor Gray
        return $summary
    } catch {
        Write-Host "❌ Erro ao obter summary: $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

# 1. CENÁRIO NORMAL (baseline)
Write-Host "1. Cenário Normal (baseline)..." -ForegroundColor Yellow
Test-Payment "Baseline"
Get-Summary
Write-Host ""

# 2. SIMULAR DELAY NO DEFAULT
Write-Host "2. Simulando delay no processador Default (2s)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "$DefaultProcessor/admin/configurations/delay" -Method PUT -Body '{"delay": 2000}' -Headers $Headers
    Write-Host "   Delay configurado no Default" -ForegroundColor Gray

    $startTime = Get-Date
    Test-Payment "Com Delay Default"
    $endTime = Get-Date
    $duration = ($endTime - $startTime).TotalMilliseconds
    Write-Host "   Tempo de resposta: $([math]::Round($duration))ms" -ForegroundColor Gray

    Get-Summary
} catch {
    Write-Host "❌ Erro ao configurar delay: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 3. SIMULAR FALHA NO DEFAULT (deve usar fallback)
Write-Host "3. Simulando falha no Default (deve usar Fallback)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "$DefaultProcessor/admin/configurations/failure" -Method PUT -Body '{"failure": true}' -Headers $Headers
    Write-Host "   Falha configurada no Default" -ForegroundColor Gray

    Test-Payment "Default com Falha"
    Get-Summary
} catch {
    Write-Host "❌ Erro ao configurar falha: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 4. SIMULAR FALHA EM AMBOS
Write-Host "4. Simulando falha em AMBOS processadores..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "$FallbackProcessor/admin/configurations/failure" -Method PUT -Body '{"failure": true}' -Headers $Headers
    Write-Host "   Falha configurada no Fallback também" -ForegroundColor Gray

    Test-Payment "Ambos com Falha"
    Get-Summary
} catch {
    Write-Host "❌ Erro ao configurar falha no fallback: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 5. RESTAURAR CONFIGURAÇÕES
Write-Host "5. Restaurando configurações normais..." -ForegroundColor Yellow
try {
    # Remover delay do default
    Invoke-RestMethod -Uri "$DefaultProcessor/admin/configurations/delay" -Method PUT -Body '{"delay": 0}' -Headers $Headers
    # Remover falha do default
    Invoke-RestMethod -Uri "$DefaultProcessor/admin/configurations/failure" -Method PUT -Body '{"failure": false}' -Headers $Headers
    # Remover falha do fallback
    Invoke-RestMethod -Uri "$FallbackProcessor/admin/configurations/failure" -Method PUT -Body '{"failure": false}' -Headers $Headers

    Write-Host "   Configurações restauradas" -ForegroundColor Green

    Test-Payment "Após Restauração"
    Get-Summary
} catch {
    Write-Host "❌ Erro ao restaurar configurações: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 6. TESTE DE PERFORMANCE FINAL
Write-Host "6. Teste de Performance (20 pagamentos rápidos)..." -ForegroundColor Yellow
$successCount = 0
$startTime = Get-Date

for ($i = 1; $i -le 20; $i++) {
    if (Test-Payment "Perf $i") {
        $successCount++
    }
    Start-Sleep -Milliseconds 50
}

$endTime = Get-Date
$totalTime = ($endTime - $startTime).TotalMilliseconds
$avgTime = $totalTime / 20

Write-Host "   Resultado: $successCount/20 pagamentos" -ForegroundColor Cyan
Write-Host "   Tempo total: $([math]::Round($totalTime))ms" -ForegroundColor Cyan
Write-Host "   Tempo médio: $([math]::Round($avgTime))ms" -ForegroundColor Cyan
Write-Host ""

Write-Host "=== SUMMARY FINAL ===" -ForegroundColor Cyan
Get-Summary

Write-Host ""
Write-Host "🏆 SISTEMA PRONTO PARA A RINHA DE BACKEND 2025!" -ForegroundColor Green