import os

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd


BASE_DIR = os.path.dirname(os.path.abspath(__file__))
GRAFICOS_DIR = os.path.join(BASE_DIR, "graficos")
os.makedirs(GRAFICOS_DIR, exist_ok=True)

TECNOLOGIAS = ["rest", "soap", "graphql", "grpc"]
LINGUAGENS = ["go", "java"]
CARGAS = [50, 100]

SERIES = [
    ("go", 50, "Go (50 usuarios)", "#79c7f2"),
    ("go", 100, "Go (100 usuarios)", "#08aeca"),
    ("java", 50, "Java (50 usuarios)", "#ffd7af"),
    ("java", 100, "Java (100 usuarios)", "#ff981a"),
]


def limpar_graficos_antigos():
    for nome in os.listdir(GRAFICOS_DIR):
        if nome.lower().endswith(".png"):
            os.remove(os.path.join(GRAFICOS_DIR, nome))


def caminho_csv(tecnologia, linguagem, carga):
    return os.path.join(
        BASE_DIR,
        f"resultados_{tecnologia}_{linguagem}_carga2_{carga}u_stats.csv",
    )


def carregar_dados():
    linhas_agregado = []
    linhas_endpoint = []

    for tecnologia in TECNOLOGIAS:
        for linguagem in LINGUAGENS:
            for carga in CARGAS:
                arquivo = caminho_csv(tecnologia, linguagem, carga)
                if not os.path.exists(arquivo):
                    raise FileNotFoundError(f"Arquivo nao encontrado: {arquivo}")

                df = pd.read_csv(arquivo)
                linhas_setup = df["Name"].astype(str).str.startswith("[setup]").sum()
                if linhas_setup:
                    raise ValueError(
                        f"{arquivo} contem linhas [setup]. Rode os testes novamente."
                    )

                agregado = df[df["Name"] == "Aggregated"]
                if agregado.empty:
                    raise ValueError(f"{arquivo} nao contem linha Aggregated.")

                agregado = agregado.iloc[0]
                falhas = int(agregado["Failure Count"])
                if falhas:
                    raise ValueError(
                        f"{arquivo} contem {falhas} falhas. Nao gere grafico com cenario invalido."
                    )

                linhas_agregado.append(
                    {
                        "tecnologia": tecnologia,
                        "linguagem": linguagem,
                        "usuarios": carga,
                        "requests": int(agregado["Request Count"]),
                        "latencia_media_ms": float(agregado["Average Response Time"]),
                        "latencia_mediana_ms": float(agregado["Median Response Time"]),
                        "p95_ms": float(agregado["95%"]),
                        "p99_ms": float(agregado["99%"]),
                        "payload_bytes": float(agregado["Average Content Size"]),
                    }
                )

                endpoints = df[df["Name"] != "Aggregated"]
                for _, endpoint in endpoints.iterrows():
                    linhas_endpoint.append(
                        {
                            "tecnologia": tecnologia,
                            "linguagem": linguagem,
                            "usuarios": carga,
                            "endpoint": endpoint["Name"],
                            "requests": int(endpoint["Request Count"]),
                            "latencia_media_ms": float(endpoint["Average Response Time"]),
                            "latencia_mediana_ms": float(endpoint["Median Response Time"]),
                            "p95_ms": float(endpoint["95%"]),
                            "payload_bytes": float(endpoint["Average Content Size"]),
                        }
                    )

    return pd.DataFrame(linhas_agregado), pd.DataFrame(linhas_endpoint)


def valor_agregado(df, tecnologia, linguagem, usuarios, coluna):
    linha = df[
        (df["tecnologia"] == tecnologia)
        & (df["linguagem"] == linguagem)
        & (df["usuarios"] == usuarios)
    ]
    if linha.empty:
        return 0
    return float(linha.iloc[0][coluna])


def formatar_eixo(ax, ylabel):
    ax.set_ylabel(ylabel, fontweight="bold", fontsize=12)
    ax.set_xticks(np.arange(len(TECNOLOGIAS)))
    ax.set_xticklabels([t.upper() for t in TECNOLOGIAS], fontweight="bold", fontsize=12)
    ax.grid(axis="y", linestyle="--", alpha=0.55)
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)


def anotar_barras(ax, barras, formato="{:.1f}"):
    for barra in barras:
        altura = barra.get_height()
        deslocamento = 4 if altura >= 0 else -12
        va = "bottom" if altura >= 0 else "top"
        ax.annotate(
            formato.format(altura),
            xy=(barra.get_x() + barra.get_width() / 2, altura),
            xytext=(0, deslocamento),
            textcoords="offset points",
            ha="center",
            va=va,
            fontsize=9,
            fontweight="bold",
            rotation=45,
        )


def plotar_comparativo_carga(df, coluna, titulo, ylabel, arquivo_saida):
    x = np.arange(len(TECNOLOGIAS))
    largura = 0.2
    deslocamentos = [-1.5 * largura, -0.5 * largura, 0.5 * largura, 1.5 * largura]

    fig, ax = plt.subplots(figsize=(14, 7))

    max_valor = 0
    for deslocamento, (linguagem, usuarios, label, cor) in zip(deslocamentos, SERIES):
        valores = [
            valor_agregado(df, tecnologia, linguagem, usuarios, coluna)
            for tecnologia in TECNOLOGIAS
        ]
        max_valor = max(max_valor, max(valores))
        barras = ax.bar(x + deslocamento, valores, largura, label=label, color=cor)
        anotar_barras(ax, barras)

    ax.set_title(titulo, fontsize=16, fontweight="bold", pad=20)
    formatar_eixo(ax, ylabel)
    ax.set_ylim(0, max_valor * 1.28 if max_valor else 1)
    ax.legend(loc="upper center", bbox_to_anchor=(0.5, -0.1), ncol=4, frameon=False)

    plt.tight_layout()
    plt.savefig(os.path.join(GRAFICOS_DIR, arquivo_saida), dpi=300)
    plt.close()


