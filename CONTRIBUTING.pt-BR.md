# Contribuindo

*Read this in [English](CONTRIBUTING.md)*

Este projeto segue o [Código de Conduta](CODE_OF_CONDUCT.md). Ao
participar, espera-se que você o siga.

## Antes de começar

Leia, nesta ordem:

1. [AGENTS.md](AGENTS.md) (em inglês) - arquitetura, a linguagem Soros
   resumida, limitações conhecidas/deliberadas, decisões de design
   não óbvias.
2. [docs/coding-style.pt-BR.md](docs/coding-style.pt-BR.md) - convenções
   de código (formatação, tratamento de erros/wrapcheck, dependências
   explícitas).
3. [README.md](README.md) / [README.pt-BR.md](README.pt-BR.md) - a API
   pública e exemplos de uso.
4. `data/SOURCE.md` (em inglês) - como os arquivos `.sor` de idiomas são
   obtidos e como adicionar um novo.

## Ambiente

`make help` lista todos os comandos. Antes de abrir um PR, rode:

```sh
make check     # lint + arch-lint + verify-docs + testes com detector de race + test-select
make lint-md   # lint de markdown (requer Node.js >= 20)
```

É o mínimo que o CI roda em cada PR
([.github/workflows/ci.yml](.github/workflows/ci.yml)) - dividido em
jobs separados lá (`test` por versão do Go, `lint` uma única vez, mais
o lint de markdown) em vez de uma chamada única a `make check`; veja a
seção "Commands" do AGENTS.md (em inglês) para entender o motivo. Se a
mudança altera
comportamento — não é só refatoração — inclua um teste cobrindo o novo
caso e, para um idioma novo/alterado, algumas conversões conhecidas
**verificadas rodando o código de fato** (`go run ./cmd/numbertext
-lang <x> -cardinal <n>` ou um teste descartável), não adivinhadas —
veja a nota "Verify locale output before writing it down" em
AGENTS.md.

## Fluxo de branches

Rebase, não merge commits: atualize sua branch com `git rebase` contra
a base antes de abrir/atualizar um PR, em vez de fazer merge da base
na sua branch. Histórico linear, sem merge commits.

A convenção de commit abaixo é obrigatória na `main`. Uma branch de
trabalho que será squashed no merge não precisa segui-la à risca — o
histórico intermediário não é o que permanece.

## Mensagens de commit

[Conventional Commits](https://www.conventionalcommits.org/), em
inglês, no imperativo:

```text
<tipo>(<escopo>): <descrição>
```

* `<descrição>` no imperativo, descrevendo a ação (`add`, `fix`,
  `remove`, `harden`, `resolve`, `guarantee`), não o estado resultante
  (`added`) ou o histórico (`fixed`, "adds/added").
* `<escopo>` é o pacote ou área afetada. Não é uma lista fechada, mas
  os mais comuns são: `soros` (o interpretador), `numbertext` (a API
  pública/registro), `cmd` (a CLI), `data` (os arquivos `.sor` de
  idiomas) — mais os transversais `release`, `docs`, `ci`, `deps`.

### Tipos

| Tipo | Quando usar |
| --- | --- |
| `feat` | Novo recurso (semver MINOR). |
| `fix` | Correção de bug (semver PATCH). |
| `docs` | Só documentação (README, comentários) — sem mudança de código. |
| `test` | Só teste (adicionado, alterado ou removido) — sem mudança de código de produção. |
| `refactor` | Muda como o código é escrito/organizado sem mudar o comportamento observável. |
| `perf` | Mudança cujo propósito é performance. |
| `style` | Formatação, ponto e vírgula, espaço em branco, lint — sem mudança de código. |
| `build` | Build e dependências (`go.mod`, `Makefile`, ferramentas). |
| `ci` | Integração contínua (`.github/workflows/`). |
| `chore` | Tarefas de manutenção que não se encaixam acima (config, `.gitignore`, ...). |
| `cleanup` | Remove código comentado, morto ou desnecessário — sem mudança de comportamento. |
| `remove` | Remove um arquivo, diretório ou recurso obsoleto/não usado. |
| `raw` | Mudança em arquivo de configuração/dado/parâmetro que não se encaixa acima (ex.: sincronizar `data/*.sor`). |

### Exemplos

```text
feat(numbertext): add RegisterLocale for adding a language at runtime
fix(soros): propagate the leading boundary through a piped call argument
docs(readme): explain regional variant fallback for en-GB/pt-BR
raw(data): sync data/*.sor to libnumbertext@a4b0225
```

(as mensagens em si ficam em inglês, seguindo a convenção acima — os
exemplos mostram o formato esperado.)

## Alterando a API pública

O projeto está em **v0.x**, então uma mudança quebrando compatibilidade
é permitida — veja "API stability" no README para o que isso significa
para quem usa a biblioteca. Ainda assim precisa ser deliberada e
visível:

* Diga isso na descrição do PR e adicione uma entrada `### Changed` ou
  `### Removed` no `CHANGELOG.md`, com a migração na mesma entrada.
* Quando um rename puder manter a grafia antiga compilando, mantenha-a
  como alias descontinuado em vez de simplesmente remover.
* Prefira uma mudança aditiva quando existir uma: um parâmetro opcional
  novo via opção variádica, uma função nova ao lado da antiga.

Mensagens de erro explicitamente **não** fazem parte da API pública;
não há suporte para checar o texto delas. Erros sentinela
(`ErrUnknownLanguage`, `ErrEmptyLocaleCode`, pensados para `errors.Is`)
fazem parte.

## Contribuindo com um assistente de IA

Este projeto é construído com um, e não há nada a esconder ou declarar
sobre isso: o código é julgado da mesma forma de qualquer jeito, e um
PR não é marcado.

O que se pede é a mesma coisa pedida de qualquer pessoa:

* **Entenda o que você está enviando.** Se você não consegue explicar
  por que uma mudança está correta, ela não está pronta — independente
  de quem/o que a escreveu.
* **Verifique a saída do idioma, não adivinhe.** Um arquivo `.sor`
  codifica escolhas estilísticas reais e às vezes surpreendentes por
  idioma (veja o próprio exemplo em AGENTS.md: `en` vs `en-GB` com
  "and", `Cardinal` vs `Year` removendo o "and"). Rode o código antes
  de escrever um valor esperado.
* **Não deixe uma ferramenta reabrir uma decisão já tomada.**
  `AGENTS.md` documenta as limitações conhecidas e deliberadas deste
  port (as duas classes de regra do upstream que ele não consegue
  expressar) porque esse limite foi alcançado uma vez, de propósito —
  um PR "corrigindo" isso adivinhando o comportamento faltante é
  fechado com uma referência, não um debate.
* **Rode as checagens localmente.** `make check` e `make lint-md`,
  antes de abrir o PR em vez de depois que o CI apontar.

`AGENTS.md` é lido diretamente por Codex e Cursor; `CLAUDE.md` o importa
para o Claude Code. Qualquer coisa que você adicionar para um deles
pertence a `AGENTS.md`, para que os outros também recebam.

## Pull requests

* `make check` passar é obrigatório, não opcional. Se a mudança tocar
  em algum arquivo Markdown, `make lint-md` também.
* Um PR pequeno e focado em uma mudança é preferível a um PR grande
  cobrindo várias coisas não relacionadas — mas isso é julgamento, não
  uma regra estrita.
* Se a mudança altera comportamento documentado, atualize a
  documentação no mesmo PR.
