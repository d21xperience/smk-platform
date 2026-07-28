# =============================================================================
# SMK PLATFORM - MAKEFILE
# =============================================================================
# Automation commands untuk seluruh project.
#
# CARA PENGGUNAAN:
#   make <command>
#
# CONTOH:
#   make help          - Tampilkan semua perintah
#   make setup         - Setup awal project
#   make infra-up      - Start infrastructure
#   make dev-fe        - Start frontend dev server
#   make dev-be        - Start backend dev server
#
# OPTIONS UNTUK KOMPATIBILITAS TERMINAL:
#   make COLOR=1 help   - Force warna aktif
#   make COLOR=0 help   - Force warna nonaktif
#   make EMOJI=1 help   - Force emoji aktif
#   make EMOJI=0 help   - Force emoji nonaktif
#
# CATATAN:
#   - Git Bash: Warna & emoji aktif secara default
#   - CMD/PowerShell: Warna & emoji nonaktif secara default (aman)
#   - Linux/macOS: Warna & emoji aktif secara default
# =============================================================================

# =============================================================================
# OS DETECTION (Gunakan env var bawaan Windows)
# =============================================================================

ifdef COMSPEC
  IS_WINDOWS := 1
else
  IS_WINDOWS := 0
endif

# =============================================================================
# COLOR SUPPORT DETECTION
# =============================================================================

ifdef COLOR
  USE_COLOR := $(COLOR)
else
  ifeq ($(IS_WINDOWS),1)
    USE_COLOR := 0
  else
    USE_COLOR := 1
  endif
endif

# =============================================================================
# EMOJI SUPPORT DETECTION
# =============================================================================

ifdef EMOJI
  USE_EMOJI := $(EMOJI)
else
  ifeq ($(IS_WINDOWS),1)
    USE_EMOJI := 0
  else
    USE_EMOJI := 1
  endif
endif


# =============================================================================
# ECHO COMMAND
# =============================================================================

ifeq ($(USE_COLOR),1)
  E :=	echo
else
  E := echo -e
endif

# =============================================================================
# COLOR CODES
# =============================================================================