def plotar_variacao_p95(df):
    tecnologias_variacao = ["rest", "soap", "graphql"]
    x = np.arange(len(tecnologias_variacao))
    largura = 0.32
    series = [("go", "Go", "#08aeca"), ("java", "Java", "#ff981a")]

    fig, ax = plt.subplots(figsize=(14, 7))

    todos_valores = []
    for indice, (linguagem, label, cor) in enumerate(series):
        valores = []
        for tecnologia in tecnologias_variacao:
            p95_50 = valor_agregado(df, tecnologia, linguagem, 50, "p95_ms")
            p95_100 = valor_agregado(df, tecnologia, linguagem, 100, "p95_ms")
            variacao = ((p95_100 - p95_50) / p95_50) * 100 if p95_50 else 0
            valores.append(variacao)
        todos_valores.extend(valores)
        barras = ax.bar(
            x + (indice - 0.5) * largura,
            valores,
            largura,
            label=label,
            color=cor,
        )
        anotar_barras(ax, barras, "{:.0f}%")

    ax.axhline(0, color="#555555", linewidth=1)
    ax.set_title(
        "Impacto da Carga no P95\n(cenarios com carga comparavel entre 50 e 100 usuarios)",
        fontsize=16,
        fontweight="bold",
        pad=20,
    )
    formatar_eixo(ax, "Variacao do P95 (%)")
    ax.set_xticks(x)
    ax.set_xticklabels(
        [t.upper() for t in tecnologias_variacao],
        fontweight="bold",
        fontsize=12,
    )
    menor = min(todos_valores + [0])
    maior = max(todos_valores + [0])
    margem = max((maior - menor) * 0.25, 10)
    ax.set_ylim(menor - margem, maior + margem)
    ax.legend(loc="upper center", bbox_to_anchor=(0.5, -0.1), ncol=2, frameon=False)

    plt.tight_layout()
    plt.savefig(os.path.join(GRAFICOS_DIR, "3_variacao_p95_50_100.png"), dpi=300)
    plt.close()


def plotar_payload_consolidado(df):
    x = np.arange(len(TECNOLOGIAS))
    largura = 0.34
    series = [("go", "Go", "#08aeca"), ("java", "Java", "#ff981a")]

    fig, ax = plt.subplots(figsize=(14, 7))

    max_valor = 0
    for indice, (linguagem, label, cor) in enumerate(series):
        valores = []
        for tecnologia in TECNOLOGIAS:
            subset = df[
                (df["tecnologia"] == tecnologia)
                & (df["linguagem"] == linguagem)
            ]
            valores.append(float(subset["payload_bytes"].mean()))
        max_valor = max(max_valor, max(valores))
        barras = ax.bar(
            x + (indice - 0.5) * largura,
            valores,
            largura,
            label=label,
            color=cor,
        )
        anotar_barras(ax, barras, "{:.0f}")

    ax.set_title(
        "Tamanho Medio da Resposta\n(media entre os cenarios de 50 e 100 usuarios)",
        fontsize=16,
        fontweight="bold",
        pad=20,
    )
    formatar_eixo(ax, "Bytes por resposta")
    ax.set_ylim(0, max_valor * 1.25 if max_valor else 1)
    ax.legend(loc="upper center", bbox_to_anchor=(0.5, -0.1), ncol=2, frameon=False)

    plt.tight_layout()
    plt.savefig(os.path.join(GRAFICOS_DIR, "4_payload_size_comp.png"), dpi=300)
    plt.close()


def salvar_resumos(df_agregado, df_endpoint):
    df_agregado.sort_values(["tecnologia", "linguagem", "usuarios"]).to_csv(
        os.path.join(GRAFICOS_DIR, "resumo_agregado.csv"),
        index=False,
    )
    df_endpoint.sort_values(["tecnologia", "linguagem", "usuarios", "endpoint"]).to_csv(
        os.path.join(GRAFICOS_DIR, "resumo_por_endpoint.csv"),
        index=False,
    )


def main():
    print("Carregando CSVs do Locust...")
    df_agregado, df_endpoint = carregar_dados()

    limpar_graficos_antigos()
    salvar_resumos(df_agregado, df_endpoint)

    plotar_comparativo_carga(
        df_agregado,
        "latencia_media_ms",
        "Latencia Media por Tecnologia\n(comparativo entre 50 e 100 usuarios)",
        "Tempo medio (ms)",
        "1_latencia_media.png",
    )
    plotar_comparativo_carga(
        df_agregado,
        "p95_ms",
        "Latencia P95 por Tecnologia\n(cauda de resposta sob carga)",
        "Tempo P95 (ms)",
        "2_latencia_p95.png",
    )
    plotar_variacao_p95(df_agregado)
    plotar_payload_consolidado(df_agregado)

    print("Graficos gerados em locust/graficos:")
    print("- 1_latencia_media.png")
    print("- 2_latencia_p95.png")
    print("- 3_variacao_p95_50_100.png")
    print("- 4_payload_size_comp.png")


if __name__ == "__main__":
    main()
