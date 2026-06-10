import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import os

# Garante que a pasta existe
os.makedirs('graficos', exist_ok=True)

tecnologias = ['rest', 'soap', 'graphql', 'grpc']
linguagens = ['go', 'java']
cargas = ['50', '100']

def ler_dados(coluna):
    # Armazena os cenários separados
    dados = {'go_50': [], 'go_100': [], 'java_50': [], 'java_100': []}
    
    for tech in tecnologias:
        for lang in linguagens:
            for carga in cargas:
                arquivo = f"resultados_{tech}_{lang}_carga2_{carga}u_stats.csv"
                valor = 0
                if os.path.exists(arquivo):
                    try:
                        df = pd.read_csv(arquivo)
                        linha = df[df['Name'] == 'Aggregated']
                        if not linha.empty:
                            valor = float(linha.iloc[0][coluna])
                    except Exception as e:
                        print(f"Aviso ao ler {arquivo}: {e}")
                
                dados[f"{lang}_{carga}"].append(valor)
    return dados

def plotar_agrupado(dados, titulo, arquivo_saida, ylabel):
    x = np.arange(len(tecnologias))
    width = 0.2  # Largura de cada barra

    fig, ax = plt.subplots(figsize=(14, 7))

    # Paleta de Cores Claras (50u) e Escuras (100u)
    cor_go50 = '#87CEFA'   # Azul claro
    cor_go100 = '#00ADD8'  # Azul escuro
    cor_java50 = '#FFDAB9' # Laranja claro
    cor_java100 = '#f89820'# Laranja escuro

    # Desenhando as 4 barras para cada protocolo
    bars1 = ax.bar(x - 1.5 * width, dados['go_50'], width, label='Go (50 Usuários)', color=cor_go50)
    bars2 = ax.bar(x - 0.5 * width, dados['go_100'], width, label='Go (100 Usuários)', color=cor_go100)
    bars3 = ax.bar(x + 0.5 * width, dados['java_50'], width, label='Java (50 Usuários)', color=cor_java50)
    bars4 = ax.bar(x + 1.5 * width, dados['java_100'], width, label='Java (100 Usuários)', color=cor_java100)

    # Estilização
    ax.set_ylabel(ylabel, fontweight='bold', fontsize=12)
    ax.set_title(f"{titulo}\n(Comparativo de Escalabilidade: 50 vs 100 Usuários)", fontsize=16, fontweight='bold', pad=20)
    ax.set_xticks(x)
    ax.set_xticklabels([t.upper() for t in tecnologias], fontweight='bold', fontsize=12)
    ax.legend(loc='upper center', bbox_to_anchor=(0.5, -0.1), ncol=4, frameon=False, fontsize=11)

    # Adiciona os números no topo das barras com rotação para não embolar
    for bars in [bars1, bars2, bars3, bars4]:
        for bar in bars:
            height = bar.get_height()
            if height > 0:
                ax.annotate(f'{height:.1f}',
                            xy=(bar.get_x() + bar.get_width() / 2, height),
                            xytext=(0, 5),  # Deslocamento vertical de 5 pontos
                            textcoords="offset points",
                            ha='center', va='bottom', fontsize=9, fontweight='bold', rotation=45)

    # Dá um respiro no topo do gráfico para os números não cortarem
    max_y = max([max(dados[k]) for k in dados]) if any(max(dados[k]) > 0 for k in dados) else 1
    ax.set_ylim(0, max_y * 1.25) 

    plt.grid(axis='y', linestyle='--', alpha=0.7)
    plt.tight_layout()
    plt.savefig(os.path.join('graficos', arquivo_saida), dpi=300)
    plt.close()

print("Processando os dados de 50 e 100 usuários...")

# 1. Mediana
plotar_agrupado(ler_dados('Median Response Time'), 'Tempo de Resposta - Mediana', '1_latencia_mediana_comp.png', 'Tempo (ms)')

# 2. P95
plotar_agrupado(ler_dados('95%'), 'Estabilidade (P95)', '2_latencia_p95_comp.png', 'Tempo (ms)')

# 3. Vazão
plotar_agrupado(ler_dados('Requests/s'), 'Vazão (Throughput)', '3_throughput_comp.png', 'Requisições / Segundo')

# 4. Tamanho do Payload (O triunfo final do gRPC!)
plotar_agrupado(ler_dados('Average Content Size'), 'Consumo de Rede - Tamanho Médio da Resposta', '4_payload_size_comp.png', 'Tamanho Médio (Bytes)')

print("✅ Todos os 4 gráficos finais gerados na pasta 'graficos'!")