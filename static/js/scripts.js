// Funções para números
function formatCurrency(value) {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(value);
}

function formatPercent(value) {
  return new Intl.NumberFormat("pt-BR", {
    style: "percent",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value / 100);
}

// Funções para manipulação do DOM
function scrollToElement(elementId) {
  const element = document.getElementById(elementId);
  if (element) {
    element.scrollIntoView({ behavior: "smooth", block: "start" });
  }
}

// Event listeners
document.addEventListener("DOMContentLoaded", function () {
  setupForm();
  setupNavigation();
  setupThemeToggle();
});

// Configuração do tema claro/escuro
function setupThemeToggle() {
  const toggleButton = document.getElementById("theme-toggle");
  if (toggleButton) {
    // Verificar se há preferência salva
    const isDarkMode = localStorage.getItem("darkMode") === "true";
    if (isDarkMode) {
      document.body.classList.add("dark-mode");
      toggleButton.innerHTML = '<i class="fas fa-sun"></i> Modo Claro';
    }

    toggleButton.addEventListener("click", function () {
      document.body.classList.toggle("dark-mode");
      const isDark = document.body.classList.contains("dark-mode");
      localStorage.setItem("darkMode", isDark);

      if (isDark) {
        toggleButton.innerHTML = '<i class="fas fa-sun"></i> Modo Claro';
      } else {
        toggleButton.innerHTML = '<i class="fas fa-moon"></i> Modo Escuro';
      }

      // Atualizar gráficos se existirem
      if (typeof initCharts === "function") {
        initCharts();
      }
    });
  }
}

// Configuração do formulário
function setupForm() {
  const form = document.getElementById("calculator-form");
  if (form) {
    form.addEventListener("submit", function (event) {
      event.preventDefault();

      // Validar o formulário
      const investmentAmountField =
        document.getElementById("investment-amount");
      const investmentAmountValue = investmentAmountField.value.trim();

      // Verificar se está vazio
      if (investmentAmountValue === "") {
        showAlert("Por favor, insira um valor para investimento.", "danger");
        return;
      }

      // Mostrar indicador de carregamento
      const loadingIndicator = document.getElementById("loading-indicator");
      loadingIndicator.classList.remove("d-none");

      // Limpar resultados anteriores
      const resultsContainer = document.getElementById("results-container");
      resultsContainer.innerHTML = "";

      // Criar FormData e adicionar os valores
      const formData = new FormData();
      formData.append("valorInvestimento", investmentAmountValue);

      console.log("Enviando valor inicial:", investmentAmountValue);

      // Verificar se a distribuição personalizada está ativada
      if (
        document.getElementById("distribuicao-personalizada") &&
        document.getElementById("distribuicao-personalizada").checked
      ) {
        const tiposSelecionados = Array.from(
          document.querySelectorAll(".tipo-investimento:checked")
        ).map((checkbox) => checkbox.value);

        formData.append("tiposInvestimento", JSON.stringify(tiposSelecionados));
        formData.append("distribuicaoPersonalizada", "true");

        console.log("Tipos de investimento selecionados:", tiposSelecionados);
      }

      // Enviar a requisição usando fetch com FormData
      fetch("/calcular", {
        method: "POST",
        body: formData,
      })
        .then((response) => {
          console.log("Status da resposta:", response.status);
          return response.json();
        })
        .then((data) => {
          // Ocultar indicador de carregamento
          loadingIndicator.classList.add("d-none");

          console.log("Resposta recebida:", data);

          if (data.status === "success") {
            // Mostrar resultados
            resultsContainer.innerHTML = data.dados_html;

            // Scroll para os resultados
            scrollToElement("results");

            // Inicializar componentes do Bootstrap
            initBootstrapComponents();

            // Inicializar gráficos
            if (typeof window.initCharts === "function") {
              window.initCharts();
            }
          } else {
            // Mostrar mensagem de erro
            showAlert(data.message, "danger");
          }
        })
        .catch((error) => {
          // Ocultar indicador de carregamento
          loadingIndicator.classList.add("d-none");

          console.error("Erro na requisição:", error);

          // Mostrar mensagem de erro
          showAlert(
            `Erro ao processar a solicitação: ${error.message}`,
            "danger"
          );
        });
    });

    // Formatar o campo de valor com máscara de moeda
    const investmentField = document.getElementById("investment-amount");
    if (investmentField) {
      // Formatação de moeda
      investmentField.addEventListener("input", function (e) {
        let value = e.target.value;

        // Remover todos os caracteres não numéricos exceto vírgula e ponto
        value = value.replace(/[^\d.,]/g, "");

        // Garantir que só exista uma vírgula ou ponto decimal
        const match = value.match(/[.,]/g);
        if (match && match.length > 1) {
          // Manter apenas o último separador decimal
          const lastIndex = value.lastIndexOf(match[match.length - 1]);
          value =
            value.substring(0, lastIndex).replace(/[.,]/g, "") +
            value.substring(lastIndex);
        }

        // Formatar com separadores de milhares
        if (value.includes(",") || value.includes(".")) {
          const parts = value.split(/[,.]/);
          const integerPart = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ".");
          const decimalPart = parts[1] || "";
          value = integerPart + "," + decimalPart;
        } else if (value) {
          value = value.replace(/\B(?=(\d{3})+(?!\d))/g, ".");
        }

        e.target.value = value;
      });
    }
  }

  setupFormDistribuicao();
}

