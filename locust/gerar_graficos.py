import pandas as pd
import matplotlib.pyplot as plt
import os

# Garante que a pasta de gráficos existe
os.makedirs('graficos', exist_ok=True)

tecnologias = ['rest', 'soap', 'graphql', 'grpc']
linguagens = ['go', 'java']
# Cores para o gráfico de evolução
cores_tech = {'rest': '#1f77b4', 'soap': '#d62728', 'graphql': '#2ca02c', 'grpc': '#9467bd'}

# ==========================================
# PARTE 1: GRÁFICOS DE BARRAS (RESULTADO FINAL)
# ==========================================
def ler_dados(coluna):
    dados = {'Go': [], 'Java': []}
    for tech in tecnologias:
        for lang in linguagens:
            arquivo = f"resultados_{tech}_{lang}_carga2_stats.csv"
            valor = 0
            if os.path.exists(arquivo):
                try:
                    df = pd.read_csv(arquivo)
                    linha = df[df['Name'] == 'Aggregated']
                    if not linha.empty:
                        valor = float(linha.iloc[0][coluna])
                except Exception as e:
                    pass
            if lang == 'go': dados['Go'].append(valor)
            else: dados['Java'].append(valor)
    return dados

def plotar_barras(dados, titulo, arquivo_saida, ylabel):
    x = range(len(tecnologias))
    width = 0.35
    fig, ax = plt.subplots(figsize=(10, 6))
    
    bars1 = ax.bar([pos - width/2 for pos in x], dados['Go'], width, label='Go', color='#00ADD8')
    bars2 = ax.bar([pos + width/2 for pos in x], dados['Java'], width, label='Java', color='#f89820')

    ax.set_ylabel(ylabel, fontweight='bold')
    ax.set_title(f"{titulo}\n(Carga 2 | 50 Usuários Simultâneos)", fontsize=14, fontweight='bold')
    ax.set_xticks(x)
    ax.set_xticklabels([t.upper() for t in tecnologias], fontweight='bold')
    ax.legend()

    for bars in [bars1, bars2]:
        for bar in bars:
            height = bar.get_height()
            ax.annotate(f'{height:.1f}',
                        xy=(bar.get_x() + bar.get_width() / 2, height),
                        xytext=(0, 3), textcoords="offset points",
                        ha='center', va='bottom', fontsize=10, fontweight='bold')

    ax.margins(y=0.15)
    plt.tight_layout()
    plt.savefig(os.path.join('graficos', arquivo_saida), dpi=300)
    plt.close()

# ==========================================
# PARTE 2: GRÁFICOS DE EVOLUÇÃO (LINHA DO TEMPO)
# ==========================================
def plotar_evolucao(lang, titulo, arquivo_saida):
    fig, ax1 = plt.subplots(figsize=(12, 6))
    ax2 = ax1.twinx()  # Cria um segundo eixo Y (lado direito) para os usuários

    user_plotted = False

    for tech in tecnologias:
        # Lê o histórico, que mostra os dados segundo a segundo
        arquivo = f"resultados_{tech}_{lang}_carga2_stats_history.csv"
        if os.path.exists(arquivo):
            try:
                df = pd.read_csv(arquivo)
                if 'Name' in df.columns and 'Aggregated' in df['Name'].values:
                    df = df[df['Name'] == 'Aggregated'].copy()
                
                if not df.empty:
                    # Cria a coluna de tempo começando do zero
                    df['Tempo (s)'] = df['Timestamp'] - df['Timestamp'].min()
                    
                    # Tenta pegar a latência da janela atual ou a média total
                    coluna_latencia = '50%' if '50%' in df.columns else 'Total Average Response Time'
                    
                    # Plota a linha do protocolo
                    ax1.plot(df['Tempo (s)'], df[coluna_latencia], label=tech.upper(), color=cores_tech[tech], linewidth=2)
                    
                    # Plota a linha de usuários apenas uma vez
                    if not user_plotted:
                        ax2.plot(df['Tempo (s)'], df['User Count'], label='Usuários Simultâneos', color='black', linestyle='--', linewidth=2, alpha=0.6)
                        user_plotted = True
            except Exception as e:
                print(f"Erro no histórico de {tech}: {e}")

    ax1.set_xlabel('Tempo de Teste (segundos)', fontweight='bold')
    ax1.set_ylabel('Tempo de Resposta (ms)', fontweight='bold')
    ax2.set_ylabel('Quantidade de Usuários', fontweight='bold', color='black')
    
    ax1.set_title(f"{titulo}\nComo a latência reage à entrada de usuários", fontsize=14, fontweight='bold')
    
    # Junta as duas legendas em uma caixa só
    lines_1, labels_1 = ax1.get_legend_handles_labels()
    lines_2, labels_2 = ax2.get_legend_handles_labels()
    ax1.legend(lines_1 + lines_2, labels_1 + labels_2, loc='upper left')

    plt.grid(True, linestyle=':', alpha=0.6)
    plt.tight_layout()
    plt.savefig(os.path.join('graficos', arquivo_saida), dpi=300)
    plt.close()

# ==========================================
# EXECUÇÃO DO SCRIPT
# ==========================================
print("Lendo os dados finais e gerando gráficos de barras...")
plotar_barras(ler_dados('Median Response Time'), 'Tempo de Resposta - Mediana', '1_latencia_mediana.png', 'Tempo (ms)')
plotar_barras(ler_dados('95%'), 'Tempo de Resposta - Percentil 95 (P95)', '2_latencia_p95.png', 'Tempo (ms)')
plotar_barras(ler_dados('Requests/s'), 'Vazão (Throughput)', '3_throughput.png', 'Requisições / Segundo')

print("Analisando o histórico para gerar gráficos de evolução temporal...")
plotar_evolucao('go', 'Evolução Temporal: Go (Goroutines)', '4_evolucao_go.png')
plotar_evolucao('java', 'Evolução Temporal: Java (Spring Boot)', '5_evolucao_java.png')

print("✅ Todos os gráficos foram gerados perfeitamente em locust/graficos/!")