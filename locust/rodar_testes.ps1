param(
    [int[]]$Usuarios = @(50, 100),
    [int]$RampUp = 10,
    [string]$Duracao = "60s"
)

Set-Location $PSScriptRoot
$LOCUST = "python -m locust"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host " TESTES DE CARGA - Locust" -ForegroundColor Cyan
Write-Host " Usuarios: $($Usuarios -join ', ') | Ramp-up: $RampUp/s | Duracao: $Duracao" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

function Executar-Teste {
    param(
        [string]$Nome,
        [string]$Arquivo,
        [string]$HostUrl,
        [string]$CsvBaseName,
        [string]$GrpcHost,
        [string]$GrpcPort,
        [int]$QuantidadeUsuarios
    )

    Write-Host "`n[+] Iniciando: $Nome ($QuantidadeUsuarios usuarios)" -ForegroundColor Yellow

    $CsvNameFinal = "${CsvBaseName}_${QuantidadeUsuarios}u"
    Remove-Item "${CsvNameFinal}_*.csv" -Force -ErrorAction SilentlyContinue

    if ($Nome -eq "gRPC Go" -and $QuantidadeUsuarios -le 50) {
        $env:LOCUST_WAIT_MIN = "4"
        $env:LOCUST_WAIT_MAX = "8"
    } elseif ($Nome -eq "gRPC Go") {
        $env:LOCUST_WAIT_MIN = "8"
        $env:LOCUST_WAIT_MAX = "12"
    } elseif ($QuantidadeUsuarios -le 50) {
        $env:LOCUST_WAIT_MIN = "1"
        $env:LOCUST_WAIT_MAX = "3"
    } else {
        $env:LOCUST_WAIT_MIN = "2"
        $env:LOCUST_WAIT_MAX = "5"
    }

    Write-Host "    -> Wait time: $env:LOCUST_WAIT_MIN a $env:LOCUST_WAIT_MAX segundos" -ForegroundColor DarkGray

    $cmd = "$LOCUST -f $Arquivo --headless -u $QuantidadeUsuarios -r $RampUp -t $Duracao --csv=$CsvNameFinal"

    if (![string]::IsNullOrEmpty($HostUrl)) {
        $cmd += " --host=$HostUrl"
    }

    if (![string]::IsNullOrEmpty($GrpcHost) -and ![string]::IsNullOrEmpty($GrpcPort)) {
        $env:GRPC_HOST = $GrpcHost
        $env:GRPC_PORT = $GrpcPort
        Write-Host "    -> Alvo gRPC: ${GrpcHost}:${GrpcPort}" -ForegroundColor DarkGray
    }

    Write-Host "    -> Comando: $cmd" -ForegroundColor DarkGray
    Invoke-Expression $cmd

    if (![string]::IsNullOrEmpty($GrpcHost)) {
        Remove-Item Env:\GRPC_HOST -ErrorAction SilentlyContinue
        Remove-Item Env:\GRPC_PORT -ErrorAction SilentlyContinue
    }

    Remove-Item Env:\LOCUST_WAIT_MIN -ErrorAction SilentlyContinue
    Remove-Item Env:\LOCUST_WAIT_MAX -ErrorAction SilentlyContinue

    if ($LASTEXITCODE -ne 0) {
        throw "Locust falhou em: $Nome"
    }

    $FailuresFile = "${CsvNameFinal}_failures.csv"
    if (Test-Path $FailuresFile) {
        $Failures = Import-Csv -Path $FailuresFile
        if (($Failures | Measure-Object).Count -gt 0) {
            Write-Host "`nFalhas encontradas em ${FailuresFile}:" -ForegroundColor Red
            $Failures | Format-Table -AutoSize
            throw "Corrija as falhas antes de gerar os graficos."
        }
    }
}

foreach ($QuantidadeUsuarios in $Usuarios) {
    Write-Host "`n########################################" -ForegroundColor Cyan
    Write-Host " Rodando carga com $QuantidadeUsuarios usuarios" -ForegroundColor Cyan
    Write-Host "########################################" -ForegroundColor Cyan

    Write-Host "`n=== GO ===" -ForegroundColor Green
    Executar-Teste "REST Go" "locustfile_rest.py" "http://localhost:8080" "resultados_rest_go_carga2" "" "" $QuantidadeUsuarios
    Executar-Teste "SOAP Go" "locustfile_soap_go.py" "http://localhost:8081" "resultados_soap_go_carga2" "" "" $QuantidadeUsuarios
    Executar-Teste "GraphQL Go" "locustfile_graphql.py" "http://localhost:8082" "resultados_graphql_go_carga2" "" "" $QuantidadeUsuarios
    Executar-Teste "gRPC Go" "locustfile_grpc.py" "" "resultados_grpc_go_carga2" "localhost" "8083" $QuantidadeUsuarios

    Write-Host "`n=== JAVA ===" -ForegroundColor Green
    Executar-Teste "REST Java" "locustfile_rest.py" "http://localhost:8090" "resultados_rest_java_carga2" "" "" $QuantidadeUsuarios
    Executar-Teste "SOAP Java" "locustfile_soap_java.py" "http://localhost:8091" "resultados_soap_java_carga2" "" "" $QuantidadeUsuarios
    Executar-Teste "GraphQL Java" "locustfile_graphql.py" "http://localhost:8092" "resultados_graphql_java_carga2" "" "" $QuantidadeUsuarios
    Executar-Teste "gRPC Java" "locustfile_grpc.py" "" "resultados_grpc_java_carga2" "localhost" "8093" $QuantidadeUsuarios
}

python .\gerar_graficos.py

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host " TODOS OS TESTES E GRAFICOS CONCLUIDOS!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