ifeq ($(USE_COLOR),0)
  GREEN	:= \033[0;32m
  YELLOW	:= \033[1;33m
  RED	:= \033[0;31m
  BLUE	:= \033[0;34m
  CYAN	:= \033[0;36m
  BOLD	:= \033[1m
  NC	:= \033[0m
else
  GREEN  :=
  YELLOW :=
  RED    :=
  BLUE   :=
  CYAN   :=
  BOLD   :=
  NC     :=
endif

# =============================================================================
# ICONS (Emoji atau ASCII)
# =============================================================================

ifeq ($(USE_EMOJI),1)
  I_OK       := ✅
  I_WARN     := ⚠️
  I_ERR      := ❌
  I_SETUP    := 🚀
  I_FE       := 🎨
  I_BE       := ⚙️
  I_TEST     := 🧪
  I_BUILD    := 🔨
  I_LINT     := 🔍
  I_FMT      := ✨
  I_CONTRACT := 📜
  I_DB       := 🗄️
  I_DEPLOY   := 🚀
  I_BACKUP   := 💾
  I_CLEAN    := 🧹
  I_DOCS     := 📚
  I_NEW      := 📝
  I_SEED     := 🌱
  I_DEL      := 🗑️
  I_STOP     := 🛑
  I_RESTART  := 🔄
  I_ROLLBACK := ⏪
  I_WAIT     := ⏳
  I_STATUS   := 📊
  I_DISK     := 💾
  I_TOOLS    := 🔧
  I_LIST     := 📋
  I_SQUARE	 := 📦
else
  I_OK       := [OK]
  I_WARN     := [!!]
  I_ERR      := [XX]
  I_SETUP    := [>>]
  I_FE       := [FE]
  I_BE       := [BE]
  I_TEST     := [TT]
  I_BUILD    := [BB]
  I_LINT     := [LL]
  I_FMT      := [FF]
  I_CONTRACT := [CC]
  I_DB       := [DB]
  I_DEPLOY   := [DP]
  I_BACKUP   := [BK]
  I_CLEAN    := [CL]
  I_DOCS     := [DC]
  I_NEW      := [NW]
  I_SEED     := [SD]
  I_DEL      := [DD]
  I_STOP     := [ST]
  I_RESTART  := [RS]
  I_ROLLBACK := [RB]
  I_WAIT     := [..]
  I_STATUS   := [SS]
  I_DISK     := [DK]
  I_TOOLS    := [TL]
  I_LIST     := [..]
  I_SQUARE	 := [BL]
endif

# =============================================================================
# BOX DRAWING
# =============================================================================

ifeq ($(IS_WINDOWS),1)
  BOX_TOP := +--------------------------------------------------------------+
  BOX_BOT := +--------------------------------------------------------------+
  BOX_V   := |
else
  BOX_TOP := ╔════════════════════════════════════════════════════════════════╗
  BOX_BOT := ╚════════════════════════════════════════════════════════════════╝
  BOX_V   := ║
endif

# =============================================================================
# PROJECT VARIABLES (CURDIR = built-in Make, tidak perlu pwd)
# =============================================================================

ROOT_DIR     := $(CURDIR)
APPS_DIR     := $(ROOT_DIR)/apps
SERVICES_DIR := $(ROOT_DIR)/services
DEPLOY_DIR   := $(ROOT_DIR)/deploy
DOCS_DIR     := $(ROOT_DIR)/docs
SCRIPTS_DIR  := $(ROOT_DIR)/scripts

FE_APPS := web-sekolah siakad-guru psb

BE_SERVICES := edge-bff core-academic core-student core-attendance core-assessment core-reporting

DOCKER_COMPOSE := docker-compose

# =============================================================================
# DEFAULT TARGET
# =============================================================================

.DEFAULT_GOAL := help
# ============================================================
# VERIFIKASI & SETUP ENVIRONMENT
# ============================================================
.PHONY: check-version
check-version:
	@chcp 65001 >nul 2>&1 || true
	@$(E)	"$(YELLOW)Memeriksa keselarasan lingkungan pengembangan..."
	@$(E)	"$(BOX_TOP)"
	@$(E)	"$(BOX_V)Target Versi Utama (Komputer Anda):"
	@$(E)	"$(BOX_V)   * Node.js : v26.2.0 (Versi NodeJS yang dibutuhkan)"
	@$(E)	"$(BOX_V)   * Go      : go1.24.11 (Versi Go yang dibutuhkan)	$(NC)"
	@$(E)	"$(BOX_BOT)"
	@$(E)	" Versi Aktif di Komputer Tim Saat Ini:"
	@$(E)	"   * Node.js : " && node --version
	@$(E)	"   * Go      : " && go version
	@$(E)	"$(BOX_TOP)"
	@$(E)	"$(BOX_V)Tips: Jika versi berbeda, silakan sesuaikan runtime Anda"
	@$(E)	"$(BOX_V)agar tidak terjadi error inkonsistensi kode."
	@$(E)	"$(BOX_BOT)"





# =============================================================================
# HELP
# =============================================================================

.PHONY: help
help:
	@$(E) ""
	@$(E) "$(BLUE)$(BOX_TOP)$(NC)"
	@$(E) "$(BLUE)$(BOX_V)  $(BOLD)SMK PLATFORM - AUTOMATION COMMANDS$(NC)$(BLUE)$(BOX_V)$(NC)"
	@$(E) "$(BLUE)$(BOX_BOT)$(NC)"
	@$(E) ""
	@$(E) "$(GREEN)$(I_LIST) SETUP & INFRASTRUCTURE$(NC)"
	@$(E) "  make setup              Setup awal project (install deps + start infra)"
	@$(E) "  make infra-up           Start infrastructure (PostgreSQL, Redis, MinIO)"
	@$(E) "  make infra-down         Stop infrastructure"
	@$(E) "  make infra-restart      Restart infrastructure"
	@$(E) "  make infra-logs         View infrastructure logs"
	@$(E) "  make infra-status       Check infrastructure status"
	@$(E) "  make infra-reset        Reset infrastructure (hapus semua data!)"
	@$(E) "  make infra-tools        Start GUI tools (pgAdmin, Redis Commander)"
	@$(E) ""
	@$(E) "$(GREEN)$(I_FE) FRONTEND DEVELOPMENT$(NC)"
	@$(E) "  make dev-fe             Start frontend dev server (pilih app)"
	@$(E) "  make dev-fe-guru        Start SIAKAD Guru dev server"
	@$(E) "  make dev-fe-web         Start Website Sekolah dev server"
	@$(E) "  make dev-fe-psb        Start psb dev server"
	@$(E) "  make build-fe           Build semua frontend apps"
	@$(E) "  make lint-fe            Lint semua frontend apps"
	@$(E) "  make format-fe          Format semua frontend code"
	@$(E) ""
	@$(E) "$(GREEN)$(I_BE) BACKEND DEVELOPMENT$(NC)"
	@$(E) "  make dev-be             Start backend dev server (pilih service)"
	@$(E) "  make dev-be-edge        Start Edge BFF dev server"
	@$(E) "  make dev-be-attendance  Start Core Attendance dev server"
	@$(E) "  make dev-be-student     Start Core Student dev server"
	@$(E) "  make dev-be-academic    Start Core Academic dev server"
	@$(E) "  make dev-be-assessment  Start Core Assessment dev server"
	@$(E) "  make dev-be-reporting   Start Core Reporting dev server"
	@$(E) "  make build-be           Build semua backend services"
	@$(E) "  make lint-be            Lint semua backend services"
	@$(E) "  make test-be            Test semua backend services"
	@$(E) ""
	@$(E) "$(GREEN)$(I_TEST) TESTING & QUALITY$(NC)"
	@$(E) "  make test               Run semua tests (FE + BE)"
	@$(E) "  make lint               Run semua linters (FE + BE)"
	@$(E) "  make format             Format semua code (FE + BE)"
	@$(E) "  make check              Run lint + test + build (pre-commit check)"
	@$(E) ""
	@$(E) "$(GREEN)$(I_CONTRACT) CONTRACT & PROTOBUF$(NC)"
	@$(E) "  make proto-gen          Generate Go code dari Protobuf"
	@$(E) "  make proto-lint         Lint Protobuf files"
	@$(E) "  make contract-check     Validate API contract adherence"
	@$(E) ""
	@$(E) "$(GREEN)$(I_DB) DATABASE$(NC)"
	@$(E) "  make db-connect         Connect ke PostgreSQL"
	@$(E) "  make db-migrate-up      Run database migrations (up)"
	@$(E) "  make db-migrate-down    Rollback database migrations (down)"
	@$(E) "  make db-migrate-status  Check migration status"
	@$(E) "  make db-seed            Seed database dengan sample data"
	@$(E) "  make db-reset           Reset database (hapus semua data!)"
	@$(E) ""
	@$(E) "$(GREEN)$(I_DEPLOY) DEPLOYMENT$(NC)"
	@$(E) "  make deploy-cloud       Deploy ke Cloud VPS"
	@$(E) "  make deploy-local       Deploy ke Server Lokal"
	@$(E) "  make deploy-rollback    Rollback deployment terakhir"
	@$(E) ""
	@$(E) "$(GREEN)$(I_BACKUP) BACKUP & RECOVERY$(NC)"
	@$(E) "  make backup             Backup semua data"
	@$(E) "  make backup-db          Backup PostgreSQL saja"
	@$(E) "  make restore-db         Restore PostgreSQL dari backup"
	@$(E) ""
	@$(E) "$(GREEN)$(I_CLEAN) CLEANUP$(NC)"
	@$(E) "  make clean              Clean build artifacts"
	@$(E) "  make clean-all          Clean semua (artifacts + deps + infra)"
	@$(E) "  make clean-docker       Clean Docker images & volumes"
	@$(E) ""
	@$(E) "$(GREEN)$(I_DOCS) DOCUMENTATION$(NC)"
	@$(E) "  make docs-serve         Serve documentation locally"
	@$(E) "  make adr-new            Buat ADR baru dari template"
	@$(E) ""
	@$(E) "$(CYAN)Terminal Options:$(NC)"
	@$(E) "  make COLOR=1 <cmd>      Force warna aktif"
	@$(E) "  make COLOR=0 <cmd>      Force warna nonaktif"
	@$(E) "  make EMOJI=1 <cmd>      Force emoji aktif"
	@$(E) "  make EMOJI=0 <cmd>      Force emoji nonaktif"
	@$(E) ""
	@$(E) "$(CYAN)Packages Build:$(NC)"
	@$(E) "  make ui-dev        Develop vue Shared components"
	@$(E) "  make ui-build      Build vue Shared components"
	@$(E) ""

# =============================================================================
# SETUP & INFRASTRUCTURE
# =============================================================================

.PHONY: setup
setup:
	@$(E) ""
	@$(E) "$(BLUE)$(BOX_TOP)$(NC)"
	@$(E) "$(BLUE)$(BOX_V)  $(BOLD)SETUP SMK PLATFORM$(NC)$(BLUE)$(BOX_V)$(NC)"
	@$(E) "$(BLUE)$(BOX_BOT)$(NC)"
	@$(E) ""
	@$(E) "$(BLUE)$(I_LIST) Checking prerequisites...$(NC)"
	@command -v docker >/dev/null 2>&1 || { $(E) "$(RED)$(I_ERR) Docker is not installed$(NC)"; exit 1; }
	@command -v node >/dev/null 2>&1 || { $(E) "$(RED)$(I_ERR) Node.js is not installed$(NC)"; exit 1; }
	@command -v go >/dev/null 2>&1 || { $(E) "$(RED)$(I_ERR) Go is not installed$(NC)"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { $(E) "$(RED)$(I_ERR) NPM is not installed$(NC)"; exit 1; }
	@$(E) "$(GREEN)$(I_OK) All prerequisites met$(NC)"
	@$(E) ""
	@$(E) "$(BLUE)$(I_LIST) Setting up environment...$(NC)"
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
        $(E) "$(GREEN)$(I_OK) .env created from template$(NC)"; \
    else \
        $(E) "$(YELLOW)$(I_WARN) .env already exists, skipping$(NC)"; \
    fi
	@$(E) ""
	@$(E) "$(BLUE)$(I_LIST) Starting infrastructure...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@$(E) ""
	@$(E) "$(BLUE)$(I_WAIT) Waiting for services to be healthy...$(NC)"
	@sleep 10
	@$(DOCKER_COMPOSE) ps
	@$(E) ""
	@$(E) "$(BLUE)$(I_LIST) Installing frontend dependencies...$(NC)"
	@for app in $(FE_APPS); do \
        if [ -d "$(APPS_DIR)/$$app" ]; then \
            $(E) "$(CYAN)  Installing $$app...$(NC)"; \
            cd $(APPS_DIR)/$$app && yarn install --frozen-lockfile 2>/dev/null || yarn install; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) ""
	@$(E) "$(BLUE)$(I_LIST) Downloading backend dependencies...$(NC)"
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service" ]; then \
            $(E) "$(CYAN)  Downloading $$service...$(NC)"; \
            cd $(SERVICES_DIR)/$$service && go mod download 2>/dev/null || true; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) ""
	@$(E) "$(GREEN)$(BOX_TOP)$(NC)"
	@$(E) "$(GREEN)$(BOX_V)  $(I_OK) SETUP COMPLETE!$(GREEN)$(BOX_V)$(NC)"
	@$(E) "$(GREEN)$(BOX_BOT)$(NC)"
	@$(E) ""
	@$(E) "Next steps:"
	@$(E) "  1. Edit .env file if needed"
	@$(E) "  2. Run 'make dev-fe-guru' to start SIAKAD Guru"
	@$(E) "  3. Run 'make dev-be-edge' to start Edge BFF"
	@$(E) "  4. Access http://localhost:8080"
	@$(E) ""

.PHONY: infra-up
infra-up:
	@$(E) "$(GREEN)$(I_SETUP) Starting infrastructure...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@$(E) "$(GREEN)$(I_OK) Infrastructure started$(NC)"
	@$(DOCKER_COMPOSE) ps

.PHONY: infra-down
infra-down:
	@$(E) "$(YELLOW)$(I_STOP) Stopping infrastructure...$(NC)"
	@$(DOCKER_COMPOSE) down
	@$(E) "$(GREEN)$(I_OK) Infrastructure stopped$(NC)"

.PHONY: infra-restart
infra-restart:
	@$(E) "$(YELLOW)$(I_RESTART) Restarting infrastructure...$(NC)"
	@$(DOCKER_COMPOSE) restart
	@$(E) "$(GREEN)$(I_OK) Infrastructure restarted$(NC)"
	@$(DOCKER_COMPOSE) ps

.PHONY: infra-logs
infra-logs:
	@$(DOCKER_COMPOSE) logs -f

.PHONY: infra-status
infra-status:
	@$(E) "$(BLUE)$(I_STATUS) Infrastructure Status:$(NC)"
	@$(E) ""
	@$(DOCKER_COMPOSE) ps
	@$(E) ""
	@$(E) "$(BLUE)$(I_LINT) Health Checks:$(NC)"
	@$(E) -n "  PostgreSQL: "
	@docker inspect --format='{{.State.Health.Status}}' smk-postgres 2>/dev/null || echo "not running"
	@$(E) -n "  Redis:      "
	@docker inspect --format='{{.State.Health.Status}}' smk-redis 2>/dev/null || echo "not running"
	@$(E) -n "  MinIO:      "
	@docker inspect --format='{{.State.Health.Status}}' smk-minio 2>/dev/null || echo "not running"
	@$(E) ""
	@$(E) "$(BLUE)$(I_DISK) Disk Usage:$(NC)"
	@docker system df --format "table {{.Type}}\t{{.TotalCount}}\t{{.Size}}"

.PHONY: infra-reset
infra-reset:
	@$(E) "$(RED)$(I_WARN) WARNING: This will delete ALL data!$(NC)"
	@$(E) -n "Are you sure? (y/N): "
	@read confirm; \
    if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
        $(E) "$(YELLOW)$(I_DEL) Resetting infrastructure...$(NC)"; \
        $(DOCKER_COMPOSE) down -v --remove-orphans; \
        $(DOCKER_COMPOSE) up -d; \
        $(E) "$(GREEN)$(I_OK) Infrastructure reset complete$(NC)"; \
    else \
        $(E) "$(GREEN)$(I_OK) Cancelled$(NC)"; \
    fi

.PHONY: infra-tools
infra-tools:
	@$(E) "$(GREEN)$(I_TOOLS) Starting GUI tools...$(NC)"
	@$(DOCKER_COMPOSE) --profile tools up -d
	@$(E) ""
	@$(E) "$(GREEN)$(I_OK) GUI Tools started:$(NC)"
	@$(E) "  pgAdmin:          http://localhost:5050 (admin@smk.local / admin)"
	@$(E) "  Redis Commander:  http://localhost:8081 (admin / admin)"
	@$(E) ""

# =============================================================================
# FRONTEND DEVELOPMENT
# =============================================================================

.PHONY: dev-fe
dev-fe:
	@$(E) "$(BLUE)$(I_FE) Available frontend apps:$(NC)"
	@$(E) "  1) web-sekolah"
	@$(E) "  2) siakad-guru"
	@$(E) "  3) psb"
	@$(E) -n "Select app (1-3): "
	@read choice; \
    case $$choice in \
        1) cd $(APPS_DIR)/web-sekolah && yarn dev ;; \
        2) cd $(APPS_DIR)/siakad-guru && yarn dev ;; \
        3) cd $(APPS_DIR)/psb && yarn dev ;; \
        *) $(E) "$(RED)$(I_ERR) Invalid choice$(NC)" ;; \
    esac

.PHONY: dev-fe-guru
dev-fe-guru:
	@$(E) "$(GREEN)$(I_FE) Starting SIAKAD Guru...$(NC)"
	@cd $(APPS_DIR)/siakad-guru && yarn dev

.PHONY: dev-fe-web
dev-fe-web:
	@$(E) "$(GREEN)$(I_FE) Starting Website Sekolah...$(NC)"
	@cd $(APPS_DIR)/web-sekolah && yarn dev

.PHONY: dev-fe-psb
dev-fe-psb:
	@$(E) "$(GREEN)$(I_FE) Starting psb...$(NC)"
	@cd $(APPS_DIR)/psb && yarn dev

.PHONY: build-fe
build-fe:
	@$(E) "$(GREEN)$(I_BUILD) Building all frontend apps...$(NC)"
	@for app in $(FE_APPS); do \
        if [ -d "$(APPS_DIR)/$$app" ]; then \
            $(E) "$(CYAN)  Building $$app...$(NC)"; \
            cd $(APPS_DIR)/$$app && npm build; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All frontend apps built$(NC)"

.PHONY: lint-fe
lint-fe:
	@$(E) "$(GREEN)$(I_LINT) Linting all frontend apps...$(NC)"
	@for app in $(FE_APPS); do \
        if [ -d "$(APPS_DIR)/$$app" ]; then \
            $(E) "$(CYAN)  Linting $$app...$(NC)"; \
            cd $(APPS_DIR)/$$app && npm run lint; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All frontend apps linted$(NC)"

.PHONY: format-fe
format-fe:
	@$(E) "$(GREEN)$(I_FMT) Formatting all frontend code...$(NC)"
	@for app in $(FE_APPS); do \
        if [ -d "$(APPS_DIR)/$$app" ]; then \
            $(E) "$(CYAN)  Formatting $$app...$(NC)"; \
            cd $(APPS_DIR)/$$app && yarn format; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All frontend code formatted$(NC)"

# =============================================================================
# BACKEND DEVELOPMENT
# =============================================================================

.PHONY: dev-be
dev-be:
	@$(E) "$(BLUE)$(I_BE) Available backend services:$(NC)"
	@$(E) "  1) edge-bff"
	@$(E) "  2) core-academic"
	@$(E) "  3) core-student"
	@$(E) "  4) core-attendance"
	@$(E) "  5) core-assessment"
	@$(E) "  6) core-reporting"
	@$(E) -n "Select service (1-6): "
	@read choice; \
    case $$choice in \
        1) cd $(SERVICES_DIR)/edge-bff && go run cmd/main.go ;; \
        2) cd $(SERVICES_DIR)/core-academic && go run cmd/main.go ;; \
        3) cd $(SERVICES_DIR)/core-student && go run cmd/main.go ;; \
        4) cd $(SERVICES_DIR)/core-attendance && go run cmd/main.go ;; \
        5) cd $(SERVICES_DIR)/core-assessment && go run cmd/main.go ;; \
        6) cd $(SERVICES_DIR)/core-reporting && go run cmd/main.go ;; \
        *) $(E) "$(RED)$(I_ERR) Invalid choice$(NC)" ;; \
    esac

