# numbertext-go

[![CI](https://github.com/casara/numbertext-go/actions/workflows/ci.yml/badge.svg)](https://github.com/casara/numbertext-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/casara/numbertext-go.svg)](https://pkg.go.dev/github.com/casara/numbertext-go)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Uma biblioteca Go que escreve números por extenso — cardinais, ordinais,
variantes de gênero/caso, negativos e valores monetários — em todos os
idiomas oferecidos pelo [Numbertext.org](https://numbertext.github.io/).

Read this in [English](README.md).

## Por quê

O [Numbertext.org](https://numbertext.github.io/) (projeto original:
[libnumbertext](https://github.com/Numbertext/libnumbertext)) descreve a
conversão de números em texto para dezenas de idiomas como pequenos
arquivos de regras declarativas, escritos em uma DSL minúscula de
reescrita por expressões regulares chamada **Soros**, em vez de código
escrito à mão para cada idioma. Adicionar um idioma — ou uma variante
regional/de gênero de um idioma já existente — significa escrever um
novo arquivo `.sor`, não escrever código novo.

Este repositório é uma implementação em Go, feita do zero, do
interpretador Soros (veja [`internal/soros`](internal/soros)), além de
uma API pública amigável (veja [`numbertext.go`](numbertext.go)),
reaproveitando os mesmos arquivos de regras `.sor` do projeto original
(veja [`data/`](data), sincronizado conforme
[`data/SOURCE.md`](data/SOURCE.md)) sem modificações.

## Instalação

```sh
go get github.com/casara/numbertext-go
```

## Uso

```go
package main

import (
	"fmt"

	numbertext "github.com/casara/numbertext-go"
)

func main() {
	cardinal, _ := numbertext.Cardinal("pt", "12345")
	fmt.Println(cardinal) // "doze mil trezentos e quarenta e cinco"

	ordinal, _ := numbertext.Ordinal("pt", "21")
	fmt.Println(ordinal) // "vigésimo primeiro"

	feminino, _ := numbertext.Convert("pt", "feminine", "2")
	fmt.Println(feminino) // "duas"

	valor, _ := numbertext.Currency("pt", "BRL", "2.50")
	fmt.Println(valor)

	negativo, _ := numbertext.Cardinal("pt", "-5")
	fmt.Println(negativo) // "menos cinco"
}
```

### API

| Função | Finalidade |
| --- | --- |
| `Cardinal(lang, n)` | "123" → número por extenso (cardinal) |
| `Ordinal(lang, n)` | "123" → número por extenso (ordinal) |
| `OrdinalNumber(lang, n)` | "123" → forma abreviada, ex. "123º" (varia por idioma) |
| `Year(lang, n)` | "1999" → ano por extenso |
| `Currency(lang, isoCode, amount)` | estilo NUMBERTEXT(): código ISO 4217 + valor |
| `Money(lang, isoCode, amount)` | estilo MONEYTEXT(): mesma ideia, formatação de cheque |
| `Help(lang)` | seção de ajuda autodescritiva do idioma |
| `Convert(lang, prefix, arg)` | a primitiva sobre a qual tudo o mais é construído — veja abaixo |
| `RegisterLocale(code, sorSource)` | adiciona um idioma em tempo de execução a partir de texto `.sor` |
| `Languages()` | todos os códigos de idioma conhecidos no momento |

Todos os argumentos numéricos são strings (`"123"`, `"-5"`, `"3.14"`),
não `int`/`float64`, porque o motor por trás da biblioteca é baseado em
texto, e isso preserva precisão arbitrária (um arquivo `.sor` consegue
escrever até decilhões por extenso). `CardinalInt(lang, n int64)` é
oferecida como atalho para o caso comum.

### Variantes regionais (ex.: "en-GB", "pt-BR")

Um código regional não tem (nem precisa ter) um arquivo `.sor` próprio:
`data/en.sor` e `data/pt.sor` trazem as regras do idioma-base mais
algumas linhas marcadas com `[:en-GB:]`, `[:pt-BR:]` etc. para diferenças
de escrita por região. Basta passar o código completo —
`Languages()` continua listando só os códigos-base (`"en"`, `"pt"`, ...),
já que a região não é um arquivo separado, mas qualquer string no
formato `"<base>-<REGIÃO>"` funciona e ativa automaticamente as linhas
marcadas daquela região:

```go
gb, _ := numbertext.Cardinal("en-GB", "101")
fmt.Println(gb) // "one hundred and one"

us, _ := numbertext.Cardinal("en", "101")
fmt.Println(us) // "one hundred one" (sem "and" — escrita padrão)

br, _ := numbertext.Cardinal("pt-BR", "16")
fmt.Println(br) // "dezesseis"

pt, _ := numbertext.Cardinal("pt", "16")
fmt.Println(pt) // "dezasseis" (português europeu, padrão do arquivo)
```

Se a parte antes do primeiro `-` não for um idioma conhecido, `Cardinal`
(e as demais funções) retornam erro — não há arquivo para usar como
fallback.

### Escolhendo quais idiomas são embutidos

Por padrão, `go build` embute os 52 arquivos `.sor` no binário via
`go:embed` (veja `locale_all.go`), independentemente de quais idiomas seu
programa realmente usa. Se você só precisa de alguns, dá pra optar por um
binário menor com build tags: passe `numbertext_select` mais uma tag
`numbertext_lang_<code>` para cada idioma desejado (`<code>` é o nome-base
do arquivo `.sor`, em minúsculas, ex.: `en`, `pt`, `hu_hung`):

```sh
go build -tags "numbertext_select numbertext_lang_en numbertext_lang_pt" .
```

Isso embute só `en` e `pt` — `numbertext.Languages()` retorna exatamente
`["en", "pt"]`, e chamar `Cardinal("de", ...)`, por exemplo, retorna erro
de "idioma desconhecido". Passar `numbertext_select` **sem** nenhuma tag
`numbertext_lang_<code>` não embute nada (`Languages()` retorna slice
vazio); só faz sentido se você planeja fornecer todos os locales via
`RegisterLocale`.

Não passar `numbertext_select` (o padrão) sempre embute todos os idiomas
— essa opção nunca muda o comportamento a menos que você opte por ela. Os
arquivos por idioma (`locale_<code>.go`) são gerados por
`scripts/gen-locale-embeds.py`; regenere-os depois que
`scripts/sync-data.sh` adicionar ou remover um `data/*.sor`.

### Gênero, caso gramatical e outras variantes específicas de idioma

Os arquivos de regras do Numbertext não se limitam a um formato fixo de
cardinal/ordinal — um idioma pode definir seções nomeadas arbitrárias
para qualquer variante gramatical de que precise (`pt.sor` tem
`feminine`/`masculine`; `ru.sor` tem
`cardinal-feminine`/`cardinal-neuter`; `it.sor` tem
`ordinal-masculine`; e assim por diante). Esses nomes de seção **não**
são padronizados entre idiomas, então esta biblioteca não finge o
contrário com uma função genérica `Gender(...)` que serviria para todos.
Em vez disso, chame `Convert` com o nome da seção que você precisa:

```go
fem, _ := numbertext.Convert("pt", "feminine", "2")
fmt.Println(fem) // "duas"
```

Execute `Help(lang)` (ou abra o arquivo `.sor` diretamente) para ver
quais seções um determinado idioma define.

### Adicionando um novo idioma

Como toda a biblioteca é orientada a dados, não é preciso tocar em
nenhum código Go para adicionar um idioma:

```go
sorSource := "1 uno\n2 dos\n3 tres\n" // um arquivo real precisa de muito mais regras
err := numbertext.RegisterLocale("es-mini", sorSource)
```

ou adicione um arquivo em [`data/`](data) — veja
[`data/SOURCE.md`](data/SOURCE.md#adding-a-new-language) para um guia
rápido e um link para a especificação da linguagem Soros.

## Estabilidade da API

O `numbertext-go` é lançado como **v0.x**, e em Semantic Versioning isso
carrega um significado específico que vale a pena deixar explícito em vez
de deixar você inferir: **uma versão minor pode quebrar a API pública.**
Fixe uma versão, leia o [CHANGELOG](CHANGELOG.md) antes de atualizar, e
espere fazer pequenos ajustes mecânicos quando o fizer.

O que conta como API pública: todo símbolo exportado no pacote raiz
`numbertext` e as flags do `cmd/numbertext`. Não é público: qualquer
coisa sob `internal/` (incluindo o próprio interpretador Soros) e o texto
exato das mensagens de erro — `ErrUnknownLanguage`/`ErrEmptyLocaleCode`
(pensados para `errors.Is`) são a parte estável do tratamento de erros;
não há suporte para checar o texto da mensagem.

Mudanças que quebram compatibilidade serão listadas em `### Changed` ou
`### Removed` no changelog, com a migração na mesma entrada. Quando um
rename puder manter a grafia antiga compilando, ele manterá — como um
alias descontinuado.

## Skill para assistentes de IA

[skills/using-numbertext-go/SKILL.md](skills/using-numbertext-go/SKILL.md)
documenta, em um formato que ferramentas de IA podem carregar sob
demanda (Claude Skill), como usar o `numbertext-go` de forma idiomática:
as funções principais, o primitivo `Convert`, variantes regionais,
seções de gênero/caso e a incorporação seletiva de idiomas. Para usá-la
em um projeto que depende do `numbertext-go`, copie o diretório para
dentro dele:

```sh
cp -r skills/using-numbertext-go <seu-projeto>/.claude/skills/using-numbertext-go
```

## Desenvolvimento

`make help` lista todos os comandos. Antes de abrir um PR:

```sh
make check     # lint + arch-lint + verify-docs + testes com detector de race + test-select
make lint-md   # lint de markdown (requer Node.js >= 20)
```

Convenções de estilo estão documentadas em
[docs/coding-style.pt-BR.md](docs/coding-style.pt-BR.md).

Veja [CONTRIBUTING.pt-BR.md](CONTRIBUTING.pt-BR.md) para o fluxo
completo de contribuição (convenção de commits, processo de
branch/PR) e [AGENTS.md](AGENTS.md) (em inglês) para uma visão mais
profunda da arquitetura, voltada para colaboradores humanos e agentes
de código (`CLAUDE.md` só o importa, para o Claude Code ler as mesmas
instruções).

## Licença

O código Go deste repositório é licenciado sob MIT — veja
[LICENSE](LICENSE). Os arquivos de regras `.sor` em [`data/`](data) são
copiados sem modificações do projeto original libnumbertext e permanecem
sob a licença BSD-3-Clause original — veja [`data/LICENSE`](data/LICENSE)
e [`data/SOURCE.md`](data/SOURCE.md).
