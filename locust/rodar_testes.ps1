param(
    [int]$Usuarios = 50,
    [int]$RampUp = 50,
    [string]$Duracao = "60s"
)

Set-Location $PSScriptRoot
$LOCUST = "python -m locust"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host " TESTES DE CARGA - Locust (Carga Fixa)" -ForegroundColor Cyan
Write-Host " Usuarios: $Usuarios | Ramp-up: $RampUp/s | Duracao: $Duracao" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

function Executar-Teste {
    param(
        [string]$Nome,
        [string]$Arquivo,
        [string]$HostUrl,
        [string]$CsvBaseName,
        [string]$GrpcHost,
        [string]$GrpcPort
    )

    Write-Host "`n[+] Iniciando: $Nome" -ForegroundColor Yellow

    # Adiciona o número de usuários no nome do arquivo para não sobrescrever
    $CsvNameFinal = "${CsvBaseName}_${Usuarios}u"
    $cmd = "$LOCUST -f $Arquivo --headless -u $Usuarios -r $RampUp -t $Duracao --csv=$CsvNameFinal"
    
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
}

Write-Host "`n=== GO ===" -ForegroundColor Green
Executar-Teste "REST Go" "locustfile_rest.py" "http://localhost:8080" "resultados_rest_go_carga2" "" ""
Executar-Teste "SOAP Go" "locustfile_soap.py" "http://localhost:8081" "resultados_soap_go_carga2" "" ""
Executar-Teste "GraphQL Go" "locustfile_graphql.py" "http://localhost:8082" "resultados_graphql_go_carga2" "" ""
Executar-Teste "gRPC Go" "locustfile_grpc.py" "" "resultados_grpc_go_carga2" "localhost" "8083"

Write-Host "`n=== JAVA ===" -ForegroundColor Green
Executar-Teste "REST Java" "locustfile_rest.py" "http://localhost:8090" "resultados_rest_java_carga2" "" ""
Executar-Teste "SOAP Java" "locustfile_soap.py" "http://localhost:8091" "resultados_soap_java_carga2" "" ""
Executar-Teste "GraphQL Java" "locustfile_graphql.py" "http://localhost:8092" "resultados_graphql_java_carga2" "" ""
Executar-Teste "gRPC Java" "locustfile_grpc.py" "" "resultados_grpc_java_carga2" "localhost" "8093"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host " TODOS OS TESTES CONCLUIDOS COM SUCESSO!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan