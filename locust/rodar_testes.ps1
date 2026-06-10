param(
    [int]$Usuarios = 50,
    [int]$RampUp = 5,
    [string]$Duracao = "60s"
)

Set-Location $PSScriptRoot
$LOCUST = "python -m locust"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host " TESTES DE CARGA - Locust" -ForegroundColor Cyan
Write-Host " Usuarios: $Usuarios | Ramp-up: $RampUp | Duracao: $Duracao" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# A função precisa existir antes de ser chamada
function Executar-Teste {
    param(
        [string]$Nome,
        [string]$Arquivo,
        [string]$HostUrl,
        [string]$CsvName,
        [string]$GrpcHost,
        [string]$GrpcPort
    )

    Write-Host "`n[+] Iniciando: $Nome" -ForegroundColor Yellow

    $cmd = "$LOCUST -f $Arquivo --headless -u $Usuarios -r $RampUp -t $Duracao --csv=$CsvName"
    
    if (![string]::IsNullOrEmpty($HostUrl)) {
        $cmd += " --host=$HostUrl"
    }

    # Tratamento especial para o gRPC
    if (![string]::IsNullOrEmpty($GrpcHost) -and ![string]::IsNullOrEmpty($GrpcPort)) {
        $env:GRPC_HOST = $GrpcHost
        $env:GRPC_PORT = $GrpcPort
        Write-Host "    -> Alvo gRPC: ${GrpcHost}:${GrpcPort}" -ForegroundColor DarkGray
    }

    Write-Host "    -> Comando: $cmd" -ForegroundColor DarkGray
    
    # Executa o comando
    Invoke-Expression $cmd

    # Limpa as variáveis de ambiente para não vazar pro próximo teste
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

Write-Host "`n=== Gerando graficos ===" -ForegroundColor Cyan
python gerar_graficos.py

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host " TODOS OS TESTES CONCLUIDOS COM SUCESSO!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan