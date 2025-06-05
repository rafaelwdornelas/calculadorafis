# 📊 Calculadora de Investimentos

Uma aplicação web desenvolvida em Go (Golang) para otimização de carteiras de investimentos, oferecendo recomendações inteligentes de alocação de ativos com base em princípios de diversificação e balanceamento.

![Go Version](https://img.shields.io/badge/Go-1.24.1-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/status-active-success?style=for-the-badge)

## 🚀 Funcionalidades

### 📈 Análise de Carteira
- Integração com API do Investidor10 para obter dados da carteira atual
- Suporte para múltiplas classes de ativos: FIIs, Ações, ETFs e Renda Fixa
- Cálculo automático da distribuição atual vs ideal

### 💡 Recomendações Inteligentes
- Algoritmo de otimização para sugerir as melhores alocações
- Análise de data ex-dividendo para maximizar rendimentos
- Distribuição personalizada ou padrão (30% FIIs, 30% Ações, 10% ETFs, 30% Renda Fixa)

### 📊 Visualizações e Relatórios
- Dashboard interativo com gráficos e tabelas
- Projeção de rendimentos mensais e anuais
- Análise por segmento e tipo de ativo
- Relatório completo para impressão

### ⚡ Performance
- Sistema de cache inteligente para cotações (30 minutos)
- Processamento otimizado de grandes volumes de dados
- Interface responsiva com carregamento assíncrono

## 🛠️ Tecnologias Utilizadas

- **Backend**: Go 1.24.1
- **Frontend**: HTML5, CSS3, JavaScript, Bootstrap 5
- **APIs**: BrAPI (cotações), Investidor10 (carteiras)
- **Bibliotecas**: 
 - `gofpdf` - Geração de PDFs
 - `Chart.js` - Gráficos interativos
 - `Bootstrap` - Interface responsiva

## 📋 Pré-requisitos

- Go 1.24.1 ou superior
- Conexão com internet para acessar APIs externas
- Token de API da BrAPI (incluído no código)

## 🔧 Instalação

1. Clone o repositório:
git clone https://github.com/rafaelwdornelas/calculadorafis.git
cd calculadorafis

2. Instale as dependências:
go mod download

3. Configure o arquivo `internal/config/config.go` com suas credenciais:
IDInvestidor10: "SEU_ID_AQUI",  // ID da sua carteira no Investidor10

4. Execute a aplicação:
go run main.go

5. Acesse no navegador:
http://localhost:5000

## 🎯 Como Usar

1. **Digite o valor do investimento**: Informe quanto deseja investir
2. **Personalize a distribuição** (opcional): Escolha quais classes de ativos incluir
3. **Calcule**: Clique no botão para gerar as recomendações
4. **Analise os resultados**: 
  - Veja as recomendações de compra por classe de ativo
  - Confira a análise de data ex-dividendo
  - Visualize como ficará sua carteira final
  - Veja a projeção de rendimentos

## 📊 Dados de Entrada

Os arquivos de recomendações em `data/` devem seguir o formato:

**recomendados_fiis.txt:**
TICKER	NOME	SEGMENTO	TIPO	PESO%
HGLG11	PÁTRIA LOG	Logístico	Fundo de Tijolo	7,14%

**recomendados_acoes.txt:**
NOME	TICKER	PESO%
Banco do Brasil	BBAS3	4,00%

## 🔒 Segurança

- Não exponha suas credenciais do Investidor10
- O token da BrAPI incluído é público e tem limite de requisições
- Recomenda-se usar variáveis de ambiente para dados sensíveis

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:

1. Faça um Fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## ⚠️ Disclaimer

Esta calculadora é uma ferramenta educacional e não constitui recomendação de investimento. Sempre consulte um profissional qualificado antes de tomar decisões financeiras.

## 👨‍💻 Autor

**Rafael W. Dornelas**
- GitHub: [@rafaelwdornelas](https://github.com/rafaelwdornelas)

## 🙏 Agradecimentos

- [BrAPI](https://brapi.dev) - API de cotações
- [Investidor10](https://investidor10.com.br) - Dados de carteiras
- Comunidade Go Brasil

---

⭐ Se este projeto foi útil para você, considere dar uma estrela!