import os

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd


BASE_DIR = os.path.dirname(os.path.abspath(__file__))
GRAFICOS_DIR = os.path.join(BASE_DIR, "graficos")

os.makedirs(GRAFICOS_DIR, exist_ok=True)


TECNOLOGIAS = ["rest", "soap", "graphql", "grpc"]
LINGUAGENS = ["go", "java"]
CARGAS = [100, 300]

SERIES = [
    ("go", 100, "Go (100 usuários)", "#79c7f2"),
    ("go", 300, "Go (300 usuários)", "#08aeca"),
    ("java", 100, "Java (100 usuários)", "#ffd7af"),
    ("java", 300, "Java (300 usuários)", "#ff981a"),
]


def limpar_graficos_antigos():
    for nome in os.listdir(GRAFICOS_DIR):
        if nome.lower().endswith(".png"):
            os.remove(os.path.join(GRAFICOS_DIR, nome))


def caminho_csv(tecnologia, linguagem, carga):
    return os.path.join(
        BASE_DIR,
        f"resultados_{tecnologia}_{linguagem}_{carga}u_stats.csv",
    )


def obter_coluna(df, coluna, arquivo):
    if coluna not in df.columns:
        raise ValueError(f"Coluna '{coluna}' não encontrada em: {arquivo}")
    return df[coluna]


def carregar_dados():
    linhas_agregado = []
    linhas_endpoint = []
    arquivos_faltando = []

    for tecnologia in TECNOLOGIAS:
        for linguagem in LINGUAGENS:
            for carga in CARGAS:
                arquivo = caminho_csv(tecnologia, linguagem, carga)

                if not os.path.exists(arquivo):
                    arquivos_faltando.append(arquivo)
                    continue

                df = pd.read_csv(arquivo)

                linhas_setup = df["Name"].astype(str).str.startswith("[setup]").sum()
                if linhas_setup:
                    raise ValueError(
                        f"{arquivo} contém linhas [setup]. "
                        "Rode os testes novamente sem misturar setup com medição."
                    )

                agregado = df[df["Name"] == "Aggregated"]
                if agregado.empty:
                    raise ValueError(f"{arquivo} não contém linha Aggregated.")

                agregado = agregado.iloc[0]

                falhas = int(agregado["Failure Count"])
                if falhas:
                    raise ValueError(
                        f"{arquivo} contém {falhas} falhas. "
                        "Não gere gráfico com cenário inválido."
                    )

                linhas_agregado.append(
                    {
                        "tecnologia": tecnologia,
                        "linguagem": linguagem,
                        "usuarios": carga,
                        "requests": int(agregado["Request Count"]),
                        "falhas": int(agregado["Failure Count"]),
                        "latencia_mediana_ms": float(agregado["Median Response Time"]),
                        "p95_ms": float(agregado["95%"]),
                        "payload_bytes": float(agregado["Average Content Size"]),
                        "throughput_rps": float(agregado["Requests/s"]),
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
                            "falhas": int(endpoint["Failure Count"]),
                            "latencia_mediana_ms": float(endpoint["Median Response Time"]),
                            "p95_ms": float(endpoint["95%"]),
                            "payload_bytes": float(endpoint["Average Content Size"]),
                            "throughput_rps": float(endpoint["Requests/s"]),
                        }
                    )

    if arquivos_faltando:
        mensagem = "Arquivos CSV não encontrados:\n"
        mensagem += "\n".join(f"- {arquivo}" for arquivo in arquivos_faltando)
        raise FileNotFoundError(mensagem)

    return pd.DataFrame(linhas_agregado), pd.DataFrame(linhas_endpoint)


def valor_agregado(df, tecnologia, linguagem, usuarios, coluna):
    linha = df[
        (df["tecnologia"] == tecnologia)
        & (df["linguagem"] == linguagem)
        & (df["usuarios"] == usuarios)
    ]

    if linha.empty:
        return 0.0

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


def plotar_comparativo(df, coluna, titulo, ylabel, arquivo_saida, formato="{:.1f}"):
    x = np.arange(len(TECNOLOGIAS))
    largura = 0.2
    deslocamentos = [-1.5 * largura, -0.5 * largura, 0.5 * largura, 1.5 * largura]

    fig, ax = plt.subplots(figsize=(14, 7))

    max_valor = 0.0

    for deslocamento, (linguagem, usuarios, label, cor) in zip(deslocamentos, SERIES):
        valores = [
            valor_agregado(df, tecnologia, linguagem, usuarios, coluna)
            for tecnologia in TECNOLOGIAS
        ]

        max_valor = max(max_valor, max(valores))

        barras = ax.bar(
            x + deslocamento,
            valores,
            largura,
            label=label,
            color=cor,
        )

        anotar_barras(ax, barras, formato)

    ax.set_title(titulo, fontsize=16, fontweight="bold", pad=20)
    formatar_eixo(ax, ylabel)

    ax.set_ylim(0, max_valor * 1.28 if max_valor else 1)
    ax.legend(loc="upper center", bbox_to_anchor=(0.5, -0.1), ncol=4, frameon=False)

    plt.tight_layout()
    plt.savefig(os.path.join(GRAFICOS_DIR, arquivo_saida), dpi=300)
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

    plotar_comparativo(
        df_agregado,
        "latencia_mediana_ms",
        "Tempo Mediano por Tecnologia\n(comparativo entre 100 e 300 usuários)",
        "Tempo mediano (ms)",
        "1_tempo_mediano.png",
        "{:.1f}",
    )

    plotar_comparativo(
        df_agregado,
        "p95_ms",
        "Latência P95 por Tecnologia\n(comparativo entre 100 e 300 usuários)",
        "Tempo P95 (ms)",
        "2_p95.png",
        "{:.1f}",
    )

    plotar_comparativo(
        df_agregado,
        "payload_bytes",
        "Tamanho Médio do Payload por Tecnologia\n(comparativo entre 100 e 300 usuários)",
        "Bytes por resposta",
        "3_payload.png",
        "{:.0f}",
    )

    plotar_comparativo(
        df_agregado,
        "throughput_rps",
        "Throughput por Tecnologia\n(comparativo entre 100 e 300 usuários)",
        "Requisições por segundo",
        "4_throughput.png",
        "{:.1f}",
    )

    print("Gráficos gerados em locust/graficos:")
    print("- 1_tempo_mediano.png")
    print("- 2_p95.png")
    print("- 3_payload.png")
    print("- 4_throughput.png")
    print("- resumo_agregado.csv")
    print("- resumo_por_endpoint.csv")


if __name__ == "__main__":
    main()