// Configuração da navegação
function setupNavigation() {
  const navLinks = document.querySelectorAll(".nav-link");

  navLinks.forEach((link) => {
    link.addEventListener("click", function (event) {
      const targetId = this.getAttribute("href");

      if (targetId.startsWith("#")) {
        event.preventDefault();

        // Scroll para a seção
        scrollToElement(targetId.substring(1));
      }
    });
  });
}

// Inicializar componentes do Bootstrap
function initBootstrapComponents() {
  // Ativar tooltips
  const tooltipTriggerList = [].slice.call(
    document.querySelectorAll('[data-bs-toggle="tooltip"]')
  );
  tooltipTriggerList.forEach(function (tooltipTriggerEl) {
    new bootstrap.Tooltip(tooltipTriggerEl);
  });

  // Ativar popovers
  const popoverTriggerList = [].slice.call(
    document.querySelectorAll('[data-bs-toggle="popover"]')
  );
  popoverTriggerList.forEach(function (popoverTriggerEl) {
    new bootstrap.Popover(popoverTriggerEl);
  });

  // Ativar tabs
  const tabTriggerList = [].slice.call(
    document.querySelectorAll('[data-bs-toggle="tab"]')
  );
  tabTriggerList.forEach(function (tabTriggerEl) {
    new bootstrap.Tab(tabTriggerEl);
  });
}

// Função para mostrar alertas
function showAlert(message, type) {
  const alertsContainer = document.getElementById("results-container");
  const alertElement = document.createElement("div");
  alertElement.className = `alert alert-${type} alert-dismissible fade show`;
  alertElement.innerHTML = `
    ${message}
    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Fechar"></button>
  `;

  alertsContainer.innerHTML = "";
  alertsContainer.appendChild(alertElement);

  // Scroll para o alerta
  alertElement.scrollIntoView({ behavior: "smooth" });

  // Auto-remover após 5 segundos para alertas de sucesso
  if (type === "success") {
    setTimeout(() => {
      alertElement.classList.remove("show");
      setTimeout(() => alertElement.remove(), 150);
    }, 5000);
  }
}

function setupFormDistribuicao() {
  const distribuicaoCheckbox = document.getElementById(
    "distribuicao-personalizada"
  );
  const opcoesDistribuicao = document.getElementById("opcoes-distribuicao");

  if (distribuicaoCheckbox && opcoesDistribuicao) {
    distribuicaoCheckbox.addEventListener("change", function () {
      if (this.checked) {
        opcoesDistribuicao.classList.remove("d-none");
      } else {
        opcoesDistribuicao.classList.add("d-none");
      }
    });

    // Garantir que pelo menos uma opção esteja marcada
    const checkboxesTipos = document.querySelectorAll(".tipo-investimento");
    checkboxesTipos.forEach((checkbox) => {
      checkbox.addEventListener("change", function () {
        const peloMenosUmMarcado = Array.from(checkboxesTipos).some(
          (cb) => cb.checked
        );
        if (!peloMenosUmMarcado) {
          this.checked = true;
          showAlert(
            "Você deve selecionar pelo menos um tipo de investimento.",
            "warning"
          );
        }
      });
    });
  }
}
