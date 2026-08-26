# Guia de Estilo Go

*[Read in English](coding-style.md)*

## Objetivo

Este documento define convenções adicionais adotadas pelo projeto além das regras já aplicadas automaticamente por
ferramentas como `gofmt`, `goimports` e `golangci-lint`.

Sempre que possível, as decisões devem priorizar:

* legibilidade;
* consistência;
* facilidade de manutenção;
* qualidade dos diffs.

---

## Formatação

Todo código deve ser formatado com:

* gofmt
* goimports

Não devem ser realizadas alterações manuais para contrariar a formatação produzida por essas ferramentas.

---

## Legibilidade acima da concisão

Prefira código explícito e fácil de entender em vez de versões excessivamente compactas.

Preferível:

```go
if err != nil {
    return err
}
```

Evite:

```go
if err != nil { return err }
```

---

## Chamadas de função

Chamadas simples devem permanecer em uma única linha.

```go
prog, err := reg.program(lang)

groups := match.Groups()
```

Chamadas com múltiplos argumentos, opções ou estruturas aninhadas devem utilizar o formato vertical.

```go
return nil, fmt.Errorf(
    "%w %q (see Languages for the supported base codes; "+
        "a regional variant like %q must build on one of them, e.g. \"en-GB\" on \"en\")",
    ErrUnknownLanguage, lang, lang,
)
```

---

## Qualidade dos diffs

Sempre que uma construção possuir tendência natural de crescimento, prefira o formato vertical.

Exemplo:

```go
declare -a targets=(
    "README.md|the prose around Usage/Regional variants/Gender"
    "CHANGELOG.md|every public change belongs here"
)
```

Esse formato reduz conflitos de merge e produz diffs menores quando novos elementos são adicionados.

---

## Comentários

Comentários devem explicar:

* propósito;
* comportamento;
* limitações;
* decisões de projeto.

Comentários que apenas repetem o nome da função devem ser evitados.

Ruim:

```go
// Compile compiles a program.
```

Melhor:

```go
// Compile parses Soros source code into a runnable Program for the
// given language code. It only fails outright if the source cannot be
// tokenized at all; individual rules it cannot compile are recorded in
// Program.Skipped instead of failing the whole locale.
```

---

## Tratamento de Erros (wrapcheck)

Erros retornados por dependências externas ou por camadas inferiores devem receber contexto adicional antes de
serem propagados.

O objetivo é tornar a origem da falha evidente para quem a captura com `errors.Is`, e facilitar o diagnóstico
quando um arquivo `.sor` falha ao compilar.

### Regra

Ao retornar um erro recebido de outra função, adicione contexto usando `fmt.Errorf` e `%w`.

Correto:

```go
prog, err := soros.Compile(source, lang)
if err != nil {
    return nil, fmt.Errorf("numbertext: compiling locale %q: %w", lang, err)
}
```

Evite:

```go
prog, err := soros.Compile(source, lang)
if err != nil {
    return nil, err
}
```

---

### Mensagem de erro

A mensagem deve descrever a operação que falhou, não repetir o texto do erro original.

Correto:

```go
return fmt.Errorf("numbertext: RegisterLocale(%q): %w", code, err)
```

Evite:

```go
return fmt.Errorf("error: %w", err)
```

```go
return fmt.Errorf("failed: %w", err)
```

```go
return fmt.Errorf("unexpected error: %w", err)
```

---

### Quando o wrap não é necessário

Não faça wrap quando:

* estiver criando um erro novo;
* estiver retornando um erro sentinela (ex.: `ErrUnknownLanguage`, `ErrEmptyLocaleCode` em `errors.go`);
* o erro já contém contexto suficiente e a camada atual não adiciona informação relevante.

Exemplos:

```go
return ErrEmptyLocaleCode
```

```go
return errUnterminated
```

---

### Orientação para IAs

Ao corrigir violações de `wrapcheck`:

1. Preserve a cadeia de erro usando `%w`.
2. Descreva a operação que falhou (ex.: "compiling locale %q", não "error compiling").
3. Não utilize mensagens genéricas como:

   * "error"
   * "failed"
   * "unexpected error"
4. Não repita contexto já presente em camadas inferiores.
5. Prefira mensagens curtas em minúsculas, prefixadas com o nome do pacote onde o código já existente já faz isso
   (`"numbertext: ..."`).
6. Inclua identificadores relevantes quando agregarem valor:

```go
return fmt.Errorf("numbertext: compiling locale %q: %w", lang, err)
```

---

## Dependências

Prefira dependências explícitas via parâmetros de função e valores construídos em vez de variáveis globais.

Preferível:

```go
func compilePattern(wrapped string) (matcher, error)
```

Evite:

```go
var globalCompiledPattern matcher
```

Salvo quando a natureza do componente justificar claramente um singleton compartilhado, como `reg` (o registro
de locales em nível de pacote, em `locales.go`) faz: ele existe uma vez por processo, é seguro para uso
concorrente, e é exatamente o identificador que `RegisterLocale`/`Languages`/`Convert` precisam compartilhar.