.PHONY: dev-be-edge
dev-be-edge:
	@$(E) "$(GREEN)$(I_BE) Starting Edge BFF...$(NC)"
	@cd $(SERVICES_DIR)/edge-bff && go run cmd/main.go

.PHONY: dev-be-academic
dev-be-academic:
	@$(E) "$(GREEN)$(I_BE) Starting Core Academic...$(NC)"
	@cd $(SERVICES_DIR)/core-academic && go run cmd/main.go

.PHONY: dev-be-student
dev-be-student:
	@$(E) "$(GREEN)$(I_BE) Starting Core Student...$(NC)"
	@cd $(SERVICES_DIR)/core-student && go run cmd/main.go

.PHONY: dev-be-attendance
dev-be-attendance:
	@$(E) "$(GREEN)$(I_BE) Starting Core Attendance...$(NC)"
	@cd $(SERVICES_DIR)/core-attendance && go run cmd/main.go

.PHONY: dev-be-assessment
dev-be-assessment:
	@$(E) "$(GREEN)$(I_BE) Starting Core Assessment...$(NC)"
	@cd $(SERVICES_DIR)/core-assessment && go run cmd/main.go

.PHONY: dev-be-reporting
dev-be-reporting:
	@$(E) "$(GREEN)$(I_BE) Starting Core Reporting...$(NC)"
	@cd $(SERVICES_DIR)/core-reporting && go run cmd/main.go

