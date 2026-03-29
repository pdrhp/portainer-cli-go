# Env Vars With Spaces on Stack Creation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Permitir variaveis de ambiente com espacos no valor (ex.: cron expression) em `stacks create-swarm-git` e `stacks redeploy`, mantendo validacao de chave e mensagens de erro claras.

**Architecture:** Centralizar o parsing/validacao de `KEY=VALUE` em um parser compartilhado e reaproveitar no fluxo por flags e no wizard interativo. A regra principal sera: validar a chave (`KEY`), preservar o valor como texto livre (incluindo espacos), e reportar erro com contexto (linha/entrada invalida) quando necessario.

**Tech Stack:** Go, Cobra, Huh (wizard TUI), Testify

---

### Task 1: Criar testes de parser (falhando primeiro)

**Files:**
- Create: `internal/envvars/parser_test.go`
- Test: `internal/envvars/parser_test.go`

**Step 1: Write the failing test**

```go
package envvars

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestParseAssignment_AllowsSpacesInValue(t *testing.T) {
    pair, err := ParseAssignment(`SCHEDULER_SYNC_PEAK=*/15 * * * *`)
    require.NoError(t, err)
    assert.Equal(t, "SCHEDULER_SYNC_PEAK", pair.Name)
    assert.Equal(t, "*/15 * * * *", pair.Value)
}

func TestParseAssignment_StripsOnlyOuterQuotes(t *testing.T) {
    pair, err := ParseAssignment(`SCHEDULER_SYNC_PEAK="*/15 * * * *"`)
    require.NoError(t, err)
    assert.Equal(t, "*/15 * * * *", pair.Value)
}

func TestParseAssignment_RejectsMissingEquals(t *testing.T) {
    _, err := ParseAssignment("INVALID")
    require.Error(t, err)
    assert.Contains(t, err.Error(), "formato")
}

func TestParseAssignment_RejectsInvalidKey(t *testing.T) {
    _, err := ParseAssignment("BAD KEY=value")
    require.Error(t, err)
    assert.Contains(t, err.Error(), "chave")
}

func TestParseAssignments_MultipleLines(t *testing.T) {
    pairs, err := ParseAssignments([]string{
        `A=1`,
        `SCHEDULER_SYNC_PEAK=*/15 * * * *`,
    })
    require.NoError(t, err)
    assert.Len(t, pairs, 2)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/envvars -run Parse -v`
Expected: FAIL com erro de compilacao (`undefined: ParseAssignment`).

**Step 3: Write minimal implementation**

Criar stubs temporarios em `internal/envvars/parser.go` apenas para compilar e manter os testes falhando por comportamento.

```go
package envvars

import "github.com/pdrhp/portainer-go-cli/pkg/types"

func ParseAssignment(input string) (types.Pair, error) {
    return types.Pair{}, nil
}

func ParseAssignments(inputs []string) ([]types.Pair, error) {
    return nil, nil
}
```

**Step 4: Run test to verify it fails**

Run: `go test ./internal/envvars -run Parse -v`
Expected: FAIL por asserts de comportamento.

**Step 5: Commit**

```bash
git add internal/envvars/parser_test.go internal/envvars/parser.go
git commit -m "test: add failing env var parser scenarios with spaced values"
```

### Task 2: Implementar parser compartilhado para KEY=VALUE

**Files:**
- Modify: `internal/envvars/parser.go`
- Test: `internal/envvars/parser_test.go`

**Step 1: Write the failing test**

Adicionar cenarios finais antes da implementacao completa:

```go
func TestParseAssignment_AllowsEmptyValue(t *testing.T) {
    pair, err := ParseAssignment("EMPTY=")
    require.NoError(t, err)
    assert.Equal(t, "", pair.Value)
}

func TestParseAssignments_ReturnsLineContextOnError(t *testing.T) {
    _, err := ParseAssignments([]string{"A=1", "BAD KEY=2"})
    require.Error(t, err)
    assert.Contains(t, err.Error(), "linha 2")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/envvars -v`
Expected: FAIL com os novos cenarios.

**Step 3: Write minimal implementation**

Implementacao completa sugerida:

