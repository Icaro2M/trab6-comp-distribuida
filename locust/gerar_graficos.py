import pandas as pd
import matplotlib.pyplot as plt
import matplotlib
import os

matplotlib.rcParams['font.family'] = 'DejaVu Sans'

protocolos = ['rest', 'soap', 'graphql', 'grpc']
linguagens = ['go', 'java']
cargas     = ['carga1', 'carga2']
cores      = {'rest': '#4C72B0', 'soap': '#DD8452', 'graphql': '#55A868', 'grpc': '#C44E52'}

os.makedirs('graficos', exist_ok=True)

def ler_stats(protocolo, linguagem, carga):
    fname = f"resultados_{protocolo}_{linguagem}_{carga}_stats.csv"
    if not os.path.exists(fname):
        return None
    df = pd.read_csv(fname)
    agg = df[df['Name'] == 'Aggregated']
    if agg.empty:
        agg = df.iloc[-1:]
    return agg.iloc[0]

# ── Gráfico 1: Latência mediana por protocolo (Go vs Java, por carga) ──
for carga in cargas:
    fig, axes = plt.subplots(1, 2, figsize=(12, 5), sharey=True)
    fig.suptitle(f'Latência Mediana por Protocolo — {carga.upper()}', fontsize=14)
    for ax, lang in zip(axes, linguagens):
        medianas, labels, bar_cores = [], [], []
        for proto in protocolos:
            row = ler_stats(proto, lang, carga)
            if row is not None:
                medianas.append(row.get('50%', row.get('Median Response Time', 0)))
                labels.append(proto.upper())
                bar_cores.append(cores[proto])
        bars = ax.bar(labels, medianas, color=bar_cores)
        ax.bar_label(bars, fmt='%.0f ms', padding=3)
        ax.set_title(lang.upper())
        ax.set_ylabel('Latência Mediana (ms)')
        ax.set_ylim(0, max(medianas) * 1.3 if medianas else 100)
    plt.tight_layout()
    plt.savefig(f'graficos/latencia_mediana_{carga}.png', dpi=150)
    plt.close()
    print(f'✅ latencia_mediana_{carga}.png')

# ── Gráfico 2: Throughput (req/s) ──
for carga in cargas:
    fig, axes = plt.subplots(1, 2, figsize=(12, 5), sharey=True)
    fig.suptitle(f'Throughput (req/s) por Protocolo — {carga.upper()}', fontsize=14)
    for ax, lang in zip(axes, linguagens):
        rps, labels, bar_cores = [], [], []
        for proto in protocolos:
            row = ler_stats(proto, lang, carga)
            if row is not None:
                rps.append(row.get('Requests/s', 0))
                labels.append(proto.upper())
                bar_cores.append(cores[proto])
        bars = ax.bar(labels, rps, color=bar_cores)
        ax.bar_label(bars, fmt='%.1f', padding=3)
        ax.set_title(lang.upper())
        ax.set_ylabel('Requisições/s')
    plt.tight_layout()
    plt.savefig(f'graficos/throughput_{carga}.png', dpi=150)
    plt.close()
    print(f'✅ throughput_{carga}.png')

# ── Gráfico 3: P95 vs P99 — Carga 2 ──
for lang in linguagens:
    fig, ax = plt.subplots(figsize=(10, 5))
    x = range(len(protocolos))
    p95, p99 = [], []
    for proto in protocolos:
        row = ler_stats(proto, lang, 'carga2')
        p95.append(row.get('95%', 0) if row is not None else 0)
        p99.append(row.get('99%', 0) if row is not None else 0)
    bar_w = 0.35
    b1 = ax.bar([i - bar_w/2 for i in x], p95, bar_w, label='P95', color='#4C72B0')
    b2 = ax.bar([i + bar_w/2 for i in x], p99, bar_w, label='P99', color='#C44E52')
    ax.bar_label(b1, fmt='%.0f ms', padding=2, fontsize=8)
    ax.bar_label(b2, fmt='%.0f ms', padding=2, fontsize=8)
    ax.set_xticks(list(x))
    ax.set_xticklabels([p.upper() for p in protocolos])
    ax.set_ylabel('Latência (ms)')
    ax.set_title(f'Percentis P95 vs P99 — {lang.upper()} — Carga 2')
    ax.legend()
    plt.tight_layout()
    plt.savefig(f'graficos/percentis_{lang}_carga2.png', dpi=150)
    plt.close()
    print(f'✅ percentis_{lang}_carga2.png')

# ── Gráfico 4: Go vs Java — comparação direta ──
for carga in cargas:
    fig, ax = plt.subplots(figsize=(10, 5))
    x = range(len(protocolos))
    go_med   = []
    java_med = []
    for proto in protocolos:
        rg = ler_stats(proto, 'go', carga)
        rj = ler_stats(proto, 'java', carga)
        go_med.append(rg.get('50%', rg.get('Median Response Time', 0)) if rg is not None else 0)
        java_med.append(rj.get('50%', rj.get('Median Response Time', 0)) if rj is not None else 0)
    b1 = ax.bar([i - 0.2 for i in x], go_med,   0.35, label='Go',   color='#00ADD8')
    b2 = ax.bar([i + 0.2 for i in x], java_med, 0.35, label='Java', color='#E76F00')
    ax.bar_label(b1, fmt='%.0f ms', padding=2, fontsize=8)
    ax.bar_label(b2, fmt='%.0f ms', padding=2, fontsize=8)
    ax.set_xticks(list(x))
    ax.set_xticklabels([p.upper() for p in protocolos])
    ax.set_ylabel('Latência Mediana (ms)')
    ax.set_title(f'Go vs Java — Latência Mediana — {carga.upper()}')
    ax.legend()
    plt.tight_layout()
    plt.savefig(f'graficos/go_vs_java_{carga}.png', dpi=150)
    plt.close()
    print(f'✅ go_vs_java_{carga}.png')

print("\n✅ Todos os gráficos gerados em locust/graficos/")