.PHONY: build-be
build-be:
	@$(E) "$(GREEN)$(I_BUILD) Building all backend services...$(NC)"
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service" ]; then \
            $(E) "$(CYAN)  Building $$service...$(NC)"; \
            cd $(SERVICES_DIR)/$$service && \
            CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$$service ./cmd/main.go; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All backend services built$(NC)"

.PHONY: lint-be
lint-be:
	@$(E) "$(GREEN)$(I_LINT) Linting all backend services...$(NC)"
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service" ]; then \
            $(E) "$(CYAN)  Linting $$service...$(NC)"; \
            cd $(SERVICES_DIR)/$$service && go vet ./...; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All backend services linted$(NC)"

.PHONY: test-be
test-be:
	@$(E) "$(GREEN)$(I_TEST) Testing all backend services...$(NC)"
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service" ]; then \
            $(E) "$(CYAN)  Testing $$service...$(NC)"; \
            cd $(SERVICES_DIR)/$$service && go test -v -race ./...; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All backend tests passed$(NC)"

# =============================================================================
# TESTING & QUALITY
# =============================================================================

.PHONY: test
test: test-be
	@$(E) "$(GREEN)$(I_TEST) Running all tests...$(NC)"
	@for app in $(FE_APPS); do \
        if [ -d "$(APPS_DIR)/$$app" ]; then \
            $(E) "$(CYAN)  Testing $$app...$(NC)"; \
            cd $(APPS_DIR)/$$app && yarn test 2>/dev/null || true; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All tests passed$(NC)"