```go
package envvars

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/pdrhp/portainer-go-cli/pkg/types"
)

var envKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func ParseAssignment(input string) (types.Pair, error) {
    raw := strings.TrimSpace(input)
    if raw == "" {
        return types.Pair{}, fmt.Errorf("variavel vazia")
    }

    parts := strings.SplitN(raw, "=", 2)
    if len(parts) != 2 {
        return types.Pair{}, fmt.Errorf("formato esperado: KEY=VALUE")
    }

    key := strings.TrimSpace(parts[0])
    value := parts[1]

    if key == "" {
        return types.Pair{}, fmt.Errorf("chave da variavel nao pode ser vazia")
    }
    if !envKeyPattern.MatchString(key) {
        return types.Pair{}, fmt.Errorf("chave da variavel invalida: %s", key)
    }

    if len(value) >= 2 {
        if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
            (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
            value = value[1 : len(value)-1]
        }
    }

    return types.Pair{Name: key, Value: value}, nil
}

func ParseAssignments(inputs []string) ([]types.Pair, error) {
    pairs := make([]types.Pair, 0, len(inputs))
    for i, line := range inputs {
        trimmed := strings.TrimSpace(line)
        if trimmed == "" {
            continue
        }

        pair, err := ParseAssignment(trimmed)
        if err != nil {
            return nil, fmt.Errorf("linha %d: %w", i+1, err)
        }
        pairs = append(pairs, pair)
    }
    return pairs, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/envvars -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/envvars/parser.go internal/envvars/parser_test.go
git commit -m "feat: parse env assignments allowing spaces in values"
```

### Task 3: Integrar parser no create-swarm-git (flags) com erro explicito

**Files:**
- Modify: `cmd/stacks_create_swarm_git.go`
- Modify: `cmd/stacks_create_swarm_git_test.go`
- Test: `cmd/stacks_create_swarm_git_test.go`

**Step 1: Write the failing test**

```go
func TestBuildPayloadFromFlags_WithEnvContainingSpaces(t *testing.T) {
    createSwarmGitName = "test-stack"
    createSwarmGitRepositoryURL = "https://github.com/user/repo"
    createSwarmGitSwarmID = "swarm"
    createSwarmGitEnv = []string{`SCHEDULER_SYNC_PEAK=*/15 * * * *`}

    payload, err := buildPayloadFromFlags()
    require.NoError(t, err)
    require.Len(t, payload.Env, 1)
    assert.Equal(t, "*/15 * * * *", payload.Env[0].Value)
}

func TestBuildPayloadFromFlags_InvalidEnvReturnsError(t *testing.T) {
    createSwarmGitName = "test-stack"
    createSwarmGitRepositoryURL = "https://github.com/user/repo"
    createSwarmGitSwarmID = "swarm"
    createSwarmGitEnv = []string{"BAD KEY=value"}

    _, err := buildPayloadFromFlags()
    require.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd -run BuildPayloadFromFlags -v`
Expected: FAIL por assinatura/fluxo atual sem erro.

**Step 3: Write minimal implementation**

- Alterar assinatura de `buildPayloadFromFlags` para retornar `(types.StackCreateSwarmGitPayload, error)`.
- Usar `envvars.ParseAssignments(createSwarmGitEnv)`.
- Em `RunE`, tratar o erro com mensagem contextual:

```go
payload, err = buildPayloadFromFlags()
if err != nil {
    return fmt.Errorf("invalid --env value: %w", err)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cmd -run BuildPayloadFromFlags -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add cmd/stacks_create_swarm_git.go cmd/stacks_create_swarm_git_test.go
git commit -m "fix: support spaced env values in create-swarm-git flags"
```

### Task 4: Integrar parser no redeploy (flags) com validacao igual

**Files:**
- Modify: `cmd/stacks_redeploy_git.go`
- Modify: `cmd/stacks_redeploy_git_test.go`
- Test: `cmd/stacks_redeploy_git_test.go`

**Step 1: Write the failing test**

```go
func TestBuildRedeployPayloadFromFlags_WithEnvContainingSpaces(t *testing.T) {
    redeployGitEnv = []string{`SCHEDULER_SYNC_PEAK=*/15 * * * *`}

    payload, err := buildRedeployPayloadFromFlags()
    require.NoError(t, err)
    require.Len(t, payload.Env, 1)
    assert.Equal(t, "*/15 * * * *", payload.Env[0].Value)
}

func TestBuildRedeployPayloadFromFlags_InvalidEnvReturnsError(t *testing.T) {
    redeployGitEnv = []string{"BAD KEY=value"}

    _, err := buildRedeployPayloadFromFlags()
    require.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd -run BuildRedeployPayloadFromFlags -v`
Expected: FAIL por assinatura/fluxo atual sem erro.

**Step 3: Write minimal implementation**

- Alterar assinatura de `buildRedeployPayloadFromFlags` para `(types.StackGitRedeployPayload, error)`.
- Usar `envvars.ParseAssignments(redeployGitEnv)`.
- Em `RunE`, retornar erro amigavel para entrada invalida.