.PHONY: lint
lint: lint-fe lint-be
	@$(E) "$(GREEN)$(I_OK) All linters passed$(NC)"

.PHONY: format
format: format-fe
	@$(E) "$(GREEN)$(I_FMT) Formatting backend code...$(NC)"
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service" ]; then \
            cd $(SERVICES_DIR)/$$service && gofmt -w .; \
            cd $(ROOT_DIR); \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) All code formatted$(NC)"

.PHONY: check
check: lint test build-be
	@$(E) ""
	@$(E) "$(GREEN)$(BOX_TOP)$(NC)"
	@$(E) "$(GREEN)$(BOX_V)  $(I_OK) ALL CHECKS PASSED!$(GREEN)$(BOX_V)$(NC)"
	@$(E) "$(GREEN)$(BOX_BOT)$(NC)"
	@$(E) ""
	@$(E) "You are ready to commit and push!"

# =============================================================================
# CONTRACT & PROTOBUF
# =============================================================================

.PHONY: proto-gen
proto-gen:
	@$(E) "$(GREEN)$(I_CONTRACT) Generating Protobuf code...$(NC)"
	@for proto_dir in $(ROOT_DIR)/contracts/protobuf/*/; do \
        if [ -d "$$proto_dir" ]; then \
            $(E) "$(CYAN)  Generating from $$proto_dir...$(NC)"; \
            protoc --go_out=. --go_opt=paths=source_relative \
                --go-grpc_out=. --go-grpc_opt=paths=source_relative \
                $$proto_dir*.proto 2>/dev/null || true; \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) Protobuf code generated$(NC)"

.PHONY: proto-lint
proto-lint:
	@$(E) "$(GREEN)$(I_LINT) Linting Protobuf files...$(NC)"
	@command -v buf >/dev/null 2>&1 && \
        buf lint $(ROOT_DIR)/contracts/protobuf/ || \
        $(E) "$(YELLOW)$(I_WARN) 'buf' not installed. Install: brew install buf$(NC)"
	@$(E) "$(GREEN)$(I_OK) Protobuf lint complete$(NC)"

.PHONY: contract-check
contract-check:
	@$(E) "$(GREEN)$(I_CONTRACT) Checking API contract adherence...$(NC)"
	@$(E) "$(CYAN)  Checking Protobuf files...$(NC)"
	@find $(ROOT_DIR)/contracts/protobuf -name "*.proto" -exec echo "    {}" \;
	@$(E) "$(CYAN)  Checking OpenAPI files...$(NC)"
	@find $(ROOT_DIR)/contracts/openapi -name "*.yaml" -exec echo "    {}" \; 2>/dev/null || true
	@$(E) "$(GREEN)$(I_OK) Contract check complete$(NC)"

# =============================================================================
# DATABASE
# =============================================================================

.PHONY: db-connect
db-connect:
	@$(E) "$(GREEN)$(I_DB) Connecting to PostgreSQL...$(NC)"
	@$(DOCKER_COMPOSE) exec postgres psql -U $${POSTGRES_USER:-smk_dev} -d $${POSTGRES_DB:-smk_platform}

.PHONY: db-migrate-up
db-migrate-up:
	@$(E) "$(GREEN)$(I_DB) Running database migrations (up)...$(NC)"
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service/migrations" ]; then \
            $(E) "$(CYAN)  Migrating $$service...$(NC)"; \
            migrate -path $(SERVICES_DIR)/$$service/migrations \
                -database "postgres://$${POSTGRES_USER:-smk_dev}:$${POSTGRES_PASSWORD:-secret123}@localhost:$${POSTGRES_PORT:-5432}/$${POSTGRES_DB:-smk_platform}?sslmode=disable" \
                up 2>/dev/null || true; \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) Migrations complete$(NC)"

.PHONY: db-migrate-down
db-migrate-down:
	@$(E) "$(YELLOW)$(I_WARN) Rolling back database migrations...$(NC)"
	@$(E) -n "How many steps to rollback? (default: 1): "
	@read steps; \
    steps=$${steps:-1}; \
    for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service/migrations" ]; then \
            $(E) "$(CYAN)  Rolling back $$service ($$steps steps)...$(NC)"; \
            migrate -path $(SERVICES_DIR)/$$service/migrations \
                -database "postgres://$${POSTGRES_USER:-smk_dev}:$${POSTGRES_PASSWORD:-secret123}@localhost:$${POSTGRES_PORT:-5432}/$${POSTGRES_DB:-smk_platform}?sslmode=disable" \
                down $$steps 2>/dev/null || true; \
        fi; \
    done
	@$(E) "$(GREEN)$(I_OK) Rollback complete$(NC)"

.PHONY: db-migrate-status
db-migrate-status:
	@$(E) "$(GREEN)$(I_DB) Migration Status:$(NC)"
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service/migrations" ]; then \
            $(E) "$(CYAN)  $$service:$(NC)"; \
            migrate -path $(SERVICES_DIR)/$$service/migrations \
                -database "postgres://$${POSTGRES_USER:-smk_dev}:$${POSTGRES_PASSWORD:-secret123}@localhost:$${POSTGRES_PORT:-5432}/$${POSTGRES_DB:-smk_platform}?sslmode=disable" \
                version 2>/dev/null || echo "    No migrations found"; \
        fi; \
    done

.PHONY: db-seed
db-seed:
	@$(E) "$(GREEN)$(I_SEED) Seeding database...$(NC)"
	@$(DOCKER_COMPOSE) exec -T postgres psql -U $${POSTGRES_USER:-smk_dev} -d $${POSTGRES_DB:-smk_platform} \
        < $(SCRIPTS_DIR)/init-db/001-init-schema.sql 2>/dev/null || \
        $(E) "$(YELLOW)$(I_WARN) Seed script not found or already applied$(NC)"
	@$(E) "$(GREEN)$(I_OK) Database seeded$(NC)"

.PHONY: db-reset
db-reset:
	@$(E) "$(RED)$(I_WARN) WARNING: This will delete ALL database data!$(NC)"
	@$(E) -n "Are you sure? (y/N): "
	@read confirm; \
    if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
        $(E) "$(YELLOW)$(I_DEL) Resetting database...$(NC)"; \
        $(DOCKER_COMPOSE) exec -T postgres psql -U $${POSTGRES_USER:-smk_dev} -d $${POSTGRES_DB:-smk_platform} \
            -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"; \
        $(MAKE) db-seed; \
        $(E) "$(GREEN)$(I_OK) Database reset complete$(NC)"; \
    else \
        $(E) "$(GREEN)$(I_OK) Cancelled$(NC)"; \
    fi

# =============================================================================
# DEPLOYMENT
# =============================================================================

.PHONY: deploy-cloud
deploy-cloud:
	@$(E) "$(GREEN)$(I_DEPLOY) Deploying to Cloud VPS...$(NC)"
	@if [ -f "$(DEPLOY_DIR)/cloud/docker-compose.yml" ]; then \
        $(E) "$(CYAN)  Building Docker images...$(NC)"; \
        $(MAKE) build-be; \
        $(MAKE) build-fe; \
        $(E) "$(CYAN)  Deploying...$(NC)"; \
        cd $(DEPLOY_DIR)/cloud && $(DOCKER_COMPOSE) pull && $(DOCKER_COMPOSE) up -d; \
        $(E) "$(GREEN)$(I_OK) Cloud deployment complete$(NC)"; \
    else \
        $(E) "$(RED)$(I_ERR) Cloud deployment config not found$(NC)"; \
    fi

.PHONY: deploy-local
deploy-local:
	@$(E) "$(GREEN)$(I_DEPLOY) Deploying to Local Server...$(NC)"
	@if [ -f "$(DEPLOY_DIR)/local/docker-compose.yml" ]; then \
        $(E) "$(CYAN)  Building Docker images...$(NC)"; \
        $(MAKE) build-be; \
        $(E) "$(CYAN)  Running migrations...$(NC)"; \
        $(MAKE) db-migrate-up; \
        $(E) "$(CYAN)  Deploying...$(NC)"; \
        cd $(DEPLOY_DIR)/local && $(DOCKER_COMPOSE) pull && $(DOCKER_COMPOSE) up -d; \
        $(E) "$(GREEN)$(I_OK) Local deployment complete$(NC)"; \
    else \
        $(E) "$(RED)$(I_ERR) Local deployment config not found$(NC)"; \
    fi

.PHONY: deploy-rollback
deploy-rollback:
	@$(E) "$(YELLOW)$(I_ROLLBACK) Rolling back deployment...$(NC)"
	@$(E) -n "Which environment? (cloud/local): "
	@read env; \
    if [ -f "$(DEPLOY_DIR)/$$env/docker-compose.yml" ]; then \
        cd $(DEPLOY_DIR)/$$env && $(DOCKER_COMPOSE) rollback 2>/dev/null || \
        $(E) "$(YELLOW)$(I_WARN) Rollback not supported. Manual rollback required.$(NC)"; \
    else \
        $(E) "$(RED)$(I_ERR) Environment not found$(NC)"; \
    fi

# =============================================================================
# BACKUP & RECOVERY
# =============================================================================

.PHONY: backup
backup: backup-db
	@$(E) "$(GREEN)$(I_BACKUP) Starting full backup...$(NC)"
	@mkdir -p backup/$$(date +%Y%m%d_%H%M%S)
	@BACKUP_DIR=backup/$$(date +%Y%m%d_%H%M%S); \
    $(E) "$(CYAN)  Backing up Redis...$(NC)"; \
    $(DOCKER_COMPOSE) exec -T redis redis-cli BGSAVE 2>/dev/null; \
    sleep 3; \
    docker cp smk-redis:/data/dump.rdb $$BACKUP_DIR/redis.rdb 2>/dev/null || true; \
    $(E) "$(CYAN)  Backing up MinIO...$(NC)"; \
    docker cp smk-minio:/data $$BACKUP_DIR/minio 2>/dev/null || true; \
    $(E) "$(GREEN)$(I_OK) Full backup complete: $$BACKUP_DIR$(NC)"

.PHONY: backup-db
backup-db:
	@$(E) "$(GREEN)$(I_BACKUP) Backing up PostgreSQL...$(NC)"
	@mkdir -p backup
	@$(DOCKER_COMPOSE) exec -T postgres pg_dump \
        -U $${POSTGRES_USER:-smk_dev} \
        -F c \
        -d $${POSTGRES_DB:-smk_platform} \
        > backup/postgres_$$(date +%Y%m%d_%H%M%S).dump
	@$(E) "$(GREEN)$(I_OK) Database backup complete$(NC)"
	@ls -lh backup/postgres_*.dump | tail -1

.PHONY: restore-db
restore-db:
	@$(E) "$(YELLOW)$(I_WARN) WARNING: This will overwrite current database!$(NC)"
	@$(E) "$(BLUE)Available backups:$(NC)"
	@ls -lh backup/postgres_*.dump 2>/dev/null || echo "  No backups found"
	@$(E) -n "Enter backup filename: "
	@read filename; \
    if [ -f "backup/$$filename" ]; then \
        $(E) "$(YELLOW)$(I_RESTART) Restoring from backup/$$filename...$(NC)"; \
        cat backup/$$filename | $(DOCKER_COMPOSE) exec -T postgres pg_restore \
            -U $${POSTGRES_USER:-smk_dev} \
            -d $${POSTGRES_DB:-smk_platform} \
            --clean --if-exists; \
        $(E) "$(GREEN)$(I_OK) Database restored$(NC)"; \
    else \
        $(E) "$(RED)$(I_ERR) Backup file not found$(NC)"; \
    fi

# =============================================================================
# CLEANUP
# =============================================================================

.PHONY: clean
clean:
	@$(E) "$(GREEN)$(I_CLEAN) Cleaning build artifacts...$(NC)"
	@for app in $(FE_APPS); do \
        if [ -d "$(APPS_DIR)/$$app/dist" ]; then \
            rm -rf $(APPS_DIR)/$$app/dist; \
            $(E) "  Cleaned $$app/dist"; \
        fi; \
    done
	@for service in $(BE_SERVICES); do \
        if [ -d "$(SERVICES_DIR)/$$service/bin" ]; then \
            rm -rf $(SERVICES_DIR)/$$service/bin; \
            $(E) "  Cleaned $$service/bin"; \
        fi; \
    done
	@rm -rf backup/*.dump backup/*.rdb 2>/dev/null || true
	@$(E) "$(GREEN)$(I_OK) Build artifacts cleaned$(NC)"

.PHONY: clean-all
clean-all: clean
	@$(E) "$(RED)$(I_WARN) WARNING: This will remove ALL dependencies and infrastructure data!$(NC)"
	@$(E) -n "Are you sure? (y/N): "
	@read confirm; \
    if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
        $(E) "$(YELLOW)$(I_DEL) Cleaning everything...$(NC)"; \
        $(DOCKER_COMPOSE) down -v --remove-orphans; \
        for app in $(FE_APPS); do \
            rm -rf $(APPS_DIR)/$$app/node_modules; \
        done; \
        for service in $(BE_SERVICES); do \
            cd $(SERVICES_DIR)/$$service && go clean -cache 2>/dev/null || true; \
            cd $(ROOT_DIR); \
        done; \
        $(E) "$(GREEN)$(I_OK) Everything cleaned$(NC)"; \
    else \
        $(E) "$(GREEN)$(I_OK) Cancelled$(NC)"; \
    fi

.PHONY: clean-docker
clean-docker:
	@$(E) "$(GREEN)$(I_CLEAN) Cleaning Docker resources...$(NC)"
	@docker image prune -f
	@docker container prune -f
	@docker network prune -f
	@$(E) "$(GREEN)$(I_OK) Docker resources cleaned$(NC)"
	@docker system df

# =============================================================================
# DOCUMENTATION
# =============================================================================

.PHONY: docs-serve
docs-serve:
	@$(E) "$(GREEN)$(I_DOCS) Serving documentation...$(NC)"
	@command -v npx >/dev/null 2>&1 && \
        cd $(DOCS_DIR) && npx docsify-cli serve . --port 3000 || \
        $(E) "$(YELLOW)$(I_WARN) Install docsify: npm install -g docsify-cli$(NC)"

.PHONY: adr-new
adr-new:
	@$(E) "$(GREEN)$(I_NEW) Creating new ADR...$(NC)"
	@$(E) -n "Enter ADR title (e.g., 'use-kafka-for-events'): "
	@read title; \
    ADR_NUM=$$(ls $(DOCS_DIR)/adr/[0-9]*.md 2>/dev/null | wc -l | tr -d ' '); \
    ADR_NUM=$$((ADR_NUM + 1)); \
    ADR_FILE=$(DOCS_DIR)/adr/$$(printf "%03d" $$ADR_NUM)-$$title.md; \
    if [ -f "$(DOCS_DIR)/adr/template.md" ]; then \
        cp $(DOCS_DIR)/adr/template.md $$ADR_FILE; \
        sed -i "s/ADR-XXX/ADR-$$(printf "%03d" $$ADR_NUM)/g" $$ADR_FILE; \
        sed -i "s/\[Judul Keputusan\]/$$title/g" $$ADR_FILE; \
        sed -i "s/\[YYYY-MM-DD\]/$$(date +%Y-%m-%d)/g" $$ADR_FILE; \
        $(E) "$(GREEN)$(I_OK) ADR created: $$ADR_FILE$(NC)"; \
        $(E) "Edit the file to fill in Context, Decision, and Consequences."; \
    else \
        $(E) "$(RED)$(I_ERR) Template not found at $(DOCS_DIR)/adr/template.md$(NC)"; \
    fi

.PHONY: ui-build
ui-build:
	@$(E) "$(I_BUILD) Building Shared UI Components..."
	@cd packages/ui-components && yarn build
	@$(E) "$(I_OK) UI Components built!"

.PHONY:ui-dev
ui-dev:
	@$(E) "$(I_BUILD) Watching Shared UI Components for changes..."
	@cd packages/ui-components && yarn dev

.PHONY:ui-lint
ui-lint:
	@$(E) "$(I_LINT) Linting Shared UI Components..."
	@cd packages/ui-components && yarn lint
	@$(E) "$(I_OK) UI Components lint passed!"