```go
payload, err = buildRedeployPayloadFromFlags()
if err != nil {
    return fmt.Errorf("invalid --env value: %w", err)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cmd -run BuildRedeployPayloadFromFlags -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add cmd/stacks_redeploy_git.go cmd/stacks_redeploy_git_test.go
git commit -m "fix: align redeploy env parsing with spaced values support"
```

### Task 5: Reusar parser no wizard (create/redeploy) e cobrir com teste

**Files:**
- Modify: `internal/wizard/stacks.go`
- Create: `internal/wizard/stacks_env_test.go`
- Test: `internal/wizard/stacks_env_test.go`

**Step 1: Write the failing test**

```go
func TestParseWizardEnvText_AllowsCronExpression(t *testing.T) {
    pairs, err := parseWizardEnvText("SCHEDULER_SYNC_PEAK=*/15 * * * *\nA=1")
    require.NoError(t, err)
    require.Len(t, pairs, 2)
    assert.Equal(t, "*/15 * * * *", pairs[0].Value)
}

func TestParseWizardEnvText_ReturnsLineAwareError(t *testing.T) {
    _, err := parseWizardEnvText("A=1\nBAD KEY=2")
    require.Error(t, err)
    assert.Contains(t, err.Error(), "linha 2")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/wizard -run WizardEnv -v`
Expected: FAIL por funcao inexistente.

**Step 3: Write minimal implementation**

- Extrair parser em helper interno testavel:

```go
func parseWizardEnvText(raw string) ([]types.Pair, error) {
    if strings.TrimSpace(raw) == "" {
        return nil, nil
    }
    lines := strings.Split(strings.TrimSpace(raw), "\n")
    return envvars.ParseAssignments(lines)
}
```

- Reusar helper em `RunCreateSwarmGitWizard` e `RunRedeployGitWizard`, removendo loops duplicados com `strings.SplitN`.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/wizard -run WizardEnv -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/wizard/stacks.go internal/wizard/stacks_env_test.go
git commit -m "refactor: reuse env parser in stack wizards"
```

### Task 6: Atualizar documentacao e validar regressao completa

**Files:**
- Modify: `docs/commands/stacks.md`
- Test: `cmd/stacks_create_swarm_git_test.go`
- Test: `cmd/stacks_redeploy_git_test.go`
- Test: `internal/envvars/parser_test.go`
- Test: `internal/wizard/stacks_env_test.go`

**Step 1: Write the failing test**

Adicionar um caso de documentacao executavel nos testes de cmd (snapshot/assert de exemplo literal) para garantir que o exemplo de cron com espacos vira `Pair` correto.

```go
func TestBuildPayloadFromFlags_CronExpressionExample(t *testing.T) {
    createSwarmGitEnv = []string{`SCHEDULER_SYNC_PEAK=*/15 * * * *`}
    payload, err := buildPayloadFromFlags()
    require.NoError(t, err)
    assert.Equal(t, "*/15 * * * *", payload.Env[0].Value)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd -run CronExpressionExample -v`
Expected: FAIL se o parser nao estiver conectado corretamente.

**Step 3: Write minimal implementation**

Atualizar `docs/commands/stacks.md` com exemplos explicitos para valores com espacos:

```bash
portainer-cli stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1 \
  --env 'SCHEDULER_SYNC_PEAK=*/15 * * * *'
```

E incluir nota: chaves devem ser `A-Z`, `0-9`, `_`; valor pode conter espacos quando devidamente quotado no shell.

**Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS em toda a suite.

**Step 5: Commit**

```bash
git add docs/commands/stacks.md cmd/stacks_create_swarm_git_test.go cmd/stacks_redeploy_git_test.go internal/envvars/parser_test.go internal/wizard/stacks_env_test.go
git commit -m "docs: document env values with spaces and cron expressions"
```

### Riscos e pontos de atencao

- `StringSliceVar` ainda depende de quoting correto no shell; o parser corrige validacao/conteudo, mas nao substitui quoting.
- Alterar assinatura de `buildPayloadFromFlags`/`buildRedeployPayloadFromFlags` impacta testes existentes; ajuste em todos os call sites.
- Evitar `strings.TrimSpace` no valor apos `=` para nao remover espacos intencionais do usuario.
- Manter mensagens de erro consistentes entre create/redeploy/wizard para reduzir ambiguidade em CI.

### Verificacao final recomendada

- `go test ./internal/envvars -v`
- `go test ./internal/wizard -v`
- `go test ./cmd -v`
- `go test ./...`

Plan complete and saved to `docs/plans/2026-03-28-env-vars-with-spaces-stack-create.md`. Two execution options:

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

Which approach?
